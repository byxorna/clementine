FROM google/golang
MAINTAINER Gabe Conradi <gummybearx@gmail.com>

ENV WORKDIR $GOPATH/src/github.ewr01.tumblr.net/gabe/clementine
RUN mkdir -p $WORKDIR
COPY . $WORKDIR
WORKDIR $WORKDIR
RUN go get github.com/tools/godep
RUN godep restore
RUN godep go build
EXPOSE 8080
ENTRYPOINT ["./clementine","-conf=example/config.json"]
