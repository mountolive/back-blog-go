//! In memory implementation for token's storage

use crate::auth::JWTToken;
use crate::store::{Storage, StorageError};
use std::collections::HashMap;
use std::fmt;

#[derive(Debug)]
pub struct MemStorageError {
    message: String,
}

// An Error implementation for memory's storage
impl StorageError for MemStorageError {}

impl std::error::Error for MemStorageError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for MemStorageError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "mem storage error: {}", self.message)
    }
}

// Represents a token's storage in memory by means of a hash map
pub struct MemStorage {
    data: HashMap<String, JWTToken>,
}

impl<'a> MemStorage {
    // Creates a new storage
    fn new() -> MemStorage {
        MemStorage {
            data: HashMap::new(),
        }
    }
}

impl<'a> Storage<'a> for MemStorage {
    // Retrieves the corresponding token associated with the passed key or returns a MemStorageError
    fn get(&'a self, key: &str) -> Result<&'a JWTToken, Box<dyn StorageError>> {
        match self.data.get(&key.to_string()) {
            Some(token) => Ok(token),
            None => Err(Box::new(MemStorageError {
                message: String::from("no token found for the passed key"),
            })),
        }
    }

    // Saves the corresponding token associated with the passed key or returns a MemStorageError
    fn set(&mut self, key: &str, value: JWTToken) -> Result<(), Box<dyn StorageError>> {
        self.data.insert(String::from(key), value);
        Ok(())
    }
}
