{
  inputs = {
    # pull down the utils so i dont have to define every sys-arch triple
    utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      # the actual packages
      nixpkgs,
      # use the utils from inputs
      utils,
      ...
    }:
    # runs the given function on every system triple
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        # the dev shell
        devShell = pkgs.mkShell {
          # tools needed at build time
          nativeBuildInputs = with pkgs; [
            pkg-config
            go
            gopls
          ];
          # things needed by the package
          buildInputs = with pkgs; [

          ];
        };

        # Make a package using someoness handy buildGoModule
        packages.default = pkgs.buildGoModule {
          pname = "boogaloo";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-kXQ5b7pd+dRYItZmVFZAFWmzj0WVsvWAzOa6CKKaaLc=";
        };
      }
    );
}
