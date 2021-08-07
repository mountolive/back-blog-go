//! Standard operations related to posts

use serde::{Deserialize, Serialize};
use std::fmt;
use std::time::SystemTime;

/// DTO a full post's data
#[derive(Deserialize)]
pub struct Post {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
    pub created_at: SystemTime,
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
#[derive(Serialize, Deserialize)]
pub struct CreatePost {
    pub creator: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
}

/// DTO with the data needed for updating a post
#[derive(Deserialize)]
pub struct UpdatePost {
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

// Marking it as "thread-safe"
unsafe impl Send for PostCreator {}
unsafe impl Sync for PostCreator {}

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

/// DTO hold all the required data to update a post against the remote post service
#[derive(Serialize)]
pub struct FullUpdatePost {
    pub id: String,
    pub title: String,
    pub content: String,
    pub tags: Vec<String>,
}

/// A service for updating a post
pub struct PostUpdater {
    pub client: Box<dyn MutatorClient<FullUpdatePost>>,
}

// Marking it as "thread-safe"
unsafe impl Send for PostUpdater {}
unsafe impl Sync for PostUpdater {}

impl PostUpdater {
    /// Updates a post with the corresponding data passed
    pub fn update(&self, id: &str, post: UpdatePost) -> Result<(), MutatorError> {
        let update_post = FullUpdatePost {
            id: id.to_string(),
            title: post.title,
            content: post.content,
            tags: post.tags,
        };
        match self.client.send(update_post) {
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
    pub message: String,
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
#[derive(Deserialize)]
pub enum Filter {
    DateRange {
        from: SystemTime,
        to: SystemTime,
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
#[derive(Deserialize)]
pub struct PostSummary {
    pub id: String,
    pub creator: String,
    pub title: String,
    pub tags: Vec<String>,
    pub created_at: SystemTime,
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

    impl MutatorClient<FullUpdatePost> for MockClient {
        fn send(&self, _: FullUpdatePost) -> Result<(), MutatorError> {
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
        match creator.update("some-id", update_post()) {
            Ok(()) => assert!(false, "shouldn't return Ok"),
            Err(_) => assert!(true, "should return Err"),
        }
    }

    #[test]
    fn test_correct_client_post_updater_update() {
        let creator = PostUpdater {
            client: Box::new(MockClient { errored: false }),
        };
        match creator.update("other-id", update_post()) {
            Ok(()) => assert!(true, "should return Ok"),
            Err(e) => assert!(false, "shouldn't return Err {}", e),
        }
    }
}
