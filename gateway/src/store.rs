//! Implementation of the corresponding TokenStore

use crate::auth::{JWTToken, TokenStore, TokenStoreError};

/// Error placeholder for any error occurred inside a Storage's operation
pub trait StorageError: std::error::Error {}

/// Represents the basic contract for a storage's driver
pub trait StorageDriver {
    fn get(&self, key: &str) -> Result<JWTToken, Box<dyn StorageError>>;
    fn set(&self, key: &str, ser_token: &str) -> Result<(), Box<dyn StorageError>>;
}

/// Wrapper for any tokens' storage engine
pub struct JWTStore {
    pub storage: Box<dyn StorageDriver>,
}

impl TokenStore for JWTStore {
    /// Retrieves a token from the underlying storage
    fn retrieve(&self, key: &str) -> Result<JWTToken, TokenStoreError> {
        match self.storage.get(key) {
            Ok(token) => Ok(token),
            Err(e) => Err(TokenStoreError {
                message: e.to_string(),
            }),
        }
    }

    /// Saves a token in the underlying storage
    fn save(&self, key: &str, ser_token: &str) -> Result<(), TokenStoreError> {
        match self.storage.set(key, ser_token) {
            Ok(_) => Ok(()),
            Err(e) => Err(TokenStoreError {
                message: e.to_string(),
            }),
        }
    }
}
