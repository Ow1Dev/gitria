{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    aspire = {
      url = "github:microsoft/aspire";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, aspire, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            aspire.packages.${system}.aspire-cli  # Use from aspire flake, not nixpkgs
            dotnet-sdk_10
            openssl

            go
            gopls          # Go language server (optional but recommended)
            go-tools       # Additional tools like staticcheck

            pnpm
            nodejs_24
            typescript-language-server
          ];
        };
      });
}
