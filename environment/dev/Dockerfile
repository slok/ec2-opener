FROM slok/ec2-opener_base
MAINTAINER Xabier Larrakoetxea <slok69@gmail.com>

USER root

RUN go get github.com/derekparker/delve/cmd/dlv
RUN go get github.com/Masterminds/glide
RUN go get github.com/golang/mock/mockgen

RUN mkdir -p /go/src/github.com/slok/ec2-opener/

WORKDIR /go/src/github.com/slok/ec2-opener/

RUN chown -R app:app /go

USER app
