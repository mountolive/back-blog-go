use crate::auth::{JWTToken, TokenStore, TokenStoreError};

pub trait StorageError: std::error::Error {}

pub trait Storage {
    fn get(&self, key: &String) -> Result<JWTToken, Box<dyn StorageError>>;
    fn set(&self, key: &String, token: &JWTToken) -> Result<(), Box<dyn StorageError>>;
}

pub struct JWTStore {
    store: dyn Storage,
}

impl TokenStore for JWTStore {
    fn retrieve(&self, key: &String) -> Result<JWTToken, TokenStoreError> {
        // TODO Implement
        Ok(JWTToken {
            value: String::from("TODO"),
        })
    }

    fn save(&self, key: &String, token: &JWTToken) -> Result<(), TokenStoreError> {
        // TODO Implement
        Ok(())
    }
}
