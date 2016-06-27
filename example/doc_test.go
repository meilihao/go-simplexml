package example

import (
	"testing"

	. "github.com/meilihao/go-simplexml"
)

func TestBuildDocument(t *testing.T) {
	doc := NewDocument()

	root := NewElement("root", doc)
	doc.Root = root

	node1 := NewElement("node1", doc)
	root.AddChild(node1)
	node2 := NewElement("node2", doc)
	root.AddChild(node2)

	subnode2 := NewElement("sub", doc)
	node2.AddChild(subnode2)

	expected := `<root>
	<node1 />
	<node2>
		<sub />
	</node2>
</root>`

	if doc.String("\t") != expected {
		t.Errorf("want (%s) get (%s)", doc.String("\t"), expected)
	}
}

func TestDecodeXml(t *testing.T) {
	s := []string{
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<root xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	<foo>z</foo>
	<foo>a</foo>
	<foo>
		<test>
			<foo>a</foo>
		</test>
	</foo>
	<foo />
	<remark><![CDATA[this >< is a test]]></remark>
</root>`,
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<env:Envelope xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns="http://localhost/foo/" xmlns:bar="http://localhost/bar1/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<env:Body xmlns:bar="http://localhost/bar2/">
		<Fizz where="test">Buzz</Fizz>
		<bar:abc bar:where="test">123</bar:abc>
	</env:Body>
	<bar:def xsi:nil="true"></bar:def>
</env:Envelope>`,
		`<root xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.w3.org/2001/XMLSchema http://www.w3.org/2001/XMLSchema.xsd">
	<books>
		<book>
			<Name xsi:type="xs:string">C#入门到精通</Name>
			<Price xsi:type="xs:decimal">25.6</Price>
		</book>
	</books>
</root>`,
	}

	for i, v := range s {
		x, err := NewXml([]byte(v))
		if err != nil {
			t.Errorf("`%s` decode want (%v) get (%v)\n", v, nil, err)
		}
		if i == 1 {
			x.CloseTagNicely = false
		}
		if x.String("\t") != v {
			t.Errorf("`xml rebuild want (%v) get (%v)\n", v, x.String("\t"))
		}
	}
}
