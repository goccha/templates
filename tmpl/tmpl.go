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

func New(template string) (*template.Template, error) {
	body, err := ReadFile(context.Background(), template)
	if err != nil {
		return nil, err
	}
	return NewTextTemplate(template, string(body))
}

func NewHtml(template string) (*html.Template, error) {
	body, err := ReadFile(context.Background(), template)
	if err != nil {
		return nil, err
	}
	return NewHtmlTemplate(template, string(body))
}

func NewTextTemplate(name, value string) (*template.Template, error) {
	tm, err := template.New(name).Funcs(funcMap()).Parse(value)
	if err != nil {
		return nil, err
	}
	return tm, nil
}

func NewHtmlTemplate(name, value string) (*html.Template, error) {
	tm, err := html.New(name).Funcs(funcMap()).Parse(value)
	if err != nil {
		return nil, err
	}
	return tm, nil
}

func Search(ctx context.Context, path, name string) ([]byte, error) {
	return reader.Search(ctx, GetFullPath(path), name)
}
func Read(ctx context.Context, path, name string) ([]byte, error) {
	return reader.Read(ctx, path, name)
}
func ReadFile(ctx context.Context, path string) ([]byte, error) {
	return reader.ReadFile(ctx, path)
}
func GetFullPath(path string) string {
	return reader.GetFullPath(path)
}
