package main

import (
	"regexp"
)

type MarkdownSubstitution struct {
	regex        *regexp.Regexp
	substitution []byte
}

var MarkdownContentRegexes = []MarkdownSubstitution{
	{regexp.MustCompile(`\!\[([^\]]+)\]\(([^\)]+)\)`), []byte(`<img src="$2" alt="$1" />`)},
	{regexp.MustCompile(`\[([^\]]+)\]\((\S*)\)`), []byte(`<a href="$2" title="$4">$1</a>`)},
	{regexp.MustCompile(`\[([^\]]+)\]\((\S*)(?:\s+)?(\"(.*)\")?\)`), []byte(`<a href="$2" title="$4">$1</a>`)},
}

var MarkdownStyleRegexes = []MarkdownSubstitution{
	{regexp.MustCompile(`_{2}(.*?)_{2}`), []byte(`<u>$1</u>`)},
	{regexp.MustCompile(`~{2}(.*?)~{2}`), []byte(`<s>$1</s>`)},
	{regexp.MustCompile(`\*{2}(.*?)\*{2}`), []byte(`<strong>$1</strong>`)},
	{regexp.MustCompile(`\*{1}(.*?)\*{1}`), []byte(`<em>$1</em>`)},
}

var MarkdownCodeRegexes = []MarkdownSubstitution{
	{regexp.MustCompile("^```(.*?)```"), []byte(`<pre>$1</pre>`)},
	{regexp.MustCompile("`(.*?)`"), []byte(`<code>$1</code>`)},
}

var MarkdownHeaderRegexes = []MarkdownSubstitution{
	{regexp.MustCompile(`[#]{6}\s+(.+)`), []byte(`<h6>$1</h6>`)},
	{regexp.MustCompile(`[#]{5}\s+(.+)`), []byte(`<h5>$1</h5>`)},
	{regexp.MustCompile(`[#]{4}\s+(.+)`), []byte(`<h4>$1</h4>`)},
	{regexp.MustCompile(`[#]{3}\s+(.+)`), []byte(`<h3>$1</h3>`)},
	{regexp.MustCompile(`[#]{2}\s+(.+)`), []byte(`<h2>$1</h2>`)},
	{regexp.MustCompile(`[#]{1}\s+(.+)`), []byte(`<h1>$1</h1>`)},
}

var MarkdownRegexes = [][]MarkdownSubstitution{
	MarkdownHeaderRegexes,
	MarkdownCodeRegexes,
	MarkdownStyleRegexes,
	MarkdownContentRegexes,
}

func Parse(content []byte) string {
	for _, regexclass := range MarkdownRegexes {
		for _, mdsub := range regexclass {
			content = mdsub.regex.ReplaceAll(content, mdsub.substitution)
		}
	}
	return string(content)
}
