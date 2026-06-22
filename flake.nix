{
  description = "rust devshell and package, created by scaffolder";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          name = "rust-devshell";

          packages = with pkgs; [
            cargo
            rustc
            rustfmt
            rust-analyzer
            clippy
            pkg-config
          ];
        };

        packages.hocusfocus = pkgs.rustPlatform.buildRustPackage {
          name = "hocusfocus";
          version = "2.1.1";

          src = ./.;

          cargoLock.lockFile = ./Cargo.lock;
        };

        apps.hocusfocus = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.hocusfocus}/bin/hocusfocus";
        };
      });
}
