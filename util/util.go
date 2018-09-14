package util

import (
	"os/exec"
	"strings"
)

func NewUUID() string {
	if o, e := exec.Command("uuidgen").Output(); e != nil {
		return ""
	} else {
		return strings.Trim(string(o), "\n")
	}
}

func Shorten(uuid string) string {
	return "#" + uuid[:4]
}

func IsExisted(items []int, i int) bool {
	for _, v := range items {
		if i == v {
			return true
		}
	}
	return false
}
