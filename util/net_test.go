package util

import (
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	list := GetLocalIP()
	t.Log(list)
}
