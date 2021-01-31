module github.com/mountolive/back-blog-go/post

go 1.14

replace github.com/mountolive/back-blog-go/storefilter v0.1.0 => ../storefilter/pkg/storefilter.go

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20201202124332-91328d7c60e7 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v4 v4.9.2
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
)
