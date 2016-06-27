package simplexml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

func Version() string {
	return "0.0.1"
}

func NewFromReader(r io.Reader) (*Document, error) {
	return decode(r)
}

func NewXml(raw []byte) (*Document, error) {
	return decode(bytes.NewBuffer(raw))
}

// must have root tag
func decode(r io.Reader) (*Document, error) {
	d := xml.NewDecoder(r)

	var start xml.StartElement
	var tree []*Element
	var e *Element
	doc := NewDocument()

	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch token := t.(type) {
		case xml.StartElement:
			start = token.Copy()

			n := len(tree)
			if n == 0 {
				e = &Element{
					XmlDocument: doc,
				}

				doc.Root = e
			} else {
				e = &Element{
					XmlDocument: doc,
					Parent:      tree[n-1], // because tree[n-1] tag is open
				}

				tree[n-1].Children = append(tree[n-1].Children, e)
			}

			e.Name = start.Name
			e.Attrs = start.Attr

			// handle namespace
			for _, v := range e.Attrs {
				if v.Name.Space == XmlnsStr || v.Name.Local == XmlnsStr {
					if v.Name.Space == XmlnsStr {
						doc.NameSpace[v.Value] = v.Name.Local
					} else if v.Name.Local == XmlnsStr && v.Name.Space == "" {
						// go1.6
						// other ns : {Name:{Space:xmlns Local:env} Value:http://schemas.xmlsoap.org/soap/envelope/}
						// default ns : {Name:{Space: Local:xmlns} Value:http://localhost/foo/}
						// why default namespace is not "{Name:{Space:xmlns Local:} Value:http://localhost/foo/}"?
						// Default Namespace
						doc.NameSpace[v.Value] = ""
					}
				}
			}

			tree = append(tree, e)
		case xml.CharData:
			// skip whitespace
			if str := strings.TrimSpace(string(token)); str != "" {
				e.Value = str
			}
		case xml.EndElement:
			e = nil
			tree = tree[:len(tree)-1]
			start = xml.StartElement{}
		case xml.ProcInst:
			//XML Declaration
			doc.Declaration = fmt.Sprintf("<?%s %s?>", token.Target, string(token.Inst))
		case xml.Directive:
		case xml.Comment:
		default:
			return nil, fmt.Errorf("invalid token(%v)", token)
		}
	}

	if len(tree) != 0 {
		return nil, errors.New("invalid document")
	}

	return doc, nil
}
