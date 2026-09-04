{ stdenv, buildGoPackage, fetchgit, pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  version = "v4.0.1" # x-release-please-version

  name = "editorconfig-checker-${version}";

  goPackagePath = "github.com/editorconfig-checker/editorconfig-checker/v2";

  src = lib.cleanSourceWith {
    filter = name: type: builtins.match ".*tests.*" name == null;
    src = (lib.cleanSource ./.);
  };

  goDeps = ./deps.nix;
}
