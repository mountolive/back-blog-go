fn main() {
    tonic_build::configure()
        .build_server(false)
        .out_dir("src")
        .compile(&["../proto/user/login.proto"], &["../proto/user"])
        .unwrap();
}
