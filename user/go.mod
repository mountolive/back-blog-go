module github.com/mountolive/back-blog-go/user

go 1.16

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20201119173150-04c754faca46 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/jackc/pgx/v4 v4.9.2
	github.com/joho/godotenv v1.3.0
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc9 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/mountolive/back-blog-go/storefilter v0.1.0 => ../storefilter/pkg/storefilter.go
