package confluence

import (
	"regexp"
	"strings"
)

func ConvertToConfluence(markup string) string {
	markup = ConvertCodeBlock(markup)

	return markup
}

func ConvertCodeBlock(markup string) string {
	re := regexp.MustCompile(`(?s)<pre><code class="language-(.*?)">(.*?)</code></pre>`)

	replaceFunc := func(s string) string {
		matches := re.FindStringSubmatch(s)
		if len(matches) != 3 {
			return s
		}
		lang := matches[1]
		code := matches[2]
		code = strings.ReplaceAll(code, "&lt;", "<")
		code = strings.ReplaceAll(code, "&gt;", ">")
		code = strings.TrimSuffix(code, "\n")

		return "<ac:structured-macro ac:name=\"code\" ac:schema-version=\"1\"><ac:parameter ac:name=\"language\">" + lang + "</ac:parameter><ac:plain-text-body><![CDATA[" + code + "]]></ac:plain-text-body></ac:structured-macro>"
	}

	markup = re.ReplaceAllStringFunc(markup, replaceFunc)

	return markup
}

func PrependWarningMessage(markup string) string {
	markup = "<ac:structured-macro ac:name=\"note\" ac:schema-version=\"1\"><ac:rich-text-body>\n<p>This page is generated automatically!</p></ac:rich-text-body></ac:structured-macro>" + markup

	return markup
}

func AppendWarningMessage(markup string) string {
	markup = markup + "<ac:structured-macro ac:name=\"note\" ac:schema-version=\"1\"><ac:rich-text-body>\n<p>This page is generated automatically!</p></ac:rich-text-body></ac:structured-macro>"

	return markup
}

func PrependTableOfContents(markup string) string {
	markup = "<ac:structured-macro ac:name=\"toc\" ac:schema-version=\"1\" data-layout=\"default\" />" + markup

	return markup
}
