pub struct Token {
    pub value: String,
    pub ttl: i32,
}

pub trait StoreError: std::error::Error {}

pub trait Store {
    fn get(&self, key: &str) -> Result<Token, Box<dyn StoreError>>;
    fn set(&self, key: &str, token: &Token) -> Result<(), Box<dyn StoreError>>;
}
