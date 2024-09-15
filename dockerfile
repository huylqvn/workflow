FROM golang:1.22.1-alpine AS builder
ENV GO111MODULE=on
WORKDIR /go/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server .

FROM alpine:latest
RUN apk update && \
    apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
RUN echo "Asia/Ho_Chi_Minh" >  /etc/timezone
COPY --from=builder /go/src/app/cmd/server /server
COPY --from=builder /go/src/app/.env /go/src/app/.env

WORKDIR /go/src/app
EXPOSE 8088
CMD ["/server"]
