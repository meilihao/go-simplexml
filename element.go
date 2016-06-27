package simplexml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

// Xml Element
type Element struct {
	Name        xml.Name
	Attrs       []xml.Attr
	Parent      *Element
	Children    []*Element
	Value       string
	XmlDocument *Document
}

// need doc!=nil
func NewElement(name string, doc *Document) *Element {
	if doc == nil {
		return nil
	}

	e := &Element{Name: xml.Name{Local: name}}
	e.Children = make([]*Element, 0, 5)
	e.Attrs = make([]xml.Attr, 0, 3)
	e.XmlDocument = doc
	return e
}

func (e *Element) AddChild(child *Element) *Element {
	if child.Parent != nil {
		child.Parent.RemoveChild(child)
	}

	child.Parent = e
	e.Children = append(e.Children, child)

	return e
}

func (e *Element) RemoveChild(child *Element) *Element {
	p := -1
	for i, v := range e.Children {
		if v == child {
			p = i
			break
		}
	}

	if p == -1 {
		return e
	}

	copy(e.Children[p:], e.Children[p+1:])
	e.Children = e.Children[0 : len(e.Children)-1]
	child.Parent = nil

	return e
}

// e.XmlDocument!=nil before call String()
func (e *Element) String(indent ...string) string {
	if e.XmlDocument == nil {
		return ""
	}

	buf := new(bytes.Buffer)

	e.Bytes(buf, getIndent(indent), 0)

	if n := buf.Len(); n > 0 {
		buf.Truncate(n - 1)
	}
	return buf.String()
}

func (e *Element) Bytes(buf *bytes.Buffer, indent string, level int) {
	indentStr, isIndent := "", indent != ""
	tag := ""

	if isIndent {
		for i := 0; i < level; i++ {
			indentStr += indent
		}
	}

	if e.Name.Local != "" {
		if e.Name.Space != "" {
			tag = fmt.Sprintf("%s:%s", getNs(e, e.Name.Space), e.Name.Local)
			if strings.HasPrefix(tag, ":") {
				// default namespace
				tag = strings.TrimLeft(tag, ":")
			}
		} else {
			tag = fmt.Sprintf("%s", e.Name.Local)
		}

		fmt.Fprintf(buf, "%s<%s", indentStr, tag)
	}

	for _, v := range e.Attrs {
		if v.Name.Space != "" {
			fmt.Fprintf(buf, ` %s:%s="%s"`, getNs(e, v.Name.Space), v.Name.Local, v.Value)
		} else {
			fmt.Fprintf(buf, ` %s="%s"`, v.Name.Local, v.Value)
		}
	}

	if len(e.Children) > 0 {
		if isIndent {
			fmt.Fprintf(buf, ">\n")
		} else {
			fmt.Fprintf(buf, ">")
		}

		for _, v := range e.Children {
			v.Bytes(buf, indent, level+1)
		}

		if isIndent {
			fmt.Fprintf(buf, "%s</%s>\n", indentStr, tag)
		} else {
			fmt.Fprintf(buf, "%s</%s>", indentStr, tag)
		}
	} else {
		if e.Value == "" {
			if e.XmlDocument.CloseTagNicely {
				fmt.Fprintf(buf, " />")
			} else {
				fmt.Fprintf(buf, "></%s>", tag)
			}
		} else {
			fmt.Fprintf(buf, ">%s</%s>", Value2Xml(e.Value), tag)
		}

		if isIndent {
			fmt.Fprintf(buf, "\n")
		}
	}
}
