use crate::post::{CreatePost, MutatorClient, MutatorError, UpdatePost};
use nats;
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
}

impl Client {
    /// Creates a new NATS' client with passed config's data
    pub fn connect(config: Config) -> Result<Client, ClientError> {
        Ok(Client { config })
    }
}

impl MutatorClient<CreatePost> for Client {
    fn send(&self, _: CreatePost) -> Result<(), MutatorError> {
        Ok(())
    }
}

impl MutatorClient<UpdatePost> for Client {
    fn send(&self, _: UpdatePost) -> Result<(), MutatorError> {
        Ok(())
    }
}

mod test {
    use super::*;

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
