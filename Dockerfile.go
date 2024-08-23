# FROM golang:1.20.14-bookworm as gobuilder
# FROM golang:1.20.14-bullseye as gobuilder
FROM ubuntu:16.04 as gobuilder
RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update --quiet --quiet && \
    apt-get install --quiet --quiet --yes \
    --no-install-recommends --no-install-suggests \
    openssl osslsigncode make zip mingw-w64 wget tar ca-certificates \
    && apt-get --quiet --quiet clean \
    && rm --recursive --force /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN wget https://go.dev/dl/go1.20.14.linux-amd64.tar.gz && \
    tar -xf go1.20.14.linux-amd64.tar.gz
ENV PATH="${PATH}:/go/bin:/root/go/bin"
WORKDIR /build
COPY . .
CMD ["make","build"]

