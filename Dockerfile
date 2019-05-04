FROM golang:1.11-alpine as builder

COPY . /src
WORKDIR /src

RUN apk update && apk add make clang git gcc g++ ca-certificates && rm -rf /var/cache/apk/*
RUN make bin




FROM rust:1.28.0-stretch as rust-builder

WORKDIR /src

RUN rustup target add x86_64-unknown-linux-musl
RUN git clone https://github.com/macarrie/vidocq --branch v0.1.1 dep_vidocq && cd dep_vidocq && cargo build --release --target x86_64-unknown-linux-musl && cp target/x86_64-unknown-linux-musl/release/vidocq ../





FROM alpine:latest

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
COPY --from=rust-builder /src/vidocq /usr/bin/vidocq

# Run
EXPOSE 8400
CMD ["./flemzerd", "-d"]
