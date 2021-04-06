package tmpl

import (
	"context"
	html "html/template"
	"io"
	"text/template"
)

type TemplateReader interface {
	ReadFile(ctx context.Context, filePath ...string) (*TemplateData, error)
	Search(ctx context.Context, path, name string) ([]byte, error)
	Read(ctx context.Context, path, name string, nested ...string) (*TemplateData, error)
	GetFullPath(path string) string
}

var reader TemplateReader

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

type TemplateData struct {
	Files []File
}

func (d *TemplateData) Text() (tm *template.Template, err error) {
	for i, v := range d.Files {
		if i == 0 {
			if tm, err = template.New(v.Name).Funcs(funcMap()).Parse(v.Body); err != nil {
				return nil, err
			}
		} else {
			if _, err = tm.New(v.Name).Parse(v.Body); err != nil {
				return nil, err
			}
		}
	}
	return
}
func (d *TemplateData) Html() (tm *html.Template, err error) {
	for i, v := range d.Files {
		if i == 0 {
			if tm, err = html.New(v.Name).Funcs(funcMap()).Parse(v.Body); err != nil {
				return nil, err
			}
		} else {
			if _, err = tm.New(v.Name).Parse(v.Body); err != nil {
				return nil, err
			}
		}
	}
	return
}

type File struct {
	Name string
	Body string
}
