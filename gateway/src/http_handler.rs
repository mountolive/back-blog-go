use crate::auth::AuthService;
use crate::post::{
    CreatePost, Filter as PostFilter, PostCreator, PostUpdater, ReadClient, UpdatePost,
};
use crate::post_reader::PostReader;
use serde::Serialize;
use std::convert::Infallible;
use std::error::Error;
use std::net::SocketAddr;
use warp::http::StatusCode;
use warp::{reject, Filter, Rejection, Reply};

/// HTTPHandler wraps all needed services and defines the registered routes
pub struct HTTPHandler {
    pub auth: AuthService,
    pub creator: PostCreator,
    pub updater: PostUpdater,
    pub reader: PostReader,
}

/// APIError is self-described
#[derive(Serialize)]
struct APIError {
    code: u16,
    message: String,
}

/// Fall-thru function to handle rejections
async fn error_handler(rej: Rejection) -> Result<impl Reply, Infallible> {
    let mut code = StatusCode::INTERNAL_SERVER_ERROR;
    let mut err = APIError {
        code: code.as_u16(),
        message: "unexpected error".to_string(),
    };

    if rej.is_not_found() {
        code = StatusCode::NOT_FOUND;
        err = APIError {
            code: code.as_u16(),
            message: "not found".to_string(),
        };
        return Ok(warp::reply::with_status(warp::reply::json(&err), code));
    }

    if let Some(e) = rej.find::<warp::filters::body::BodyDeserializeError>() {
        match e.source() {
            Some(cause) => {
                err = APIError {
                    code: StatusCode::BAD_REQUEST.as_u16(),
                    message: cause.to_string(),
                };
            }
            None => {
                err.message = "unknown serialization error".to_string();
            }
        }
        return Ok(warp::reply::with_status(warp::reply::json(&err), code));
    }

    if let Some(e) = rej.find::<HandlerError>() {
        err.message = format!("external error: {}", e.message);
        return Ok(warp::reply::with_status(warp::reply::json(&err), code));
    }

    Ok(warp::reply::with_status(warp::reply::json(&err), code))
}

#[derive(Debug)]
struct HandlerError {
    message: String,
}

impl reject::Reject for HandlerError {}

impl warp::Reply for HandlerError {
    fn into_response(self) -> warp::reply::Response {
        todo!()
    }
}

impl HTTPHandler {
    async fn posts(&self, filter: PostFilter) -> Result<impl warp::Reply, HandlerError> {
        match self.reader.posts(filter) {
            Ok(posts) => Ok(warp::reply::json(&posts)),
            Err(err) => Err(HandlerError {
                message: err.message,
            }),
        }
    }

    async fn post(&self, id: &str) -> Result<impl warp::Reply, HandlerError> {
        match self.reader.post(&id[..]) {
            Ok(post) => Ok(warp::reply::json(&post)),
            Err(err) => Err(HandlerError {
                message: err.message,
            }),
        }
    }

    async fn create_post(&self, create: CreatePost) -> Result<impl warp::Reply, HandlerError> {
        match self.creator.create(create) {
            Ok(_) => Ok(warp::reply::with_status("OK", StatusCode::CREATED)),
            Err(err) => Err(HandlerError {
                message: err.message,
            }),
        }
    }

    async fn update_post(
        &self,
        id: String,
        update: UpdatePost,
    ) -> Result<impl warp::Reply, HandlerError> {
        match self.updater.update(&id[..], update) {
            Ok(_) => Ok(warp::reply::with_status("OK", StatusCode::NO_CONTENT)),
            Err(err) => Err(HandlerError {
                message: err.message,
            }),
        }
    }

    /// Starts the server
    pub async fn start(&'static self, addr: SocketAddr) {
        let posts_by_filter = warp::path!("posts")
            .and(warp::get())
            .and(warp::query().map(move |filter: PostFilter| self.posts(filter)));

        let post_by_id = warp::path!("posts")
            .and(warp::get())
            .and(warp::path::param())
            .map(move |id: String| self.post(&id[..]));

        let create_post = warp::path!("posts")
            .and(warp::post())
            .and(warp::body::json())
            .map(move |create: CreatePost| self.create_post(create));

        let update_post = warp::path!("posts")
            .and(warp::put())
            .and(warp::path::param::<String>())
            .and(warp::body::json())
            .map(move |id: String, update: UpdatePost| self.update_post(id, update));

        let filters = posts_by_filter
            .or(post_by_id)
            .or(create_post)
            .or(update_post)
            .recover(error_handler);

        warp::serve(filters).run(addr).await;
    }
}
