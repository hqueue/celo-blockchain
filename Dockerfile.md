# Dockerfile.hqueue

This image starts to run `make geth` from `/go-ethereum` directory which is mounted volume of `celo-blockchain`.

For example,

```
$ git clone git@github.com:celo-org/celo-blockchain.git
$ cd celo-blockchain
$ docker pull hqueue/celo-builder
$ docker run -v `pwd`:/go-ethereum hqueue/celo-builder
```

If you want to run into the container,

```
$ docker run -it --rm --entrypoint /bin/sh hqueue/celo-builder
```
