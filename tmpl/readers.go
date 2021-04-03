package tmpl

import (
	"context"
)

type TemplateReader interface {
	ReadFile(ctx context.Context, filePath string) ([]byte, error)
	Search(ctx context.Context, path, name string) ([]byte, error)
	Read(ctx context.Context, path, name string) ([]byte, error)
	GetFullPath(path string) string
}

var reader TemplateReader
