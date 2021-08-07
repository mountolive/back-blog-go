use crate::post::{CreatePost, MutatorClient, MutatorError, UpdatePost};
use nats;
use serde::Serialize;
use std::fmt;

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
        match bincode::serialize(&payload) {
            Ok(bytes) => match self.conn.publish(&self.config.subject[..], &bytes[..]) {
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
        self.send(payload)
    }
}

impl MutatorClient<UpdatePost> for Client {
    /// Implements the send method for post updating
    fn send(&self, payload: UpdatePost) -> Result<(), MutatorError> {
        self.send(payload)
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
