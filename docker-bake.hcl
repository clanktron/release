variable "cache_settings" {
  default = [
    {
      type = "registry",
      ref = "${image}-cache:cache"
    }
  ]
}

variable "registry" { default = "ghcr.io" }
variable "org" { default = "clanktron" }
variable "repo" { default = "git-release" }
variable "tag" { default = "dev" }
variable "image" { default = "${registry}/${org}/${repo}"}
variable "ref" { default = "${image}:${tag}"}

target "default" {
  inherits = ["cache"]
  output = [
    "type=image,name=${ref},mode=max"
  ]
}

target "default-push" {
  inherits = ["default"]
  output = [
    "type=registry,name=${ref}"
  ]
}

target "cache" {
  cache-from = cache_settings
  cache-to = cache_settings
}

target "bin" {
  target = "bin"
  output = [
    "type=local,dest=./dist,platform-split=false",
  ]
}

group "all" {
  targets = ["containers", "linux", "darwin", "windows"]
}

group "release-binaries" {
  targets = ["linux", "darwin", "windows"]
}

target "containers" {
  inherits = ["cache"]
  target = "production"
  platforms = [
    "linux/arm64",
    "linux/amd64",
    "linux/riscv64"
  ]
  output = [
    "type=registry,name=${ref}",
  ]
}

target "linux" {
  inherits = ["cache","bin"]
  matrix = {
    arch = ["arm64", "amd64", "arm", "riscv64"]
  }
  name = "linux-${arch}"
  args = {
    BINOS = "linux",
    BINARCH = "${arch}"
  }
}

target "darwin" {
  inherits = ["cache","bin"]
  matrix = {
    arch = ["arm64", "amd64"]
  }
  name = "darwin-${arch}"
  args = {
    BINOS = "darwin",
    BINARCH = "${arch}"
  }
}

target "windows" {
  inherits = ["cache","bin"]
  matrix = {
    arch = ["arm64", "amd64"]
  }
  name = "windows-${arch}"
  args = {
    BINOS = "windows",
    BINARCH = "${arch}"
  }
}
