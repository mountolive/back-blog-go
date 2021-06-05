module github.com/mountolive/back-blog-go/post

go 1.16

replace github.com/mountolive/back-blog-go/storefilter v0.1.0 => ../storefilter/pkg/storefilter.go

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/containerd/continuity v0.0.0-20201202124332-91328d7c60e7 // indirect
	github.com/golang/protobuf v1.5.0
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/jackc/pgx/v4 v4.9.2
	github.com/joho/godotenv v1.3.0
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/nats-io/nats-server/v2 v2.2.6 // indirect
	github.com/nats-io/nats.go v1.11.0
	github.com/ory/dockertest/v3 v3.6.3
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.26.0 // indirect
)
