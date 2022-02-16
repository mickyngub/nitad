FROM golang:1.17

WORKDIR /opt/code/

ADD ./ /opt/code/

RUN go build -o nitad server.go

ENTRYPOINT ["./nitad"]