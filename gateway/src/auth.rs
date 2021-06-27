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

impl JWTToken {
    // Creates a new JWT token from the username passed and with the passed TTL (in seconds)
    pub fn generate(
        username: &String,
        ttl: u64,
        key: &Hmac<Sha256>,
    ) -> Result<JWTToken, TokenError> {
        let mut claims = BTreeMap::new();
        claims.insert(USER_KEY, username);
        let value: String;
        match claims.sign_with_key(key) {
            Ok(token_val) => value = token_val,
            Err(_) => {
                return Err(TokenError {
                    message: String::from("signing token"),
                })
            }
        }
        let until: Duration;
        match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
            Ok(now) => until = now + Duration::from_secs(ttl),
            Err(_) => {
                return Err(TokenError {
                    message: String::from(WRONG_NOW),
                })
            }
        }
        Ok(JWTToken { value, until })
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

    // Explodes and a gets underlying username from JWT token
    fn get_username(&self, key: &Hmac<Sha256>) -> Result<String, TokenError> {
        let token: Token<Header, BTreeMap<String, String>, _>;
        match VerifyWithKey::verify_with_key(&self.value[..], key) {
            Ok(val) => {
                token = val;
                Ok(String::from(&token.claims()[USER_KEY]))
            }
            Err(_) => Err(TokenError {
                message: String::from("unable to retrieve username"),
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
pub trait TokenStore {
    fn retrieve(&self, key: &String) -> Result<JWTToken, TokenStoreError>;
    fn save(&self, key: &String, token: &JWTToken) -> Result<(), TokenStoreError>;
}

// Describes the basic contract expected by an authenticator's client
pub trait Authenticator {
    fn authenticate(
        &self,
        username: &String,
        password: &String,
    ) -> Result<bool, AuthenticationError>;
}

// Handles login and authorization by means of an Authenticator and a TokenStore
pub struct AuthService<T: Authenticator, V: TokenStore> {
    authenticator: T,
    store: V,
    token_ttl: u64,
    token_key: Hmac<Sha256>,
}

impl<T: Authenticator, V: TokenStore> AuthService<T, V> {
    // Creates a new AuthService
    pub fn new(
        authenticator: T,
        store: V,
        token_ttl: u64,
        secret: String,
    ) -> Result<AuthService<T, V>, AuthenticationError> {
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
    pub fn login(&self, usr: String, pass: String) -> Result<JWTToken, AuthenticationError> {
        match self.authenticator.authenticate(&usr, &pass) {
            Ok(logged_in) => {
                if !logged_in {
                    return Err(AuthenticationError {
                        message: String::from("invalid credentials"),
                    });
                }
                let token: JWTToken;
                match JWTToken::generate(&usr, self.token_ttl, &self.token_key) {
                    Ok(t) => token = t,
                    Err(e) => return Err(AuthenticationError { message: e.message }),
                }
                match self.store.save(&usr, &token) {
                    Ok(_) => return Ok(token),
                    Err(e) => return Err(AuthenticationError { message: e.message }),
                }
            }
            Err(e) => return Err(e),
        };
    }

    // Checks whether the received token is still authorized
    pub fn authorize(&self, token: &JWTToken) -> Result<bool, AuthenticationError> {
        let complete_token: JWTToken;
        match token.get_username(&self.token_key) {
            Ok(usr) => match self.store.retrieve(&usr) {
                Ok(saved_token) => {
                    if saved_token != *token {
                        return Ok(false);
                    }
                    complete_token = saved_token;
                }
                Err(e) => return Err(AuthenticationError { message: e.message }),
            },
            Err(e) => return Err(AuthenticationError { message: e.message }),
        }
        match complete_token.is_evicted() {
            Ok(evicted) => Ok(evicted),
            Err(e) => return Err(AuthenticationError { message: e.message }),
        }
    }
}
