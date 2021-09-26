mod auth;
mod grpc_authenticator;
mod http_handler;
mod mem_storage;
mod nats_client;
mod post;
mod post_reader;
mod store;
mod user;

use auth::AuthService;
use grpc_authenticator::{create_grpc_login_client, GRPCAuthenticator};
use http_handler::HTTPHandler;
use mem_storage::MemStorageDriver;
use nats_client::{Client, Config};
use parking_lot::RwLock;
use post::{PostCreator, PostUpdater};
use post_reader::{PostReader, ReaderClientConfig};
use std::collections::HashMap;
use std::env;
use std::net::SocketAddr;
use std::str::FromStr;
use store::JWTStore;

fn parse_nats_config() -> Config {
    Config {
        user: env::var("POST_SERVER_USER").expect("post server's user not set"),
        pass: env::var("POST_SERVER_PASSWORD").expect("post server's password not set"),
        subject: env::var("POST_SERVER_SUBJECT").expect("post server's subject not set"),
        host: env::var("POST_SERVER_HOST").expect("post server's host not set"),
        port: env::var("POST_SERVER_PORT").expect("post server's port not set"),
    }
}

#[tokio::main]
async fn main() {
    // Setup authenticator
    let grpc_srv_addr = env::var("USER_SERVICE_ADDRESS").expect("user service address not set");
    let user_srv_addr = grpc_srv_addr
        .parse::<http::Uri>()
        .expect("malformed user service's address");
    let authenticator = GRPCAuthenticator::new(create_grpc_login_client(user_srv_addr).await);

    let storage_driver = MemStorageDriver {
        data: RwLock::new(HashMap::new()),
    };
    let store = JWTStore {
        storage: Box::new(storage_driver),
    };

    let env_ttl = env::var("TOKEN_TTL").expect("token ttl not set");
    let ttl = env_ttl.parse::<u64>().expect("wrong ttl value set");
    let secret = env::var("TOKEN_SALT").expect("token salt not set");

    let auth_service =
        AuthService::new(Box::new(authenticator), Box::new(store), ttl, &secret[..]).unwrap();

    // Setup post creator and updater
    let nats_client_creator = Client::connect(parse_nats_config()).unwrap();
    let nats_client_updater = Client::connect(parse_nats_config()).unwrap();

    let post_creator = PostCreator {
        client: Box::new(nats_client_creator),
    };
    let post_updater = PostUpdater {
        client: Box::new(nats_client_updater),
    };

    // Setup post reader
    let post_reader = PostReader::new(ReaderClientConfig::with_default(
        env::var("POST_REST_API_URL").expect("post's rest api url not set"),
    ))
    .unwrap();

    // HTTP handler
    let forever_handler = Box::leak(Box::new(HTTPHandler {
        auth: auth_service,
        creator: post_creator,
        updater: post_updater,
        reader: post_reader,
    }));

    let port = env::var("GATEWAY_PORT").expect("port not set");
    let address =
        SocketAddr::from_str(&format!("127.0.0.1:{}", port)[..]).expect("malformed server address");
    forever_handler.start(address).await;
}
