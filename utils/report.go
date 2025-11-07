package utils

import (
	"io/ioutil"
	"strings"
)

func InsertCss(reportTemplateFile, cssFile string, index int) error {
	css, err := ioutil.ReadFile(cssFile)
	if err != nil {
		return err
	}

	html, err := ioutil.ReadFile(reportTemplateFile)
	if err != nil {
		return err
	}

	if !strings.Contains(strings.ReplaceAll(strings.ReplaceAll(string(html), " ", ""), "\t", ""), strings.ReplaceAll(strings.ReplaceAll(string(css), " ", ""), "\t", "")) {
		if err := InsertStringToFile(reportTemplateFile, "<style>"+string(css)+"</style> \n", index); err != nil {
			return err
		}
	}

	return nil
}
