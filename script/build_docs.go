package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const CMD_LIST_HEADER_LINES = 2
const CMD_LIST_FOOTER_LINES = 1
const ASSET_PATH = "../app/src/main/assets/"
const CATALOG_FILE = "../app/src/main/res/values/catalog.xml"

func main() {
	cmds := commandList()

	for _, cmd := range cmds {
		extractHelp(cmd)
	}
	genCatalog(cmds)

}

func commandList() []string {
	var cmds []string

	out, err := exec.Command("fossil", "help").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	lines = lines[CMD_LIST_HEADER_LINES : len(lines)-(CMD_LIST_FOOTER_LINES+1)]

	for _, s := range lines {
		cmds = append(cmds, strings.Fields(s)...)
	}

	return cmds
}

func extractHelp(cmd string) {
	out, err := exec.Command("fossil", "help", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	text := string(out)

	f, err := os.Create(filepath.Join(ASSET_PATH, cmd+".html"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("<pre>" + text + "</pre>")
}

func genCatalog(cmds []string) {
	template := template.Must(template.New("").Parse(`
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string-array name="catalog">
{{ range . }}        <item>{{.}}</item>
{{ end }}
    </string-array>
</resources>`))

	f, err := os.Create(CATALOG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = template.Execute(f, cmds)

}
