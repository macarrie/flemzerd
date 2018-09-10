FROM golang:1.11-alpine as builder

COPY . /src
WORKDIR /src

RUN apk update && apk add make clang git gcc g++ ca-certificates rust cargo && rm -rf /var/cache/apk/*
RUN make bin

# Build vidocq latest tagged release
RUN git clone https://github.com/macarrie/vidocq --branch v0.1 dep_vidocq && cd dep_vidocq && cargo build --release && cp target/release/vidocq ../bin/



FROM alpine:latest

RUN apk update && apk add gcc

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN mkdir -p /var/lib/flemzerd/server/ui
RUN mkdir -p /var/lib/flemzerd/db
RUN mkdir -p /downloads
RUN mkdir -p /shows
RUN mkdir -p /movies

# Copy application files
WORKDIR /app
COPY --from=builder /src/bin/flemzerd /app
COPY --from=builder /src/package/flemzerd_*/ui/* /var/lib/flemzerd/server/ui/
COPY --from=builder /src/bin/vidocq /usr/bin/vidocq

# Run
EXPOSE 8400
CMD ["./flemzerd", "-d"]
