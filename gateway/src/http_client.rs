use crate::post::{Filter, Post, PostSummary, ReadClient, ReaderError};
use reqwest;
use std::time::SystemTime;

pub struct PostReader {}

impl PostReader {
    fn build_filter_str(&self, filter: Filter) -> String {
        // TODO Implement
        "".to_string()
    }
}

impl ReadClient for PostReader {
    fn posts(&self, filter: Filter) -> Result<Vec<PostSummary>, ReaderError> {
        match reqwest::blocking::get(self.build_filter_str(filter)) {
            // TODO Handle
        };
    }

    fn post(&self, id: &str) -> Result<Post, ReaderError> {
        // TODO Implement
        Ok(Post {
            id: "some_id".to_string(),
            creator: "some_c".to_string(),
            content: "con".to_string(),
            title: "title".to_string(),
            tags: vec!["bla".to_string(), "bla".to_string()],
            created_at: SystemTime::now(),
        })
    }
}
