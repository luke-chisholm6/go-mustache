package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func splitKV(keyValue string) (key string, value string, err error) {
	keyValueArray := strings.Split(keyValue, "=")

	if len(keyValueArray) != 2 {
		err = errors.New("bad formatting kv pair")
		return
	}

	key, value = keyValueArray[0], keyValueArray[1]

	return
}

func getTemplateContext(kvSlice []string) (map[string]string, error) {
	keyValueMap := make(map[string]string)

	for _, keyValue := range kvSlice {
		key, value, err := splitKV(keyValue)
		if err != nil {
			return keyValueMap, err
		}
		keyValueMap[key] = value
	}

	return keyValueMap, nil
}

func compileTemplate(reader io.Reader) (*template.Template, error) {
	bytes, err := ioutil.ReadAll(reader)
	switch {
	case err != nil:
		return nil, err
	case len(bytes) <= 0:
		return nil, errors.New("no template context given")
	}

	inputTemplate := string(bytes)
	compiledTemplate, err := template.New("").Parse(inputTemplate)
	if err != nil {
		return nil, err
	}

	return compiledTemplate, nil
}

func render(tmpl *template.Template, context map[string]string, writer io.Writer) error {
	return tmpl.Execute(writer, context)
}

func run(inputTemplate io.Reader, inputContext []string, output io.Writer) error {
	compiledTemplate, err := compileTemplate(inputTemplate)
	if err != nil {
		return err
	}

	templateContext, err := getTemplateContext(inputContext)
	if err != nil {
		return err
	}

	if err := render(compiledTemplate, templateContext, output); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(os.Stdin, os.Args[1:], os.Stdout); err != nil {
		panic(err)
	}
}
