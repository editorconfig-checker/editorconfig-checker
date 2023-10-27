{ stdenv, buildGoPackage, fetchgit, pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  # load version from file
  versionPath = toString ./VERSION;
  versionData = builtins.readFile versionPath;
  versionLen = lib.stringLength versionData;
  # trim trailing newline
  version = lib.substring 0 (versionLen - 1) versionData;

  # set the version dynamically at build time
  buildFlagsArray = ''
    -ldflags=-X main.version=${version}
  '';

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
