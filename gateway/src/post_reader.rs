use crate::post::{Filter, Post, PostSummary, ReaderError};
use serde::Deserialize;

/// Wraps the Client's basic config
pub struct ReaderClientConfig {
    pub base_url: String,
    pub from_param_name: String,
    pub to_param_name: String,
    pub tag_param_name: String,
    pub page_param_name: String,
    pub page_size_param_name: String,
}

impl ReaderClientConfig {
    /// Builds a ReaderClientConfig with default parameter names and formatting for the passed url
    pub fn with_default(base_url: String) -> Self {
        ReaderClientConfig {
            base_url,
            from_param_name: "start_date".to_string(),
            to_param_name: "end_date".to_string(),
            tag_param_name: "tag".to_string(),
            page_param_name: "page".to_string(),
            page_size_param_name: "page_size".to_string(),
        }
    }
}

/// It's an impl ReaderClient
pub struct PostReader {
    client: reqwest::Client,
    parsed_base_url: reqwest::Url,
    config: ReaderClientConfig,
}

/// Defines a range filtering for posts retrieval
#[derive(Deserialize)]
pub struct DateRange {
    pub from: String,
    pub to: String,
    pub page: i32,
    pub page_size: i32,
}

/// Defines a tag's filtering for posts retrieval
#[derive(Deserialize)]
pub struct Tag {
    pub tag: String,
    pub page: i32,
    pub page_size: i32,
}

impl PostReader {
    /// Builds a PostReader with the passed configuration
    pub fn new(config: ReaderClientConfig) -> Result<Self, ReaderError> {
        match reqwest::Url::parse(&config.base_url[..]) {
            Ok(parsed_base_url) => Ok(PostReader {
                client: reqwest::Client::new(),
                parsed_base_url,
                config,
            }),
            Err(e) => Err(ReaderError {
                message: format!("parsing base_url: {}", e),
            }),
        }
    }

    fn date_filter(&self, filter: DateRange) -> String {
        let base_url = format!("{}posts", self.parsed_base_url.as_str());
        let url = reqwest::Url::parse_with_params(
            &base_url[..],
            &[
                (&self.config.from_param_name[..], &filter.from[..]),
                (&self.config.to_param_name[..], &filter.to[..]),
                (
                    &self.config.page_param_name[..],
                    &format!("{}", filter.page)[..],
                ),
                (
                    &self.config.page_size_param_name[..],
                    &format!("{}", filter.page_size)[..],
                ),
            ],
        )
        .unwrap();
        url.as_str().to_string()
    }

    fn tag_filter(&self, filter: Tag) -> String {
        let base_url = format!("{}posts", self.parsed_base_url.as_str());
        let url = reqwest::Url::parse_with_params(
            &base_url[..],
            &[
                (&self.config.tag_param_name[..], &filter.tag[..]),
                (
                    &self.config.page_param_name[..],
                    &format!("{}", filter.page)[..],
                ),
                (
                    &self.config.page_size_param_name[..],
                    &format!("{}", filter.page_size)[..],
                ),
            ],
        )
        .unwrap();
        url.as_str().to_string()
    }

    fn build_id_url(&self, id: &str) -> String {
        format!("{}posts/{}", self.parsed_base_url.as_str(), id)
    }

    /// Retrieves posts by the passed DateRange
    pub async fn posts_by_date(&self, filter: DateRange) -> Result<Vec<PostSummary>, ReaderError> {
        match self.client.get(self.date_filter(filter)).send().await {
            Ok(response) => match response.json().await {
                Ok(result) => Ok(result),
                Err(e) => Err(ReaderError {
                    message: format!("posts desearializing: {}", e),
                }),
            },
            Err(e) => Err(ReaderError {
                message: e.to_string(),
            }),
        }
    }

    /// Retrieves posts by the passed Tag's filter
    pub async fn posts_by_tag(&self, filter: Tag) -> Result<Vec<PostSummary>, ReaderError> {
        match self.client.get(self.tag_filter(filter)).send().await {
            Ok(response) => match response.json().await {
                Ok(result) => Ok(result),
                Err(e) => Err(ReaderError {
                    message: format!("posts desearializing: {}", e),
                }),
            },
            Err(e) => Err(ReaderError {
                message: e.to_string(),
            }),
        }
    }

    /// Retrieves the post with the passed id
    pub async fn post(&self, id: &str) -> Result<Post, ReaderError> {
        match self.client.get(self.build_id_url(id)).send().await {
            Ok(response) => match response.json().await {
                Ok(result) => Ok(result),
                Err(e) => Err(ReaderError {
                    message: format!("post desearializing: {:?}", e),
                }),
            },
            Err(e) => Err(ReaderError {
                message: e.to_string(),
            }),
        }
    }
}
