// +build ignore

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"text/template"
)

//
// docker run --rm -v $(pwd):/tmp alpine:3.6 /bin/sh -c "apk add -U mailcap; mv /etc/mime.types /tmp"
//

type mime struct {
	Name string
	Exts []string
}

func main() {
	in, err := os.Open("mime.types")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var mimes []mime

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "#") || len(txt) == 0 {
			continue
		}
		parts := strings.Fields(txt)
		if len(parts) == 1 {
			continue
		}
		mimes = append(mimes, mime{parts[0], parts[1:]})
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("mime_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	template.Must(
		template.New("_").Parse(templ),
	).Execute(out, mimes)
}

var templ = `package mimetype

import "mime"

func init() {
{{- range . -}}
{{ $name := .Name -}}
{{- range .Exts }}
	mime.AddExtensionType({{ printf "\".%s\"" . }}, {{ printf "%q" $name }})
{{- end }}
{{- end }}
}
`
