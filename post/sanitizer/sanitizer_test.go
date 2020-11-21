package sanitizer

import (
	"testing"

	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/stretchr/testify/require"
)

type sanitizeContentCase struct {
	Description,
	Content,
	Sanitized string
}

func TestSanitizer(t *testing.T) {
	sanitizer := NewSanitizer()
	genericErr := "Got: %s Expected: %s"

	t.Run("Canary", func(t *testing.T) {
		var _ usecase.ContentSanitizer = NewSanitizer()
	})

	t.Run("SanitizeContent", func(t *testing.T) {
		testCases := []sanitizeContentCase{
			{
				Description: "Removes xss attempt html",
				Content:     `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
				Sanitized:   "<p><a href=\"http://www.google.com\" rel=\"nofollow\">Google</a></p>\n",
			},
			{
				Description: "Removes xss attempt markdown 1",
				Content:     "[Click Me](javascript:alert('jaquer'))",
				Sanitized:   "<p><a title=\"jaquer\">Click Me</a>)</p>\n",
			},
			{
				Description: "Removes xss attempt markdown 2",
				Content:     `![hi](https://www.hello.com/image.png"onload="alert('XSS'))`,
				Sanitized:   "<p><img src=\"https://www.hello.com/image.png\" alt=\"hi\"/>)</p>\n",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Description, func(t *testing.T) {
				t.Log(tc.Description)
				got := sanitizer.SanitizeContent(tc.Content)
				require.True(t, got == tc.Sanitized, genericErr, got, tc.Sanitized)
			})
		}
	})
}
