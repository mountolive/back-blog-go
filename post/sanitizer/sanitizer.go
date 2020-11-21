package sanitizer

import (
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	policy *bluemonday.Policy
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{bluemonday.UGCPolicy()}
}

func (s *Sanitizer) SanitizeContent(content string) string {
	html := markdown.ToHTML([]byte(content), nil, nil)
	return string(s.policy.SanitizeBytes(html))
}
