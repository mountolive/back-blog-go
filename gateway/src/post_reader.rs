use crate::post::{Filter, Post, PostSummary, ReaderError};
use reqwest;

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

    fn build_filter_url(&self, filter: Filter) -> String {
        let base_url = format!("{}posts", self.parsed_base_url.as_str());
        match filter {
            Filter::DateRange {
                from,
                to,
                page,
                page_size,
            } => {
                let url = reqwest::Url::parse_with_params(
                    &base_url[..],
                    &[
                        (&self.config.from_param_name[..], &from[..]),
                        (&self.config.to_param_name[..], &to[..]),
                        (&self.config.page_param_name[..], &format!("{}", page)[..]),
                        (
                            &self.config.page_size_param_name[..],
                            &format!("{}", page_size)[..],
                        ),
                    ],
                )
                .unwrap();
                url.as_str().to_string()
            }
            Filter::Tags {
                tag,
                page,
                page_size,
            } => {
                let url = reqwest::Url::parse_with_params(
                    &base_url[..],
                    &[
                        (&self.config.tag_param_name[..], &tag[..]),
                        (&self.config.page_param_name[..], &format!("{}", page)[..]),
                        (
                            &self.config.page_size_param_name[..],
                            &format!("{}", page_size)[..],
                        ),
                    ],
                )
                .unwrap();
                url.as_str().to_string()
            }
        }
    }

    fn build_id_url(&self, id: &str) -> String {
        format!("{}posts/{}", self.parsed_base_url.as_str(), id)
    }

    /// Retrieves posts by the passed filters
    pub async fn posts(&self, filter: Filter) -> Result<Vec<PostSummary>, ReaderError> {
        match self.client.get(self.build_filter_url(filter)).send().await {
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
