module github.com/mountolive/back-blog-go/post

go 1.14

replace github.com/mountolive/back-blog-go/storefilter v0.1.0 => ../storefilter/pkg/storefilter.go

require (
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
)
