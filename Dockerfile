FROM golang:1.6.2

WORKDIR /go/src

RUN go get github.com/tools/godep \
&& go get github.com/fsouza/go-dockerclient \
&& apt-get update \
&& apt-get install -y apt-transport-https ca-certificates vim \
&& curl -O https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz \
&& tar -zxvf docker-latest.tgz \
&& mv docker/docker /usr/bin/docker \
&& chmod +x /usr/bin/docker \
&& rm -rf docker

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GO15VENDOREXPERIMENT 1