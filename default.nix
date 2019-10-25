{ stdenv, buildGoPackage, fetchgit, pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  name = "editorconfig-checker-${version}";
  version = "2.0.3";

  goPackagePath = "github.com/editorconfig-checker/editorconfig-checker";

  src = lib.cleanSourceWith {
    filter = name: type: builtins.match ".*tests.*" name == null;
    src = (lib.cleanSource ./.);
  };

  goDeps = ./deps.nix;
}
