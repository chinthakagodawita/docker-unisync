# docker-unisync
[![Build Status](https://travis-ci.org/chinthakagodawita/docker-unisync.svg?branch=master)](https://travis-ci.org/chinthakagodawita/docker-unisync)

docker-unisync mounts a specified directory into your Docker Machine and recursively watches that directory for changes. Changed files are copied across to your Docker Machine using [Unison](https://www.cis.upenn.edu/~bcpierce/unison/).

Files changed on your Docker Machine are copied back to the source folder as well.

At the moment, only Mac OS X is supported due to reliance on the FSEvents API. Once the Go fsnotify library [supports recursive directory watching](https://github.com/go-fsnotify/fsnotify/issues/18), this may be revisited.

## Installation

Binaries will be available very soon. For the meantime, install Go 1.5 and install via `go get`:

```bash
brew install go
go get chinthakagodawita/docker-unisync
```

## Usage
First make sure your Docker Machine is running:

```bash
docker-machine up mymachine
```

Then `cd` to the directory you want to sync and run `docker-unisync`:

```bash
cd /my/sync/dir
docker-unisync mymachine
```

Alternatively, provide source and destination parameters:

```bash
docker-unisync --source=/my/sync/dir --destination=/dir/on/machine mymachine
```

See `docker-unisync --help` for a full list of options.

## Building

All dependencies are managed via [godep](https://github.com/tools/godep) and the [Go 1.5 Vendor Experiment](https://golang.org/s/go15vendor).

To build, install `godep`, `cd` to this directory and run:

```bash
export GO15VENDOREXPERIMENT=1
godep restore
make
```

## Releasing

1. Increment the version number in the `Makefile`
2. Commit and tag the commit with the same version number with a 'v' prefix (e.g. `v0.0.1`)
3. Push to the `master` branch, Travis will take care of building a releasing a binary.
