FROM golang:alpine as builder
ADD . /src
RUN apk update && apk add make clang git
RUN go get -u golang.org/x/vgo
RUN cd /src && make bin

FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/bin/flemzerd /app
EXPOSE 8400
CMD ["./flemzerd", "-d"]
