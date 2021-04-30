package tmpl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestReadFile(t *testing.T) {
	Setup(&FileTemplateReader{RootDir: "../.test/templates"}, func() map[string]interface{} {
		return map[string]interface{}{
			"test": func() string {
				return "%テスト%"
			},
		}
	})
	files := []string{"index.html", "header.html", "about.html"}
	if data, err := ReadFile(context.TODO(), files...); err != nil {
		t.Error(err)
		return
	} else if data == nil {
		t.Errorf("File not found")
	} else {
		if tm, err := data.Html(); err != nil {
			t.Error(err)
			return
		} else {
			vars := map[string]interface{}{
				"CompanyName": "あいうえお",
				"Title":       "テスト",
			}
			w := &bytes.Buffer{}
			if err = tm.Execute(w, vars); err != nil {
				t.Error(err)
			}
			if !strings.Contains(w.String(), "%テスト%") {
				t.Errorf("Not contains expected value.")
			}
			fmt.Println(w.String())
		}
	}
}
