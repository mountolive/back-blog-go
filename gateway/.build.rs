// Rename this file to build.rs to regenerate proto's code
fn main() {
    tonic_build::configure()
        .build_server(false)
        .out_dir("src")
        .compile(&["../proto/user/login.proto"], &["../proto/user"])
        .unwrap();
}
