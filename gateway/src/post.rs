//! Standard operations related to posts

// TODO Implement
pub struct Post;

// TODO Implement
pub struct MutatorError {}

pub trait Mutator {
    fn create(&self, post: Post) -> Result<(), MutatorError>;
    fn update(&self, post: Post) -> Result<(), MutatorError>;
}

// TODO Implement
pub struct ReaderError {}

// TODO Implement
enum Filter {}

pub struct PostSummary;

pub trait Reader {
    fn posts(&self, filter: Filter) -> Result<Vec<PostSummary>, ReaderError>;
    fn post(&self, id: &str) -> Result<Post, ReaderError>;
}
