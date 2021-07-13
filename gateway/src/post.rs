//! Standard operations related to posts

use serde::Serialize;
use std::fmt;
use std::time::Instant;

/// DTO a full post's data
pub struct Post {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
    pub created_at: Instant,
}

/// Error associated to a create or update action regarding posts
#[derive(Debug)]
pub struct MutatorError {
    pub message: String,
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

/// DTO with the data needed for creating a post
#[derive(Serialize)]
pub struct CreatePost {
    pub creator: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
}

/// DTO with the data needed for updating a post
#[derive(Serialize)]
pub struct UpdatePost {
    pub id: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
}

/// Defines the mutator client
pub trait MutatorClient<T: Serialize> {
    fn send(&self, payload: T) -> Result<(), MutatorError>;
}

/// A service for creating posts
pub struct PostCreator {
    pub client: Box<dyn MutatorClient<CreatePost>>,
}

impl PostCreator {
    /// Creates a post with the corresponding data passed
    pub fn create(&self, post: CreatePost) -> Result<(), MutatorError> {
        match self.client.send(post) {
            Ok(()) => Ok(()),
            Err(e) => Err(MutatorError {
                message: format!("post creator: {}", e.message),
            }),
        }
    }
}

/// A service for updating a post
pub struct PostUpdater {
    pub client: Box<dyn MutatorClient<UpdatePost>>,
}

impl PostUpdater {
    /// Updates a post with the corresponding data passed
    pub fn update(&self, post: UpdatePost) -> Result<(), MutatorError> {
        match self.client.send(post) {
            Ok(()) => Ok(()),
            Err(e) => Err(MutatorError {
                message: format!("post update: {}", e.message),
            }),
        }
    }
}

/// Error associated to a create or update action regarding posts
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

/// Filter contains the possible options for filtering posts
pub enum Filter {
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

/// DTO holding basic data from a post
pub struct PostSummary {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub tags: Vec<String>,
    pub created_at: Instant,
}

/// Read defines the contract for listing posts
pub trait ReadClient {
    fn posts(&self, filter: Filter) -> Result<Vec<PostSummary>, ReaderError>;
    fn post(&self, id: &str) -> Result<Post, ReaderError>;
}

mod test {
    use super::*;

    struct MockClient {
        errored: bool,
    }

    impl MutatorClient<CreatePost> for MockClient {
        fn send(&self, _: CreatePost) -> Result<(), MutatorError> {
            if self.errored {
                return Err(MutatorError {
                    message: String::from("whatever error"),
                });
            }
            Ok(())
        }
    }

    fn create_post() -> CreatePost {
        CreatePost {
            creator: String::from("some-creator"),
            title: String::from("some-title"),
            content: String::from("some-content"),
            tags: vec![String::from("cool"), String::from("awesome")],
        }
    }

    #[test]
    fn test_errored_client_post_creator_create() {
        let creator = PostCreator {
            client: Box::new(MockClient { errored: true }),
        };
        match creator.create(create_post()) {
            Ok(()) => assert!(false, "shouldn't return Ok"),
            Err(_) => assert!(true, "should return Err"),
        }
    }

    #[test]
    fn test_correct_client_post_creator_create() {
        let creator = PostCreator {
            client: Box::new(MockClient { errored: false }),
        };
        match creator.create(create_post()) {
            Ok(()) => assert!(true, "should return Ok"),
            Err(e) => assert!(false, "shouldn't return Err {}", e),
        }
    }

    impl MutatorClient<UpdatePost> for MockClient {
        fn send(&self, _: UpdatePost) -> Result<(), MutatorError> {
            if self.errored {
                return Err(MutatorError {
                    message: String::from("another whatever error"),
                });
            }
            Ok(())
        }
    }

    fn update_post() -> UpdatePost {
        UpdatePost {
            id: String::from("some-id"),
            title: String::from("some-title"),
            content: String::from("some-content"),
            tags: vec![String::from("amazing"), String::from("superb")],
        }
    }

    #[test]
    fn test_errored_client_post_updater_update() {
        let creator = PostUpdater {
            client: Box::new(MockClient { errored: true }),
        };
        match creator.update(update_post()) {
            Ok(()) => assert!(false, "shouldn't return Ok"),
            Err(_) => assert!(true, "should return Err"),
        }
    }

    #[test]
    fn test_correct_client_post_updater_update() {
        let creator = PostUpdater {
            client: Box::new(MockClient { errored: false }),
        };
        match creator.update(update_post()) {
            Ok(()) => assert!(true, "should return Ok"),
            Err(e) => assert!(false, "shouldn't return Err {}", e),
        }
    }
}
