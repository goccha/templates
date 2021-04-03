package tmpl

import (
	"context"
	"strings"
)

const dirCur = "github.com/goccha/templates/cursor"

func GetCursor(ctx context.Context, path, rootDir, sep string) (context.Context, *cursor) {
	var cur *cursor
	if v := ctx.Value(dirCur); v == nil {
		cur = newCursor(path, rootDir, sep)
		ctx = context.WithValue(ctx, dirCur, cur)
	} else {
		cur = v.(*cursor)
	}
	return ctx, cur
}

type cursor struct {
	root  string
	dirs  []string
	index int
	sep   string
}

func newCursor(path, rootDir, sep string) *cursor {
	path = path[len(rootDir)+1:]
	dirs := strings.Split(path, sep)
	fp := &cursor{
		root: rootDir,
		dirs: make([]string, 0, len(dirs)),
		sep:  sep,
	}
	fp.dirs = append(fp.dirs, rootDir)
	for _, d := range dirs {
		fp.dirs = append(fp.dirs, d)
	}
	fp.index = len(fp.dirs) - 1
	return fp
}
func (c *cursor) Up() (string, error) {
	i := c.index
	if i > 0 {
		c.index--
		return c.String(), nil
	}
	return "", RootPath
}
func (c *cursor) String() string {
	builder := strings.Builder{}
	for i, d := range c.dirs {
		if i > c.index {
			break
		}
		if i > 0 {
			builder.WriteString(c.sep)
		}
		builder.WriteString(d)
	}
	return builder.String()
}
