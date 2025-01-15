{ stdenv, buildGoPackage, fetchgit, pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  version = "v3.1.2" # x-release-please-version

  # create link so the tool can also be executed as `ec`
  postInstall = ''
    ln -s $bin/bin/editorconfig-checker $bin/bin/ec
  '';

  name = "editorconfig-checker-${version}";

  goPackagePath = "github.com/editorconfig-checker/editorconfig-checker/v2";

  src = lib.cleanSourceWith {
    filter = name: type: builtins.match ".*tests.*" name == null;
    src = (lib.cleanSource ./.);
  };

  goDeps = ./deps.nix;
}
