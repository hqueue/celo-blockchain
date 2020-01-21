# Docker image to build Celo and other packages

To build celo-blockchain

```
cd <celo-blockchain>
docker run -it --rm -v `pwd`:/go-ethereum hqueue/celo-builder
```

To build other packages such as `geth_exporter`

```
cd <geth_exporter>
docker run -it --rm -v `pwd`:/geth_exporter --entrypoint /bin/sh hqueue/celo-builder

# From inside container..

cd /geth_exporter
go get
go build
```
