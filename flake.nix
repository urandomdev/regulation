{
  description = "Development environment for the Regulation monorepo";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        formatter = pkgs.nixfmt-rfc-style;

        devShells = {
          default = pkgs.mkShell {
            name = "regulation-dev";
            packages = with pkgs; [
              go
              nodejs
              git
              gnumake
              pkg-config
              openssl
              sqlite
              docker-client
              docker-compose
            ];
            env = {
              GO111MODULE = "on";
              CONFIG_DIR = builtins.toString ./config;
            };
          };

          api = pkgs.mkShell {
            name = "regulation-api";
            packages = with pkgs; [
              go
              pkg-config
              openssl
              sqlite
            ];
            env = {
              GO111MODULE = "on";
              GOPATH = "$HOME/go";
              CONFIG_DIR = builtins.toString ./config;
            };
          };

          frontend = pkgs.mkShell {
            name = "regulation-frontend";
            packages = with pkgs; [
              nodejs
            ];
          };
        };
      });
}
