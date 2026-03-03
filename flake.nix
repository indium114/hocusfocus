{
  description = "Hocusfocus devshell and package";

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
          name = "hocusfocus-devshell";

          packages = with pkgs; [
            go
            gopls
            gotools
            delve
          ];
        };

        packages.hocusfocus = pkgs.buildGoModule {
          pname = "hocusfocus";
          version = "2026.02.06-a";

          src = self;

          vendorHash = "sha256-SMhllO87YlmySHroKfPq1pHb67CwHaZ3XMp3t983etc=";

          subPackages = [ "." ];
          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "A simple TUI productivity tool";
            license = licenses.mit;
            platforms = platforms.all;
          };
        };

        apps.hocusfocus = {
          type = "app";
          program = "${self.packages.${system}.hocusfocus}/bin/hocusfocus";
        };
      });
}
