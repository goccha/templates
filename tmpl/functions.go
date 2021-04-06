package tmpl

import (
	"bytes"
	"context"
	"github.com/goccha/log"
	"html/template"
	"os"
	"time"
)

var _build func() map[string]interface{}

func embedded() map[string]interface{} {
	return map[string]interface{}{
		"timeFormat": timeFormat,
		"import":     importTemplate,
		"importHtml": importHTML,
		"env":        getEnv,
		"date":       date,
		"now":        now,
	}
}

type Format string

func importTemplate(filePath string, variables interface{}) (string, error) {
	if data, err := reader.ReadFile(context.Background(), filePath); err != nil {
		log.Error("%+v", err)
		return "", err
	} else if data == nil {
		return "", nil
	} else {
		buf := new(bytes.Buffer)
		if tm, err := data.Text(); err != nil {
			return "", err
		} else if err = tm.Execute(buf, variables); err != nil {
			return "", err
		}
		return string(buf.Bytes()), nil
	}
}

func importHTML(filePath string, variables interface{}) (template.HTML, error) {
	if data, err := reader.ReadFile(context.Background(), filePath); err != nil {
		log.Error("%+v", err)
		return "", err
	} else if data == nil {
		return "", nil
	} else {
		buf := new(bytes.Buffer)
		if tm, err := data.Html(); err != nil {
			return "", err
		} else if err = tm.Execute(buf, variables); err != nil {
			return "", err
		}
		return template.HTML(buf.Bytes()), nil
	}
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

func getEnv(name string, defaultValue ...string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	if defaultValue != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func date() string {
	return time.Now().Format("2006-01-02")
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
