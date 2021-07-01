//! Implementation of the corresponding TokenStore

use crate::auth::{JWTToken, TokenStore, TokenStoreError};

// Error placeholder for any error occurred inside a Storage's operation
pub trait StorageError: std::error::Error {}

// Represents the basic contract for a token's storage
pub trait Storage<'a> {
    fn get(&'a self, key: &str) -> Result<&'a JWTToken, Box<dyn StorageError>>;
    fn set(&mut self, key: &str, token: JWTToken) -> Result<(), Box<dyn StorageError>>;
}

// Wrapper for any tokens' storage engine
pub struct JWTStore<'a> {
    storage: dyn Storage<'a>,
}

impl<'a> TokenStore<'a> for JWTStore<'a> {
    // Retrieves a token from the underlying storage
    fn retrieve(&'a self, key: &str) -> Result<&'a JWTToken, TokenStoreError> {
        match self.storage.get(key) {
            Ok(token) => Ok(token),
            Err(e) => Err(TokenStoreError {
                message: e.to_string(),
            }),
        }
    }

    // Saves a token in the underlying storage
    fn save(&mut self, key: &str, token: JWTToken) -> Result<(), TokenStoreError> {
        match self.storage.set(key, token) {
            Ok(_) => Ok(()),
            Err(e) => Err(TokenStoreError {
                message: e.to_string(),
            }),
        }
    }
}
