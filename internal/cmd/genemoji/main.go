// Executable genemoji generates a table of emojis from an Emoji specification
// emoji-data.txt file.
package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type emojiDatum struct {
	Value            string
	Status           string
	Name             string
	Introduced       string
	FullyQualifiesAs string
}

var tmpl = template.Must(template.New("").Parse(`package emoji

var emojis = map[string]Emoji{
{{- range . }}
	{{ .Value }}: { {{ .Name }}, {{ .Status }}, {{ .Introduced }}, {{ .FullyQualifiesAs }} },
{{- end }}
}
`))

var (
	data = flag.String("data", "", "input file")
	out  = flag.String("out", "", "output file")
)

func main() {
	flag.Parse()

	f, err := os.Open(*data)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var emojiData []emojiDatum
	fullyQualifiedByName := map[string]int{} // used to track the (unique) fully-qualified version of a given emoji

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()

		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		// Here's an example of the most "complicated" sort of entry in the
		// emoji-test.txt dataset:
		//
		// 1F469 1F3FB 200D 2764 200D 1F48B 200D 1F469 1F3FB      ; minimally-qualified # 👩🏻‍❤‍💋‍👩🏻 E13.1 kiss: woman, woman, light skin tone
		//
		// The format here is: (where columns are zero-indexed)
		//
		// - Space-delimited hex-encoded codepoints
		// - Space up to column 55
		// - Semi-colon
		// - Space
		// - Status, one of "component", "unqualified", "minimally-qualified", "fully-qualified"
		// - Space up to column 77
		// - Pound sign, space
		// - The emoji, consisting of as many codepoints as indicated by the first section
		// - Space, "E"
		// - The version of Emoji the emoji was introduced
		// - Space
		// - Name
		//
		// The following code implements parsing such a format. It's very
		// possible that future versions of the Emoji standard will change the
		// exact values of the column alignments.
		col1 := 55
		col2 := 77
		extraAfterCol1 := 2 // used to trim the semi-colon and space after col1
		extraAfterCol2 := 3 // used to trim the space before and after emoji + "E" before version

		codepoints := strings.TrimSpace(l[0:col1])
		status := strings.TrimSpace(l[col1+extraAfterCol1 : col2])

		var runes []rune
		for _, cp := range strings.Split(strings.TrimSpace(codepoints), " ") {
			n, err := strconv.ParseInt(cp, 16, 32)
			if err != nil {
				panic(err)
			}

			runes = append(runes, rune(n))
		}

		trailing := string([]rune(l)[col2+extraAfterCol2+len(runes):])
		trailingParts := strings.SplitN(trailing, " ", 2)
		introduced := trailingParts[0]
		name := trailingParts[1]

		emojiData = append(emojiData, emojiDatum{
			Value: strconv.Quote(string(runes)),
			Status: map[string]string{
				"component":           "Component",
				"fully-qualified":     "FullyQualified",
				"minimally-qualified": "MinimallyQualified",
				"unqualified":         "Unqualified",
			}[status],
			Introduced:       strconv.Quote(introduced),
			Name:             strconv.Quote(name),
			FullyQualifiesAs: strconv.Quote(""), // default value
		})

		if status == "fully-qualified" {
			fullyQualifiedByName[strconv.Quote(name)] = len(emojiData) - 1
		}
	}

	for i, d := range emojiData {
		if j, ok := fullyQualifiedByName[d.Name]; ok {
			emojiData[i].FullyQualifiesAs = emojiData[j].Value
		}
	}

	f, err = os.Create(*out)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(f, emojiData); err != nil {
		panic(err)
	}
}
