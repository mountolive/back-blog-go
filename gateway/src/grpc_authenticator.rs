// ! gRPC authenticator's implementation

use crate::auth::{AuthenticationError, Authenticator};
use crate::user::login_client::LoginClient;
use crate::user::{LoginRequest, LoginResponse};
use tokio;

/// GRPCAuthenticator gRPC authenticator implementation
pub struct GRPCAuthenticator {
    client: LoginClient<tonic::transport::Channel>,
}

/// Creates a login client
pub async fn create_grpc_login_client(
    address: http::Uri,
) -> LoginClient<tonic::transport::Channel> {
    let channel = tonic::transport::Channel::builder(address)
        .connect()
        .await
        .expect("failed to connect to gRPC server");
    LoginClient::new(channel)
}

impl GRPCAuthenticator {
    /// Creates a new GRPCAuthenticator wrapping the passed LoginClient
    pub fn new(client: LoginClient<tonic::transport::Channel>) -> Self {
        GRPCAuthenticator { client }
    }

    /// Login synchronously. Let the client handle async code
    fn login(
        &self,
        login: &str,
        password: &str,
    ) -> Result<tonic::Response<LoginResponse>, tonic::Status> {
        let future_login = async move {
            // **Rolls his eyes**
            let mut copied_client = self.client.clone();
            copied_client
                .login(LoginRequest {
                    login: login.to_string(),
                    password: password.to_string(),
                })
                .await
        };

        let runtime = tokio::runtime::Runtime::new().expect("unable to start runtime");

        runtime.block_on(future_login)
    }
}

impl Authenticator for GRPCAuthenticator {
    fn authenticate(&self, username: &str, password: &str) -> Result<bool, AuthenticationError> {
        match self.login(username, password) {
            Ok(response) => Ok(response.get_ref().success),
            Err(e) => Err(AuthenticationError {
                message: e.message().to_string(),
            }),
        }
    }
}
