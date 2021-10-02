use crate::post::{CreatePost, FullUpdatePost, MutatorClient, MutatorError};
use nats;
use serde::Serialize;
use std::fmt;

/// CreatePostEvent is self described
#[derive(Serialize)]
pub struct CreatePostEvent {
    event_name: String,
    #[serde(flatten)]
    payload: CreatePost,
}

/// UpdatePostEvent is self described
#[derive(Serialize)]
pub struct UpdatePostEvent {
    event_name: String,
    #[serde(flatten)]
    payload: FullUpdatePost,
}

impl CreatePostEvent {
    /// Creates a CreatePostEvent with default event_name
    fn new(payload: CreatePost) -> Self {
        CreatePostEvent {
            event_name: "posts.v1.create".to_string(),
            payload,
        }
    }
}
impl UpdatePostEvent {
    /// Creates an UpdatePostEvent with default event_name
    fn new(payload: FullUpdatePost) -> Self {
        UpdatePostEvent {
            event_name: "posts.v1.update".to_string(),
            payload,
        }
    }
}

/// Config encodes the necessary data for a NATS connection
pub struct Config {
    pub user: String,
    pub pass: String,
    pub subject: String,
    pub host: String,
    pub port: String,
}

impl Config {
    /// URL for connecting to the corresponding NATS server
    pub fn url(&self) -> String {
        format!("nats://{}:{}", self.host, self.port)
    }
}

/// ClientError wraps any errors returned by the underlying NATS client
#[derive(Debug)]
pub struct ClientError {
    message: String,
}

impl std::error::Error for ClientError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for ClientError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "nats client error: {}", self.message)
    }
}

/// NATS client's wrapper
pub struct Client {
    config: Config,
    conn: nats::Connection,
}

impl Client {
    /// Creates a new NATS' client with passed config's data
    pub fn connect(config: Config) -> Result<Client, ClientError> {
        let mut options = nats::Options::new();
        if !config.user.is_empty() {
            options = nats::Options::with_user_pass(&config.user[..], &config.pass[..]);
        }
        match options.connect(&config.url()[..]) {
            Ok(conn) => Ok(Client { config, conn }),
            Err(e) => Err(ClientError {
                message: e.to_string(),
            }),
        }
    }

    fn send<T: Serialize>(&self, payload: T) -> Result<(), MutatorError> {
        let payload_str = serde_json::to_string(&payload).unwrap();
        match bincode::serialize(&payload_str) {
            // TODO Ignoring first element as it's sent like 'D' (coming from bincode)
            Ok(bytes) => match self.conn.publish(&self.config.subject[..], &bytes[1..]) {
                Ok(()) => Ok(()),
                Err(e) => Err(MutatorError {
                    message: e.to_string(),
                }),
            },
            Err(e) => Err(MutatorError {
                message: e.to_string(),
            }),
        }
    }
}

impl MutatorClient<CreatePost> for Client {
    /// Implements the send method for post creation
    fn send(&self, payload: CreatePost) -> Result<(), MutatorError> {
        let create_post = CreatePostEvent::new(payload);
        self.send(create_post)
    }
}

impl MutatorClient<FullUpdatePost> for Client {
    /// Implements the send method for post updating
    fn send(&self, payload: FullUpdatePost) -> Result<(), MutatorError> {
        let update_post = UpdatePostEvent::new(payload);
        self.send(update_post)
    }
}

mod test {
    use super::Config;

    #[test]
    fn test_url() {
        let config = Config {
            user: "usr".to_string(),
            pass: "pass".to_string(),
            subject: "not_important".to_string(),
            host: "127.0.0.1".to_string(),
            port: "4222".to_string(),
        };
        assert_eq!(config.url(), "nats://127.0.0.1:4222".to_string())
    }
}
