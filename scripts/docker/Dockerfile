FROM golang:1.19-alpine

USER root
ENV TZ=Asia/Tokyo \
  LANG=ja_JP.UTF-8

ARG uid
ARG gid
ARG user

RUN adduser -u $uid -g $gid -D $user \
  && chown -R $user:$user /var \
  && chmod -R 0755 /var \
  && install -m 0755 -o $user -g $user -d /var/www \
  && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# Goのパスを通す
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin

RUN apk update \
  && apk add git \
  && go install golang.org/x/tools/gopls@latest
USER $user

WORKDIR /var/batch
