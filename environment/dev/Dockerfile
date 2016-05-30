FROM slok/ec2-open-port_base
MAINTAINER Xabier Larrakoetxea <slok69@gmail.com>

USER root

RUN apk --update add tar git bash && rm -rf /var/cache/apk/*

RUN go get github.com/derekparker/delve/cmd/dlv
RUN go get github.com/Masterminds/glide

RUN mkdir -p /go/src/github.com/slok/ec2-open-port/

COPY . /go/src/github.com/slok/ec2-open-port/
WORKDIR /go/src/github.com/slok/ec2-open-port/

RUN glide install

RUN chown -R app:app /go

USER app
