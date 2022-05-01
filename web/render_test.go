package web

import (
	"fmt"
	"regexp"
	"testing"
)

func TestHtml(t *testing.T) {

	dirReg := regexp.MustCompile(`<a[\s\S]+?</a>`)
	fmt.Println(dirReg.FindAllString(DirHtmlTemplate, -1))

}
