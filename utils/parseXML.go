package utils

import (
	"fmt"
	"github.com/tinyhubs/tinydom"
	"strings"
)

func ParseXML(str, str1, str2, str3 string) string {
	doc, err := tinydom.LoadDocument(strings.NewReader(str))
	if err != nil {
		fmt.Printf("XML解析失败，err：%s\n", err.Error())
	}
	elem := doc.FirstChildElement(str1).FirstChildElement(str2).FirstChildElement(str3)
	return elem.Text()
}
