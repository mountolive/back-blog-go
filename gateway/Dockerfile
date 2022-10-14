FROM rust:1.64 as builder
WORKDIR /usr/src/gateway
COPY . .
RUN cargo install --path .

FROM debian:buster-slim
RUN apt-get update && apt-get install -y libssl-dev && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/gateway /usr/local/bin/gateway
EXPOSE 8003
CMD ["gateway"]
