package tmpl

import (
	"context"
	"github.com/goccha/envar"
	"github.com/goccha/errors"
	"github.com/goccha/log"
	html "html/template"
	"text/template"
)

var RootPath = errors.New("root dir")
var TemplatePath = "GOCCHA_TEMPLATE_PATH"

const (
	DefaultDir = "templates"
)

func init() {
	log.Debug("emails.init()")
	Setup(&FileTemplateReader{RootDir: envar.Get(TemplatePath).String(DefaultDir)})
}

func Setup(r TemplateReader, f ...func() map[string]interface{}) {
	if r != nil {
		log.Debug("setup %v", r)
		reader = r
	}
	if f != nil && len(f) > 0 {
		_build = f[0]
	}
}

var functions map[string]interface{}

func funcMap() map[string]interface{} {
	if functions == nil {
		functions = NewFuncMap()
	}
	return functions
}

func New(templates ...string) (tm *template.Template, err error) {
	var data *TemplateData
	if data, err = ReadFile(context.Background(), templates...); err != nil {
		return nil, err
	} else {
		tm, err = data.Text()
	}
	return
}

func NewHtml(template ...string) (tm *html.Template, err error) {
	var data *TemplateData
	if data, err = ReadFile(context.Background(), template...); err != nil {
		return nil, err
	} else {
		tm, err = data.Html()
	}
	return
}

func Search(ctx context.Context, path, name string) ([]byte, error) {
	return reader.Search(ctx, GetFullPath(path), name)
}
func Read(ctx context.Context, path, name string) (*TemplateData, error) {
	return reader.Read(ctx, path, name)
}
func ReadFile(ctx context.Context, filePath ...string) (*TemplateData, error) {
	return reader.ReadFile(ctx, filePath...)
}
func GetFullPath(path string) string {
	return reader.GetFullPath(path)
}
