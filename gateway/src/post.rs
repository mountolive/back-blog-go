//! Standard operations related to posts

use std::fmt;
use std::time::Instant;

// DTO a full post's data
pub struct Post {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
    pub created_at: Instant,
}

// Error associated to a create or update action regarding posts
#[derive(Debug)]
pub struct MutatorError {
    message: String,
}

impl std::error::Error for MutatorError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for MutatorError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "post mutation error: {}", self.message)
    }
}

// Defines the contract for creating or updating posts
pub trait Mutate {
    fn create(&self, post: Post) -> Result<(), MutatorError>;
    fn update(&self, post: Post) -> Result<(), MutatorError>;
}

// Error associated to a create or update action regarding posts
#[derive(Debug)]
pub struct ReaderError {
    message: String,
}

impl std::error::Error for ReaderError {
    fn description(&self) -> &str {
        &self.message[..]
    }
}

impl fmt::Display for ReaderError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "post reading error: {}", self.message)
    }
}

// Filter contains the possible options for filtering posts
enum Filter {
    DateRange {
        from: Instant,
        to: Instant,
        page: i32,
        page_size: i32,
    },
    Tags {
        tags: Vec<String>,
        page: i32,
        page_size: i32,
    },
}

// DTO holding basic data from a post
pub struct PostSummary {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub tags: Vec<String>,
    pub created_at: Instant,
}

// Read defines the contract for listing posts
pub trait Read {
    fn posts(&self, filter: Filter) -> Result<Vec<PostSummary>, ReaderError>;
    fn post(&self, id: &str) -> Result<Post, ReaderError>;
}
