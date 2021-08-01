use crate::auth::AuthService;
use crate::post::{PostCreator, PostUpdater};
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
    pub async fn start(&self, addr: SocketAddr) {
        let posts_by_filter = warp::path!("posts").and(warp::get());

        let post_by_id = warp::path!("posts" / String).and(warp::get());

        let filters = posts_by_filter.or(post_by_id).map(|_| "implement");

        warp::serve(filters).run(addr).await;
    }
}
