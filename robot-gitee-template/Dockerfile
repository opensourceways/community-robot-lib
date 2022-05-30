FROM golang:latest as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/{ROBOT_NAME}
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o {ROBOT_NAME} .

# copy binary config and utils
FROM alpine:3.14
COPY  --from=BUILDER /go/src/github.com/opensourceways/{ROBOT_NAME}/{ROBOT_NAME} /opt/app/{ROBOT_NAME}

ENTRYPOINT ["/opt/app/{ROBOT_NAME}"]