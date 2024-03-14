{
  nixConfig.bash-prompt = "[nix(cedille-hackathon)] ";
  inputs = { 
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
     
  };


  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs {
          config = { allowUnfree = true; };
        };
    in {
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = [
          pkgs.go_1_22
          pkgs.air
          pkgs.delve
          pkgs.vscode
        ];
        hardeningDisable = [ "fortify" ];
      };
    };
}

