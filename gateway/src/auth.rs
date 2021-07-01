//! Provides basic authentication and authorization API surface

use hmac::{Hmac, NewMac};
use jwt::{Header, SignWithKey, Token, VerifyWithKey};
use sha2::Sha256;
use std::collections::BTreeMap;
use std::fmt;
use std::time::{Duration, SystemTime};

// Error that wraps when a JWTToken action goes wrong
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

// Basic implementation of a JWT token with TTL
pub struct JWTToken {
    pub value: String,
    until: Duration,
}

const WRONG_NOW: &str = "unable to determine now's timestamp";
const USER_KEY: &str = "user";

// TODO make JWTToken implement Serialize
impl<'a> JWTToken {
    // Creates a new JWT token from the username passed and with the passed TTL (in seconds)
    pub fn generate(username: &str, ttl: u64, key: &Hmac<Sha256>) -> Result<JWTToken, TokenError> {
        let mut claims = BTreeMap::new();
        claims.insert(USER_KEY, username);
        match claims.sign_with_key(key) {
            Ok(token_value) => match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
                Ok(now) => Ok(JWTToken {
                    value: token_value,
                    until: now + Duration::from_secs(ttl),
                }),
                Err(_) => {
                    return Err(TokenError {
                        message: String::from(WRONG_NOW),
                    })
                }
            },
            Err(_) => {
                return Err(TokenError {
                    message: String::from("signing token"),
                })
            }
        }
    }

    // Explodes and a gets underlying username from JWT token
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

    // Checks whether a token is already evicted
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

// Error returned when an action related to the token's store goes wrong
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

// Wrapping error for authetication and authorization actions
#[derive(Debug)]
pub struct AuthenticationError {
    message: String,
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

// Describes the basic contract expected by a Token's store
pub trait TokenStore<'a> {
    fn retrieve(&'a self, key: &str) -> Result<&'a JWTToken, TokenStoreError>;
    fn save(&mut self, key: &str, token: JWTToken) -> Result<(), TokenStoreError>;
}

// Describes the basic contract expected by an authenticator's client
pub trait Authenticator {
    fn authenticate(&self, username: &str, password: &str) -> Result<bool, AuthenticationError>;
}

// Handles login and authorization by means of an Authenticator and a TokenStore
pub struct AuthService<'a> {
    authenticator: Box<dyn Authenticator>,
    store: Box<dyn TokenStore<'a>>,
    token_ttl: u64,
    token_key: Hmac<Sha256>,
}

impl<'a> AuthService<'a> {
    // Creates a new AuthService
    pub fn new(
        authenticator: Box<dyn Authenticator>,
        store: Box<dyn TokenStore<'a>>,
        token_ttl: u64,
        secret: &str,
    ) -> Result<AuthService<'a>, AuthenticationError> {
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
            authenticator,
            store,
            token_ttl,
            token_key,
        })
    }

    // Authenticates user against authentication service
    // TODO Make TokenStore store a String (serialized) representation of the token
    // pub fn login(&mut self, usr: &str, pass: &str) -> Result<String, AuthenticationError> {
    //     match self.authenticator.authenticate(usr, pass) {
    //         Ok(logged_in) => {
    //             if !logged_in {
    //                 return Err(AuthenticationError {
    //                     message: String::from("invalid credentials"),
    //                 });
    //             }
    //             let token: JWTToken;
    //             match JWTToken::generate(usr, self.token_ttl, &self.token_key) {
    //                 Ok(t) => token = t,
    //                 Err(e) => return Err(AuthenticationError { message: e.message }),
    //             }
    //             let result = token.value;
    //             match self.store.save(&usr, token) {
    //                 Ok(_) => return Ok(result),
    //                 Err(e) => return Err(AuthenticationError { message: e.message }),
    //             }
    //         }
    //         Err(e) => return Err(e),
    //     };
    // }

    // Checks whether the received token is still authorized
    pub fn authorize(&'a self, token_value: &str) -> Result<bool, AuthenticationError> {
        match JWTToken::get_username(token_value, &self.token_key) {
            Ok(usr) => match self.store.retrieve(&usr[..]) {
                Ok(saved_token) => match saved_token.is_evicted() {
                    Ok(evicted) => Ok(evicted),
                    Err(e) => return Err(AuthenticationError { message: e.message }),
                },
                Err(e) => return Err(AuthenticationError { message: e.message }),
            },
            Err(e) => return Err(AuthenticationError { message: e.message }),
        }
    }
}