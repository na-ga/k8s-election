FROM debian:jessie
MAINTAINER na-ga <katsutoshi.nagaoka@gmail.com>

# so apt-get doesn't complain
ENV DEBIAN_FRONTEND=noninteractive

RUN \
  apt-get update && \
  apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

ADD server server
ENTRYPOINT ["/server"]
