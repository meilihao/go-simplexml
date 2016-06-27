package simplexml

import (
	"fmt"
	"strings"
)

func getIndent(indent []string) string {
	if len(indent) > 0 && indent[0] != "" {
		return indent[0]
	}

	return ""
}

func getNs(e *Element, name string) string {
	if name == XmlnsStr {
		// namespace define
		return XmlnsStr
	}
	return e.XmlDocument.NameSpace[name]
}

func isCDATA(s string) bool {
	if strings.Contains(s, "<") || strings.Contains(s, ">") {
		return true
	}
	return false
}

func Value2Xml(s string) string {
	if isCDATA(s) {
		return fmt.Sprintf("<![CDATA[%s]]>", s)
	} else {
		return s
	}
}
