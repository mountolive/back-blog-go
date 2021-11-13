# My personal blog's backend

It's intended to have a "microservices" layout

_Overkill_ is my business, my business is good... But loosely following a hexagonal architecture

Not intended to be a simple and elegant solution... I'm just using this as an excuse for testing out different transports and get in touch with cumbersome ops

## TODO

-  A huge one: explain the architecture...
-  Another huge one: **code quality is frankly crappy at some layers** -> I need
    to rework some stuff a bit.

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

### Alternative (nix <3)!

- Install [`nix`](https://nixos.org/guides/install-nix.html)
- Install [`direnv`](https://direnv.net/docs/hook.html)
- Run (if using bash, follow the linked source if not):
```bash
eval "$(direnv hook bash)"
```
(add this to your `.bashrc`)
- Run:
```bash
echo "use nix" > .envrc && direnv allow
```

Off to go

### Set the .env variables for local development

If using the previous approach, these variable will be loaded on every time you
enter the project's directory

```
NATS_PORT
USERS_DB_USER
USERS_DB_PASS
USERS_DB_NAME
USERS_DB_HOST
USERS_DB_PORT
USERS_PORT
USERS_ADMIN_EMAIL
USERS_ADMIN_USERNAME
USERS_ADMIN_PASSWORD
USERS_ADMIN_FIRST_NAME
USERS_ADMIN_LAST_NAME
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
USER_SERVICE_ADDRESS
TOKEN_TTL
TOKEN_SALT
POST_REST_API_URL
GATEWAY_PORT
POST_SERVER_USER
POST_SERVER_PASSWORD
POST_SERVER_SUBJECT
POST_SERVER_HOST
POST_SERVER_PORT
API_URL
```

It's best to set them in a local `.env.local`, prefixed with `export`.

Afterwards you can run:

```
source .env.local
```
