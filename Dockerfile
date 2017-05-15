FROM golang:1.8.1-alpine
ADD . /snitch
WORKDIR /snitch
ENV GOPATH /snitch
RUN go build
ENTRYPOINT ["/bin/sh", "-c", "/snitch/snitch ${*}", "--"]
