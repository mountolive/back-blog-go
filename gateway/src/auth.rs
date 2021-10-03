//! Provides basic authentication and authorization API surface

use crate::user::login_client::LoginClient;
use hmac::{Hmac, NewMac};
use jwt::{Header, SignWithKey, Token, VerifyWithKey};
use serde::{Deserialize, Serialize};
use sha2::Sha256;
use std::collections::BTreeMap;
use std::fmt;
use std::time::{Duration, SystemTime};

/// Error that wraps when a JWTToken action goes wrong
#[derive(Debug)]
pub struct TokenError {
    message: String,
}

impl std::error::Error for TokenError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for TokenError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "jwt token error: {}", self.message)
    }
}

/// Basic implementation of a JWT token with TTL
#[derive(Serialize, Deserialize)]
pub struct JWTToken {
    pub value: String,
    until: Duration,
}

static WRONG_NOW: &str = "unable to determine now's timestamp";
static USER_KEY: &str = "user";

impl JWTToken {
    /// Creates a new JWT token from the username passed and with the passed TTL (in seconds)
    pub fn generate(username: &str, ttl: u64, key: &Hmac<Sha256>) -> Result<JWTToken, TokenError> {
        let mut claims = BTreeMap::new();
        claims.insert(USER_KEY, username);
        match claims.sign_with_key(key) {
            Ok(token_value) => match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
                Ok(now) => Ok(JWTToken {
                    value: token_value,
                    until: now + Duration::from_secs(ttl),
                }),
                Err(_) => Err(TokenError {
                    message: String::from(WRONG_NOW),
                }),
            },
            Err(_) => Err(TokenError {
                message: String::from("signing token"),
            }),
        }
    }

    /// Explodes and a gets underlying username from JWT token
    fn get_username(token_value: &str, key: &Hmac<Sha256>) -> Result<String, TokenError> {
        let token: Token<Header, BTreeMap<String, String>, _>;
        match VerifyWithKey::verify_with_key(token_value, key) {
            Ok(val) => {
                token = val;
                Ok(token.claims()[USER_KEY].to_string())
            }
            Err(_) => Err(TokenError {
                message: String::from("unable to retrieve username"),
            }),
        }
    }

    /// Checks whether a token is already evicted
    fn is_evicted(&self) -> Result<bool, TokenError> {
        match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
            Ok(now) => Ok(now > self.until),
            Err(_) => Err(TokenError {
                message: String::from(WRONG_NOW),
            }),
        }
    }
}

impl PartialEq for JWTToken {
    fn eq(&self, other: &Self) -> bool {
        self.value == other.value
    }
}

/// Error returned when an action related to the token's store goes wrong
#[derive(Debug)]
pub struct TokenStoreError {
    pub message: String,
}

impl std::error::Error for TokenStoreError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for TokenStoreError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "token store error: {}", self.message)
    }
}

/// Wrapping error for authetication and authorization actions
#[derive(Debug)]
pub struct AuthenticationError {
    pub message: String,
}

impl std::error::Error for AuthenticationError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for AuthenticationError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "authentication error: {}", self.message)
    }
}

/// Describes the basic contract expected by a Token's store
pub trait TokenStore {
    fn retrieve(&self, key: &str) -> Result<JWTToken, TokenStoreError>;
    fn save(&self, key: &str, ser_token: &str) -> Result<(), TokenStoreError>;
}

/// Handles login and authorization by means of an Authenticator and a TokenStore
pub struct AuthService {
    // It's crap to pass a concrete dep
    login_client: LoginClient<tonic::transport::Channel>,
    store: Box<dyn TokenStore>,
    token_ttl: u64,
    token_key: Hmac<Sha256>,
}

// Marking it as thread-safe
unsafe impl Send for AuthService {}
unsafe impl Sync for AuthService {}

impl AuthService {
    /// Creates a new AuthService
    pub fn new(
        // This is crap :(
        login_client: LoginClient<tonic::transport::Channel>,
        store: Box<dyn TokenStore>,
        token_ttl: u64,
        secret: &str,
    ) -> Result<AuthService, AuthenticationError> {
        let token_key: Hmac<Sha256>;
        match Hmac::new_from_slice(secret.as_bytes()) {
            Ok(key) => token_key = key,
            Err(_) => {
                return Err(AuthenticationError {
                    message: String::from("invalid key length"),
                })
            }
        }
        Ok(AuthService {
            login_client,
            store,
            token_ttl,
            token_key,
        })
    }

    /// Authenticates user against authentication service and returns a JWT token value
    pub async fn login(&self, usr: &str, pass: &str) -> Result<String, AuthenticationError> {
        match crate::grpc_authenticator::authenticate(self.login_client.clone(), usr, pass).await {
            Ok(logged_in) => {
                if !logged_in {
                    return Err(AuthenticationError {
                        message: String::from("invalid credentials"),
                    });
                }
                let token: JWTToken;
                match JWTToken::generate(usr, self.token_ttl, &self.token_key) {
                    Ok(t) => token = t,
                    Err(e) => return Err(AuthenticationError { message: e.message }),
                }
                // TODO Handle serializing Result for login's token
                let serialized_token = serde_json::to_string(&token).unwrap();
                match self.store.save(&usr, &serialized_token[..]) {
                    Ok(_) => Ok(token.value),
                    Err(e) => Err(AuthenticationError { message: e.message }),
                }
            }
            Err(e) => Err(e),
        }
    }

    /// Checks whether the received token is still authorized
    pub fn authorize(&self, token_value: &str) -> Result<bool, AuthenticationError> {
        match JWTToken::get_username(token_value, &self.token_key) {
            Ok(usr) => match self.store.retrieve(&usr[..]) {
                Ok(saved_token) => match saved_token.is_evicted() {
                    Ok(evicted) => Ok(evicted),
                    Err(e) => Err(AuthenticationError { message: e.message }),
                },
                Err(e) => Err(AuthenticationError { message: e.message }),
            },
            Err(e) => Err(AuthenticationError { message: e.message }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::thread;

    const SECRET: &str = "un secreto";
    const USER: &str = "noice";

    fn generate_key(test_case: fn(token_key: &Hmac<Sha256>)) {
        match Hmac::new_from_slice(SECRET.as_bytes()) {
            Ok(key) => test_case(&key),
            Err(e) => assert!(
                false,
                "shouldn't return an Err result when creating the key: {}",
                e,
            ),
        }
    }

    #[test]
    fn test_generate_token_correct() {
        let correct = |key: &Hmac<Sha256>| match JWTToken::generate("test", 1000, key) {
            Ok(_) => assert!(true, "should return an Ok result"),
            Err(e) => assert!(false, "shouldn't return an Err {}", e),
        };
        generate_key(correct)
    }

    #[test]
    fn test_get_username_error() {
        let test_case = |key: &Hmac<Sha256>| {
            let exp_err_msg = String::from("unable to retrieve username");
            match Hmac::new_from_slice("whatever".as_bytes()) {
                Ok(other_key) => {
                    let username = "whatever";
                    match JWTToken::generate(username, 2000, &other_key) {
                        Ok(token) => match JWTToken::get_username(&token.value[..], key) {
                            Ok(_) => assert!(
                                false,
                                "shouldn't return Ok result when evaluating with wrong key"
                            ),
                            Err(e) => assert_eq!(e.message, exp_err_msg),
                        },
                        Err(_) => {
                            assert!(false, "shoudln't return an error when creating mock token")
                        }
                    }
                }
                Err(_) => assert!(false, "shouldn't return an error when creating mock key"),
            }
        };
        generate_key(test_case)
    }

    #[test]
    fn test_get_username_correct() {
        let test_case = |key: &Hmac<Sha256>| match JWTToken::generate(USER, 2000, key) {
            Ok(token) => match JWTToken::get_username(&token.value[..], key) {
                Ok(found) => assert_eq!(found, USER),
                Err(_) => assert!(false, "shouldn't return an Err while verifying the token"),
            },
            Err(_) => {
                assert!(false, "shoudln't return an error when creating mock token")
            }
        };
        generate_key(test_case)
    }

    #[test]
    fn test_is_evicted_before_now() {
        let test_case = |key: &Hmac<Sha256>| match JWTToken::generate(USER, 5, key) {
            Ok(token) => {
                thread::sleep(Duration::from_secs(2));
                let is_evicted = token.is_evicted().unwrap();
                assert!(!is_evicted, "should return true");
            }
            Err(_) => {
                assert!(false, "shoudln't return an error when creating mock token")
            }
        };
        generate_key(test_case)
    }

    #[test]
    fn test_is_evicted_past_now() {
        let test_case = |key: &Hmac<Sha256>| match JWTToken::generate(USER, 2, key) {
            Ok(token) => {
                thread::sleep(Duration::from_secs(3));
                let is_evicted = token.is_evicted().unwrap();
                assert!(is_evicted, "should return false");
            }
            Err(_) => {
                assert!(false, "shoudln't return an error when creating mock token")
            }
        };
        generate_key(test_case)
    }
}
