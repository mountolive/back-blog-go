//! HTTP handlers' definitions
use crate::auth::AuthService;
use crate::post::{
    CreatePost, Filter as PostFilter, PostCreator, PostUpdater, ReadClient, UpdatePost,
};
use crate::post_reader::PostReader;
use serde::{Deserialize, Serialize};
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
    let mut message = "unexpected error".to_string();

    if rej.is_not_found() {
        code = StatusCode::NOT_FOUND;
        message = "not found".to_string();
    }

    if let Some(_) = rej.find::<Unauthorized>() {
        code = StatusCode::UNAUTHORIZED;
        message = "naughty".to_string();
    }

    if let Some(e) = rej.find::<warp::filters::body::BodyDeserializeError>() {
        match e.source() {
            Some(cause) => {
                code = StatusCode::BAD_REQUEST;
                message = cause.to_string();
            }
            None => {
                message = "unknown serialization error".to_string();
            }
        }
    }

    if let Some(e) = rej.find::<HandlerError>() {
        message = format!("external error: {}", e.message);
    }

    let err = APIError {
        code: code.as_u16(),
        message,
    };

    Ok(warp::reply::with_status(warp::reply::json(&err), code))
}

#[derive(Debug)]
struct HandlerError {
    message: String,
}

impl reject::Reject for HandlerError {}

impl warp::Reply for HandlerError {
    fn into_response(self) -> warp::reply::Response {
        // TODO Handle errors specifically for HandleError when converting into warp::reply::Response
        warp::reply::with_status(
            warp::reply::Response::new(warp::hyper::Body::from(self.message)),
            StatusCode::INTERNAL_SERVER_ERROR,
        )
        .into_response()
    }
}

struct JSONResponse(Result<warp::reply::WithStatus<warp::reply::Json>, HandlerError>);

impl warp::Reply for JSONResponse {
    fn into_response(self) -> warp::reply::Response {
        match self.0 {
            Ok(ok) => ok.into_response(),
            Err(err) => err.into_response(),
        }
    }
}

struct EmptyResponse(Result<warp::reply::WithStatus<String>, HandlerError>);

impl warp::Reply for EmptyResponse {
    fn into_response(self) -> warp::reply::Response {
        match self.0 {
            Ok(ok) => ok.into_response(),
            Err(err) => err.into_response(),
        }
    }
}

#[derive(Deserialize)]
struct LoginDTO {
    username: String,
    password: String,
}

const TOKEN_PREFIX: &str = "Bearer ";

#[derive(Debug)]
struct Unauthorized;

impl reject::Reject for Unauthorized {}

impl HTTPHandler {
    fn posts(&self, filter: PostFilter) -> JSONResponse {
        match self.reader.posts(filter) {
            Ok(posts) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&posts),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    fn post(&self, id: &str) -> JSONResponse {
        match self.reader.post(&id[..]) {
            Ok(post) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&post),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    fn create_post(&self, create: CreatePost) -> EmptyResponse {
        match self.creator.create(create) {
            Ok(_) => EmptyResponse(Ok(warp::reply::with_status(
                "OK".to_string(),
                StatusCode::CREATED,
            ))),
            Err(err) => EmptyResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    fn update_post(&self, id: String, update: UpdatePost) -> EmptyResponse {
        match self.updater.update(&id[..], update) {
            Ok(_) => EmptyResponse(Ok(warp::reply::with_status(
                "OK".to_string(),
                StatusCode::NO_CONTENT,
            ))),
            Err(err) => EmptyResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    fn login(&self, usr: &str, pass: &str) -> JSONResponse {
        match self.auth.login(usr, pass) {
            Ok(token) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&token),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    fn authorize(&'static self) -> impl Filter<Extract = ((),), Error = Rejection> + Copy {
        warp::header::<String>("Authorization").and_then(move |token: String| async move {
            match self.auth.authorize(token.trim_start_matches(TOKEN_PREFIX)) {
                Ok(auth) => {
                    if !auth {
                        return Err(reject::custom(Unauthorized));
                    }
                    Ok(())
                }
                Err(err) => Err(reject::custom(HandlerError {
                    message: err.message,
                })),
            }
        })
    }

    /// Starts the server
    pub async fn start(&'static self, addr: SocketAddr) {
        let login = warp::path!("user")
            .and(warp::post())
            .and(warp::body::json())
            .map(move |creds: LoginDTO| self.login(&creds.username[..], &creds.password[..]));

        let posts_by_filter = warp::path!("posts")
            .and(warp::get())
            .and(warp::query().map(move |filter: PostFilter| self.posts(filter)));

        let post_by_id = warp::path!("posts")
            .and(warp::get())
            .and(warp::path::param())
            .map(move |id: String| self.post(&id[..]));

        let create_post = warp::path!("posts")
            .and(self.authorize())
            .and(warp::post())
            .and(warp::body::json())
            .map(move |_: (), create: CreatePost| self.create_post(create));

        let update_post = warp::path!("posts")
            .and(self.authorize())
            .and(warp::put())
            .and(warp::path::param::<String>())
            .and(warp::body::json())
            .map(move |_: (), id: String, update: UpdatePost| self.update_post(id, update));

        let routes = login
            .or(posts_by_filter)
            .or(post_by_id)
            .or(create_post)
            .or(update_post)
            .recover(error_handler);

        warp::serve(routes).run(addr).await;
    }
}
