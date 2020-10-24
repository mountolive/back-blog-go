module github.com/mountolive/back-blog-go/user

go 1.14

require (
	github.com/mountolive/back-blog-go/storehelper v0.0.0-20201024170606-9df49c46ba25
	github.com/stretchr/testify v1.6.1
)

replace github.com/mountolive/back-blog-go/storehelper v0.1.0 => ../storehelper/pkg/storehelper.go
