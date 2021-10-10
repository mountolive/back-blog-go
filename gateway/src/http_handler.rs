//! HTTP handlers' definitions
use crate::auth::AuthService;
use crate::post::{CreatePost, PostCreator, PostUpdater, UpdatePost};
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

    if rej.find::<Unauthorized>().is_some() {
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

    if let Some(e) = rej.find::<warp::reject::InvalidQuery>() {
        match e.source() {
            Some(cause) => {
                code = StatusCode::BAD_REQUEST;
                message = cause.to_string();
            }
            None => {
                message = "unknown query serialization error".to_string();
            }
        }
    }

    if let Some(e) = rej.find::<HandlerError>() {
        message = format!("external error: {}", e.message);
    }

    let err = APIError {
        code: code.as_u16(),
        message: format!("{}: {:?}", message, rej),
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
    async fn _posts(&self, filter: crate::post::Filter) -> JSONResponse {
        match self.reader.posts(filter).await {
            Ok(posts) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&posts),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    async fn posts(
        &self,
        filter: crate::post::Filter,
    ) -> Result<impl warp::Reply, warp::Rejection> {
        Ok(self._posts(filter).await)
    }

    async fn _post(&self, id: &str) -> JSONResponse {
        match self.reader.post(id).await {
            Ok(post) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&post),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    async fn post(&self, id: String) -> Result<impl warp::Reply, warp::Rejection> {
        Ok(self._post(&id[..]).await)
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

    async fn login(&self, creds: LoginDTO) -> JSONResponse {
        match self
            .auth
            .login(&creds.username[..], &creds.password[..])
            .await
        {
            Ok(token) => JSONResponse(Ok(warp::reply::with_status(
                warp::reply::json(&token),
                StatusCode::OK,
            ))),
            Err(err) => JSONResponse(Err(HandlerError {
                message: err.message,
            })),
        }
    }

    async fn authenticate(&self, creds: LoginDTO) -> Result<impl warp::Reply, warp::Rejection> {
        Ok(self.login(creds).await)
    }

    fn authorize(&self) -> impl Filter<Extract = ((),), Error = Rejection> + Copy + '_ {
        warp::header::<String>("Authorization").and_then(move |token: String| async move {
            match self.auth.authorize(token.trim_start_matches(TOKEN_PREFIX)) {
                Ok(is_evicted) => {
                    if is_evicted {
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

    pub async fn start(&'static self, addr: SocketAddr) {
        let authenticate = warp::path!("user")
            .and(warp::post())
            .and(warp::body::json())
            .and_then(move |creds: LoginDTO| self.authenticate(creds));

        let post_by_id = warp::path!("posts" / String)
            .and(warp::get())
            .and_then(move |id: String| self.post(id));

        let posts_by_filter = warp::path!("posts")
            .and(warp::get())
            // TODO: post_by_filter should be able to work by means of queryParams
            .and(warp::body::json())
            .and_then(move |filter: crate::post::Filter| self.posts(filter));

        let create_post = warp::path!("posts")
            .and(warp::post())
            .and(self.authorize())
            .and(warp::body::json())
            .map(move |_, create: CreatePost| self.create_post(create));

        let update_post = warp::path!("posts" / String)
            .and(warp::put())
            .and(self.authorize())
            .and(warp::body::json())
            .map(move |id: String, _, update: UpdatePost| self.update_post(id, update));

        let cors = warp::cors().allow_any_origin().allow_methods(vec!["GET"]);

        let routes = post_by_id
            .or(authenticate)
            .or(posts_by_filter)
            .or(create_post)
            .or(update_post)
            .recover(error_handler)
            .with(cors);

        warp::serve(routes).run(addr).await;
    }
}
