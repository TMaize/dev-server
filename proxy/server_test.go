package proxy

import (
	"fmt"
	"net/url"
	"testing"
)

func TestName(t *testing.T) {
	parse, err := url.Parse("https://uptream.com:8080")
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(parse.Scheme)
	fmt.Println(parse.Host)
	fmt.Println(parse.Port())
}