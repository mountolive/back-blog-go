use crate::auth::{EvictionCheckError, Token, TokenStore, TokenStoreError};

pub struct JWTToken {
    pub value: String,
    // TODO Add ttl to JWTToken
}

impl Token for JWTToken {
    fn is_evicted(&self) -> Result<bool, EvictionCheckError> {
        // TODO Implement
        Ok(false)
    }
}

impl JWTToken {
    pub fn generate(username: &String, ttl: i32) -> JWTToken {
        // TODO Implement
        JWTToken {
            value: String::from("TODO"),
        }
    }
}

pub trait StorageError: std::error::Error {}

pub trait Storage {
    fn get(&self, key: &String) -> Result<JWTToken, Box<dyn StorageError>>;
    fn set(&self, key: &String, token: &JWTToken) -> Result<(), Box<dyn StorageError>>;
}

pub struct JWTStore {
    store: dyn Storage,
}

impl TokenStore for JWTStore {
    fn retrieve(&self, key: &String) -> Result<Box<dyn Token>, TokenStoreError> {
        // TODO Implement
        Ok(Box::new(JWTToken {
            value: String::from("TODO"),
        }))
    }

    fn save(&self, key: &String, token: &dyn Token) -> Result<(), TokenStoreError> {
        // TODO Implement
        Ok(())
    }
}
