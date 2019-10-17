package main

import "testing"

func TestFull(t *testing.T) {
	tables := []struct {
		md   []byte
		html string
	}{
		{
			[]byte(
				"# Title\n" +
					"## Subtitle\n" +
					"`printf`:__**hello** *world*__\n" +
					"![AltText](http://url.com)\n" +
					"~~**Visit my website**~~: [AltText](http://url.com)"),
			"<h1>Title</h1>\n" +
				"<h2>Subtitle</h2>\n" +
				"<code>printf</code>:<u><strong>hello</strong> <em>world</em></u>\n" +
				"<img src=\"http://url.com\" alt=\"AltText\" />\n" +
				"<s><strong>Visit my website</strong></s>: <a href=\"http://url.com\" title=\"\">AltText</a>",
		},
	}

	for _, table := range tables {
		rendered := Parse(table.md)
		if rendered != table.html {
			t.Errorf("Render of \"%s\" was incorrect, got: \"%s\", want: \"%s\".", table.md, rendered, table.html)
		}
	}
}

func TestContent(t *testing.T) {
	tables := []struct {
		md   []byte
		html string
	}{
		{[]byte(`![AltText](http://url.com)`), `<img src="http://url.com" alt="AltText" />`},
		{[]byte(`[AltText](http://url.com)`), `<a href="http://url.com" title="">AltText</a>`},
		{[]byte(`[AltText](http://url.com "Desc")`), `<a href="http://url.com" title="Desc">AltText</a>`},
	}

	for _, table := range tables {
		rendered := Parse(table.md)
		if rendered != table.html {
			t.Errorf("Render of \"%s\" was incorrect, got: \"%s\", want: \"%s\".", table.md, rendered, table.html)
		}
	}
}

func TestStyle(t *testing.T) {
	tables := []struct {
		md   []byte
		html string
	}{
		{[]byte("~~Test~~"), "<s>Test</s>"},
		{[]byte("__Test__"), "<u>Test</u>"},
		{[]byte("*Test*"), "<em>Test</em>"},
		{[]byte("**Test**"), "<strong>Test</strong>"},
	}

	for _, table := range tables {
		rendered := Parse(table.md)
		if rendered != table.html {
			t.Errorf("Render of \"%s\" was incorrect, got: \"%s\", want: \"%s\".", table.md, rendered, table.html)
		}
	}
}

func TestCode(t *testing.T) {
	tables := []struct {
		md   []byte
		html string
	}{
		{[]byte("`printf(\"Hello, world!\");`"), "<code>printf(\"Hello, world!\");</code>"},
	}

	for _, table := range tables {
		rendered := Parse(table.md)
		if rendered != table.html {
			t.Errorf("Render of \"%s\" was incorrect, got: \"%s\", want: \"%s\".", table.md, rendered, table.html)
		}
	}
}

func TestHeader(t *testing.T) {
	tables := []struct {
		md   []byte
		html string
	}{
		{[]byte("# Header"), "<h1>Header</h1>"},
		{[]byte("## Header"), "<h2>Header</h2>"},
		{[]byte("### Header"), "<h3>Header</h3>"},
		{[]byte("#### Header"), "<h4>Header</h4>"},
		{[]byte("##### Header"), "<h5>Header</h5>"},
		{[]byte("###### Header"), "<h6>Header</h6>"},
	}

	for _, table := range tables {
		rendered := Parse(table.md)
		if rendered != table.html {
			t.Errorf("Render of \"%s\" was incorrect, got: \"%s\", want: \"%s\".", table.md, rendered, table.html)
		}
	}
}
