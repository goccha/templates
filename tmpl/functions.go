package tmpl

import (
	"bytes"
	"context"
	"github.com/goccha/log"
	html "html/template"
	"os"
	"text/template"
	"time"
)

var _build func() map[string]interface{}

func embedded() map[string]interface{} {
	return map[string]interface{}{
		"timeFormat": timeFormat,
		"import":     importTemplate,
		"env":        getEnv,
		"date":       date,
		"now":        now,
	}
}

type Format string

const (
	HTML Format = "html"
	TEXT Format = "text"
)

func importTemplate(filePath string, variables interface{}, format ...Format) (string, error) {
	body, err := reader.ReadFile(context.Background(), filePath)
	if err != nil {
		log.Error("%+v", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if format != nil && len(format) > 0 && format[0] == HTML {
		if tm, err := html.New(filePath).Funcs(NewFuncMap()).Parse(string(body)); err != nil {
			return "", err
		} else if err = tm.Execute(buf, variables); err != nil {
			return "", err
		}
	} else {
		if tm, err := template.New(filePath).Funcs(NewFuncMap()).Parse(string(body)); err != nil {
			return "", err
		} else if err = tm.Execute(buf, variables); err != nil {
			return "", err
		}
	}
	return string(buf.Bytes()), nil
}

func NewFuncMap() map[string]interface{} {
	e := embedded()
	if _build != nil {
		m := _build()
		for k, v := range m {
			e[k] = v
		}
	}
	return e
}

func timeFormat(layout string, value interface{}) (string, error) {
	var t time.Time
	var err error
	switch value.(type) {
	case string:
		t, err = time.Parse(time.RFC3339Nano, value.(string))
		if err != nil {
			return "", err
		}
	case int64:
		t = time.Unix(value.(int64), 0)
	case time.Time:
		t = value.(time.Time)
	}
	return t.Format(layout), nil
}

func getEnv(names ...string) string {
	for _, key := range names {
		v, ok := os.LookupEnv(key)
		if ok {
			return v
		}
	}
	return ""
}

func date() string {
	return time.Now().Format("2006-01-02")
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
