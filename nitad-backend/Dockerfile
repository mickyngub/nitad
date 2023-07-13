FROM golang:1.17

WORKDIR /opt/code/

ADD ./ /opt/code/

RUN go build -o nitad main.go

ENTRYPOINT ["./nitad"]

