use std::fmt;

#[derive(Debug)]
pub struct EvictionCheckError {
    message: String,
}

impl std::error::Error for EvictionCheckError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for EvictionCheckError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "eviction check error: {}", self.message)
    }
}

pub struct JWTToken {
    pub value: String,
    // TODO Add ttl to JWTToken
}

impl JWTToken {
    pub fn generate(username: &String, ttl: i32) -> JWTToken {
        // TODO Implement
        JWTToken {
            value: String::from("TODO"),
        }
    }

    fn is_evicted(&self) -> Result<bool, EvictionCheckError> {
        // TODO Implement
        Ok(false)
    }
}

#[derive(Debug)]
pub struct TokenStoreError {
    message: String,
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

pub trait TokenStore {
    fn retrieve(&self, key: &String) -> Result<JWTToken, TokenStoreError>;
    fn save(&self, key: &String, token: &JWTToken) -> Result<(), TokenStoreError>;
}

pub trait Authenticator {
    fn authenticate(
        &self,
        username: &String,
        password: &String,
    ) -> Result<bool, AuthenticationError>;
}

pub struct AuthService<T: Authenticator, V: TokenStore> {
    authenticator: T,
    store: V,
    token_ttl: i32,
}

impl<T: Authenticator, V: TokenStore> AuthService<T, V> {
    pub fn new(authenticator: T, store: V, token_ttl: i32) -> AuthService<T, V> {
        AuthService {
            authenticator,
            store,
            token_ttl,
        }
    }

    pub fn login(&self, usr: String, pass: String) -> Result<JWTToken, AuthenticationError> {
        match self.authenticator.authenticate(&usr, &pass) {
            Ok(logged_in) => {
                if !logged_in {
                    return Err(AuthenticationError {
                        message: String::from("invalid credentials"),
                    });
                }
                let token = JWTToken::generate(&usr, self.token_ttl);
                match self.store.save(&usr, &token) {
                    Ok(_) => return Ok(token),
                    Err(e) => {
                        return Err(AuthenticationError {
                            message: String::from("implement"),
                        })
                    }
                }
            }
            Err(e) => return Err(e),
        };
    }
}
