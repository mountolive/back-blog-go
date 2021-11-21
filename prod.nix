let
  sources = import ./nix/sources.nix { };
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
    pkgs.deno
  ];

  shellHook = ''
    echo "---------------------"
    echo "Starting infra circus"
    echo "---------------------"
    docker-compose --file docker-compose-infra.yml --env-file .env up -d
    echo "---------------------"
    echo "Sourcing env"
    echo "---------------------"
    source .env
    rm .env
  '';
}
