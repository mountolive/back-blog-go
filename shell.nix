let
  sources = import ./nix/sources.nix { };
  pkgs = import <nixpkgs> { };
in
pkgs.mkShell {
  buildInputs = [
    # system deps
    pkgs.cmake
    pkgs.openssl
    pkgs.pkg-config

    # docker
    pkgs.docker
    pkgs.docker-compose

    # programming languages
    pkgs.go_1_17
    pkgs.rustc
    pkgs.deno

    # misc
    pkgs.cargo
    pkgs.rust-analyzer
    pkgs.protobuf
  ];
}
