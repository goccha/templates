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

func readFile(_ context.Context, filePath ...string) (data *TemplateData, err error) {
	data = &TemplateData{
		Files: make([]File, 0, len(filePath)),
	}
	for _, path := range filePath {
		fullpath := GetFullPath(path)
		var f *os.File
		if f, err = os.Open(fullpath); err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			}
			return nil, err
		} else {
			if body, err := io.ReadAll(f); err != nil {
				return nil, err
			} else {
				data.Files = append(data.Files, File{
					Name: path,
					Body: string(body),
				})
				log.Debug("%s=%x", path, sha256.Sum256(body))
			}
		}
	}
	return
}

type FileTemplateReader struct {
	RootDir string `json:"root_dir"`
}

func (r *FileTemplateReader) ReadFile(ctx context.Context, filePath ...string) (*TemplateData, error) {
	return readFile(ctx, filePath...)
}
func (r *FileTemplateReader) Search(ctx context.Context, dirPath, name string) ([]byte, error) {
	if dirPath == "" {
		return nil, nil
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := filepath.Join(dirPath, file.Name())
		if !file.IsDir() && strings.HasSuffix(filename, name) {
			body, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			log.Debug("%s=%x", filename, sha256.Sum256(body))
			return body, nil
		}
	}
	ctx, cur := GetCursor(ctx, dirPath, r.RootDir, r.separator())
	dir, err := cur.Up()
	if err == RootPath {
		return nil, nil
	}
	return r.Search(ctx, dir, name)
}
func (r *FileTemplateReader) Read(ctx context.Context, dirPath, name string, nested ...string) (*TemplateData, error) {
	if nested != nil && len(nested) > 0 {
		filePath := make([]string, 0, len(nested)+1)
		filePath = append(filePath, filepath.Join(dirPath, name))
		filePath = append(filePath, nested...)
		return r.ReadFile(ctx, filePath...)
	} else {
		return r.ReadFile(ctx, filepath.Join(dirPath, name))
	}
}
func (r *FileTemplateReader) GetFullPath(path string) string {
	return filepath.Join(r.RootDir, path)
}
func (r *FileTemplateReader) separator() string {
	return string(rune(filepath.Separator))
}
