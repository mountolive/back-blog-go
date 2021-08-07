use crate::auth::AuthService;
use crate::post::{
    CreatePost, Filter as PostFilter, PostCreator, PostUpdater, ReadClient, UpdatePost,
};
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
        let posts_by_filter = warp::path!("posts").and(warp::get()).and(warp::query().map(
            move |filter: PostFilter| {
                self.reader.posts(filter);
                // TODO: Write response, post_by_filter, gateway
            },
        ));

        let post_by_id = warp::path!("posts")
            .and(warp::get())
            .and(warp::path::param())
            .map(move |id: String| {
                self.reader.post(&id[..]);
                // TODO: Write response, post_by_id, gateway
            });

        let create_post = warp::path!("posts")
            .and(warp::post())
            .and(warp::body::json())
            .map(move |create: CreatePost| {
                self.creator.create(create);
                // TODO: Write response, create_post, gateway
            });

        let update_post = warp::path!("posts")
            .and(warp::put())
            .and(warp::path::param::<String>())
            .and(warp::body::json())
            .map(move |id: String, update: UpdatePost| {
                self.updater.update(&id[..], update);
                // TODO: Write response, update_post
            });

        let filters = posts_by_filter
            .or(post_by_id)
            .or(create_post)
            .or(update_post)
            .map(|_| "implement");

        warp::serve(filters).run(addr).await;
    }
}
