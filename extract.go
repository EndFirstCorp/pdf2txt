package pdf2txt

import "io"
import "io/ioutil"
import "bytes"

// extract the compressed streams into text files for debugging
func extract(r io.Reader) {
	uncategorized := make(map[string]*object)
	contents := []string{}
	fonts := make(map[string]*font)
	toUnicode := []string{}

	tchan := make(chan interface{}, 15)
	go Tokenize(newBufReader(r), tchan)

	for t := range tchan {
		switch v := t.(type) {
		case *object:
			oType := v.name("/Type")
			switch oType {
			case "/Page":
				page := getPage(v)
				contents = append(contents, page.Contents...)

			case "/Font":
				if _, ok := fonts[v.refString]; !ok {
					font := getFont(v)
					fonts[v.refString] = font
					if font.ToUnicode != "" {
						toUnicode = append(toUnicode, font.ToUnicode)
					}
				}

			default:
				uncategorized[v.refString] = v
			}
		}
	}

	for i := range toUnicode {
		ref := toUnicode[i]
		var buf bytes.Buffer
		buf.ReadFrom(uncategorized[ref].stream)
		ioutil.WriteFile("toUnicode "+ref+".txt", buf.Bytes(), 644)
	}

	for i := range contents {
		ref := contents[i]
		var buf bytes.Buffer
		buf.ReadFrom(uncategorized[ref].stream)
		ioutil.WriteFile("contents "+ref+".txt", buf.Bytes(), 644)
	}
}
