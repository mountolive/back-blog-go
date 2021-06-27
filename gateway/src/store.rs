use crate::auth::{JWTToken, TokenStore, TokenStoreError};

pub trait StorageError: std::error::Error {}

pub trait Storage {
    fn get(&self, key: &String) -> Result<String, Box<dyn StorageError>>;
    fn set(&self, key: &String, token: &JWTToken) -> Result<(), Box<dyn StorageError>>;
}

pub struct JWTStore<T: Storage> {
    store: T,
}

impl<T: Storage> TokenStore for JWTStore<T> {
    fn retrieve(&self, key: &String) -> Result<JWTToken, TokenStoreError> {
        // TODO Implement
        Err(TokenStoreError {
            message: String::from("implement"),
        })
    }

    fn save(&self, key: &String, token: &JWTToken) -> Result<(), TokenStoreError> {
        // TODO Implement
        Ok(())
    }
}
