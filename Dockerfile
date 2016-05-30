FROM golang:1.6-alpine
MAINTAINER Xabier Larrakoetxea <slok69@gmail.com>

# Create an user with the same uid/gid as the user running docker-compose in development to avid permissions conflicts
ARG uid=1000
ARG gid=1000

RUN addgroup -g $gid app
RUN adduser -D -u $uid -G app app

USER app

RUN echo 'PATH=$PATH:/code/bin' > ~/.bashrc

WORKDIR /code
