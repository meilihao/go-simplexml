package simplexml

import (
	"bytes"
)

const (
	XmlnsStr = "xmlns"
)

// Xml Document
type Document struct {
	Root           *Element
	Declaration    string
	NameSpace      map[string]string // value -> ns
	CloseTagNicely bool
}

func NewDocument() *Document {
	return &Document{
		NameSpace:      make(map[string]string),
		CloseTagNicely: true,
	}
}

func (doc *Document) String(indent ...string) string {
	tmpIndent := getIndent(indent)
	buf := new(bytes.Buffer)

	if doc.Declaration != "" {
		buf.WriteString(doc.Declaration)
		if tmpIndent != "" {
			buf.WriteString("\n")
		}
	}

	if doc.Root != nil {
		doc.Root.Bytes(buf, tmpIndent, 0)
	}

	// remove the last "\n"
	if n := buf.Len(); n > 0 {
		buf.Truncate(n - 1)
	}
	return buf.String()
}
