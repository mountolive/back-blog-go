let
  pkgs = import <nixpkgs> { };
in
pkgs.mkShell {
  buildInputs = [
    # system deps
    pkgs.cmake
    pkgs.openssl

    # docker
    pkgs.docker
    pkgs.docker-compose

    # programming languages
    pkgs.go_1_17
    pkgs.rustc

    # misc
    pkgs.cargo
    pkgs.rust-analyzer
    pkgs.protobuf
  ];

  shellHook = ''
    echo "Starting infra circus"
    make local-infra
    echo "Sourcing local env"
    source .env.local
  '';

  
}
