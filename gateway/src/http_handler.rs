use crate::auth::AuthService;
use crate::post::{Filter as PostFilter, PostCreator, PostUpdater, ReadClient};
use crate::post_reader::PostReader;
use std::net::SocketAddr;
use warp::Filter;

/// HTTPHandler wraps all needed services and defines the registered routes
pub struct HTTPHandler {
    pub auth: AuthService,
    pub creator: PostCreator,
    pub updater: PostUpdater,
    pub reader: PostReader,
}

impl HTTPHandler {
    /// Starts the server
    pub async fn start(&'static self, addr: SocketAddr) {
        let posts_by_filter = warp::path!("posts")
            .and(warp::get())
            .and(warp::query().map(move |filter: PostFilter| self.reader.posts(filter)));

        let post_by_id = warp::path!("posts" / String).and(warp::get());

        let create_post = warp::path!("posts").and(warp::post());

        let update_post = warp::path!("posts" / String).and(warp::put());

        let filters = posts_by_filter
            .or(post_by_id)
            .or(create_post)
            .or(update_post)
            .map(|_| "implement");

        warp::serve(filters).run(addr).await;
    }
}
