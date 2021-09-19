# My personal blog's backend

It's intended to have a "microservices" layout

_Overkill_ is my business, my business is good... But loosely following a hexagonal architecture

Not intended to be a simple and elegant solution... I'm just using this as an excuse for testing out different transports and get in touch with cumbersome ops

## TODO

-  A huge one: explain the architecture...

## CARGO HEADS-UP

If using ssh `insteadOf` for global `.gitconfig`, make sure to:

```bash
eval `ssh-agent -s`
ssh-add
```

### Requirements

- Go (>1.15)
- Rust (>1.54)
- Protoc (`sudo apt install -y protobuf-compiler`, if using `apt`)

### Set the .env variables for local development

```
NATS_PORT
USERS_DB_USER
USERS_DB_PASS
USERS_DB_NAME
USERS_DB_HOST
USERS_DB_PORT
USERS_PORT
POSTS_DB_USER
POSTS_DB_PASS
POSTS_DB_NAME
POSTS_DB_HOST
POSTS_DB_PORT
POSTS_NATS_HOST
POSTS_NATS_PORT
POSTS_USERS_GRPC_HOST
POSTS_USERS_GRPC_PORT
POSTS_NATS_SUBSCRIPTION_NAME
POSTS_NATS_DEADLETTER_NAME
POSTS_HTTP_PORT
```
