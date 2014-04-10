FROM       ubuntu:13.10
MAINTAINER Johannes 'fish' Ziemke <fish@docker.com> @discordianfish

RUN        apt-get update && apt-get install -yq curl git mercurial make
RUN        curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV        PATH    /usr/local/go/bin:$PATH
ENV        GOPATH  /go

ADD        . /usr/src/docker-exporter
RUN        cd /usr/src/docker-exporter && \
           go get -d && go build && cp docker-exporter /

ENTRYPOINT [ "/docker-exporter" ]
EXPOSE     8080
