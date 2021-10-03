// ! gRPC authenticator's implementation

use crate::auth::AuthenticationError;
use crate::user::login_client::LoginClient;
use crate::user::LoginRequest;

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

/// Logins against authentication server
pub async fn authenticate(
    client: LoginClient<tonic::transport::Channel>,
    login: &str,
    password: &str,
) -> Result<bool, AuthenticationError> {
    let mut copied_client = client.clone();
    match copied_client
        .login(LoginRequest {
            login: login.to_string(),
            password: password.to_string(),
        })
        .await
    {
        Ok(response) => Ok(response.get_ref().success),
        Err(e) => Err(AuthenticationError {
            message: e.message().to_string(),
        }),
    }
}
