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
