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

pub trait Token {
    fn is_evicted(&self) -> Result<bool, EvictionCheckError>;
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
    fn retrieve(&self, key: &String) -> Result<Box<dyn Token>, TokenStoreError>;
    fn save(&self, key: &String, token: &dyn Token) -> Result<(), TokenStoreError>;
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
}

impl<T: Authenticator, V: TokenStore> AuthService<T, V> {
    pub fn new(authenticator: T, store: V) -> AuthService<T, V> {
        AuthService {
            authenticator,
            store,
        }
    }

    pub fn login<S: Token>(&self, usr: String, pass: String) -> Result<S, AuthenticationError> {
        match self.authenticator.authenticate(&usr, &pass) {
            Ok(token) => match self.store.save(usr, Box::new(token)) {
                Ok(correct) => {
                    if correct {
                        return Ok(token);
                    }
                    return Err(AuthenticationError { message: String::from("token unable to be saved") }
                },
                Err(e) => return Err(AuthenticationError { message: e.message }),
            },
            Err(e) => return Err(e),
        };
    }
}
