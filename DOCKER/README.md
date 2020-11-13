# Docker

## Supported tags and respective `Dockerfile` links

DockerHub tags for official releases are [here](https://hub.docker.com/r/supermint/supermint/tags/). The "latest" tag will always point to the highest version number.

Official releases can be found [here](https://github.com/vbhp/supermint

The Dockerfile for supermint is not expected to change in the near future. The master file used for all builds can be found [here](https://raw.githubusercontent.com/supermint/supermint/master/DOCKER/Dockerfile).

Respective versioned files can be found <https://raw.githubusercontent.com/supermint/supermint/vX.XX.XX/DOCKER/Dockerfile> (replace the Xs with the version number).

## Quick reference

- **Where to get help:** <https://supermint.com/>
- **Where to file issues:** <https://github.com/vbhp/supermint
- **Supported Docker versions:** [the latest release](https://github.com/moby/moby/releases) (down to 1.6 on a best-effort basis)

## Supermint

Supermintore is Byzantine Fault Tolerant (BFT) middleware that takes a state transition machine, written in any programming language, and securely replicates it on many machines.

For more background, see the [the docs](https://docs.supermint.com/master/introduction/#quick-start).

To get started developing applications, see the [application developers guide](https://docs.supermint.com/master/introduction/quick-start.html).

## How to use this image

### Start one instance of the Supermintore with the `kvstore` app

A quick example of a built-in app and Supermintore in one container.

```sh
docker run -it --rm -v "/tmp:/supermint" supermint/supermint init
docker run -it --rm -v "/tmp:/supermint" supermint/supermint node --proxy_app=kvstore
```

## Local cluster

To run a 4-node network, see the `Makefile` in the root of [the repo](https://github.com/vbhp/supermintun:

```sh
make build-linux
make build-docker-localnode
make localnet-start
```

Note that this will build and use a different image than the ones provided here.

## License

- Supermint license is [Apache 2.0](https://github.com/vbhp/supermint

## Contributing

Contributions are most welcome! See the [contributing file](https://github.com/vbhp/supermint) for more information.
