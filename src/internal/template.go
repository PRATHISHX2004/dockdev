package internal

import (
	"os"
	"text/template"
)

func RenderTemplate(templatePath string, destPath string, data interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, data)
}
