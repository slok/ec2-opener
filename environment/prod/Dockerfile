FROM slok/ec2-opener_base
MAINTAINER Xabier Larrakoetxea <slok69@gmail.com>

USER root

RUN apk --update add tar git bash && rm -rf /var/cache/apk/*

RUN go get github.com/Masterminds/glide

RUN mkdir -p /go/src/github.com/slok/ec2-opener/

COPY . /go/src/github.com/slok/ec2-opener/
WORKDIR /go/src/github.com/slok/ec2-opener/

RUN glide install
RUN go build -o /usr/bin/ec2-opener ./main.go

ENTRYPOINT ["/usr/bin/ec2-opener"]
CMD ["-h"]
