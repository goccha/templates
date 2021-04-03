package tmpl

import (
	"context"
	"crypto/sha256"
	"github.com/goccha/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readFile(_ context.Context, path string) (body []byte, err error) {
	path = GetFullPath(path)
	var f *os.File
	if f, err = os.Open(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	} else {
		if body, err = io.ReadAll(f); err != nil {
			return nil, err
		}
	}
	log.Debug("%s=%x", path, sha256.Sum256(body))
	return
}

type FileTemplateReader struct {
	RootDir string `json:"root_dir"`
}

func (r *FileTemplateReader) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return readFile(ctx, path)
}
func (r *FileTemplateReader) Search(ctx context.Context, template, name string) ([]byte, error) {
	if template == "" {
		return nil, nil
	}
	files, err := os.ReadDir(template)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := filepath.Join(template, file.Name())
		if !file.IsDir() && strings.HasSuffix(filename, name) {
			body, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			log.Debug("%s=%x", filename, sha256.Sum256(body))
			return body, nil
		}
	}
	ctx, cur := GetCursor(ctx, template, r.RootDir, r.separator())
	dir, err := cur.Up()
	if err == RootPath {
		return nil, nil
	}
	return r.Search(ctx, dir, name)
}
func (r *FileTemplateReader) Read(ctx context.Context, path, name string) ([]byte, error) {
	return r.ReadFile(ctx, filepath.Join(path, name))
}
func (r *FileTemplateReader) GetFullPath(path string) string {
	return filepath.Join(r.RootDir, path)
}
func (r *FileTemplateReader) separator() string {
	return string(rune(filepath.Separator))
}
