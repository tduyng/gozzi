package builder

import (
	"path"

	"github.com/tduyng/gozzi/app/template/funcs"
)

func (b *Builder) buildTagPermalink(tag string) string {
	return path.Join("/tags", funcs.Urlize(tag)) + "/"
}

func (b *Builder) buildTagURL(tagLink string) string {
	return b.site.BaseURL + tagLink
}
