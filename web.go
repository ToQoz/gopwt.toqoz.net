package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ToQoz/gopwt/translator"
	"github.com/ToQoz/rome"
)

func main() {
	router := rome.NewRouter()

	staticServer := http.FileServer(http.Dir("statics"))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticServer.ServeHTTP(w, r)
	}))
	router.Get("/app.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticServer.ServeHTTP(w, r)
	}))
	router.Get("/app.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticServer.ServeHTTP(w, r)
	}))
	router.Get("/sandbox.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := build(r.FormValue("code"), w); err != nil {
			log.Println(err.Error())
		}
	}))
	log.Fatal(http.ListenAndServe(":5000", router))
}

var fileName = "example_test.go"

func build(code string, w io.Writer) error {
	tmp, err := ioutil.TempDir("sandbox", "t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(tmp, "example_test.go"), []byte(code), 0777)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	gopath, importpath, err := translator.Translate("./" + tmp)
	if err != nil {
		return err
	}
	defer os.RemoveAll(gopath)

	dir := filepath.Join(gopath, "src", importpath)
	// example_test.go -> example.go & extract test functions
	mv(filepath.Join(dir, "example_test.go"), filepath.Join(dir, "example.go"))
	tests, err := extractTests(filepath.Join(dir, "example.go"))
	if err != nil {
		return err
	}
	// generate testmain.go
	outpath := filepath.Join(dir, "testmain.go")
	out, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := testmainTpl.Execute(out, tests); err != nil {
		return err
	}

	// sh -c for `*.go`
	// https://groups.google.com/d/msg/golang-nuts/twD7eN-c804/2YK8ZtUMz8gJ
	cmd := exec.Command("sh", "-c", "gopherjs build *.go -o testmain.js && cat testmain.js")
	cmd.Dir = dir
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func extractTests(_filepath string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), _filepath, nil, 0)
	if err != nil {
		return nil, err
	}
	tests := []string{}
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			fn := decl.(*ast.FuncDecl)
			if strings.HasPrefix(fn.Name.Name, "Test") {
				tests = append(tests, fn.Name.Name)
			}
		}
	}
	return tests, nil
}

func mv(src, dst string) error {
	if err := cp(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

var testmainTpl = template.Must(template.New("main").Parse(`package main

import (
	"testing"
)

var benchmarks = []testing.InternalBenchmark{}
var examples = []testing.InternalExample{}
var tests = []testing.InternalTest{
{{range .}}
	{
		"{{.}}",
		{{.}},
	},
{{end}}
}

func main() {
	match := func(pat string, str string) (bool, error) {
		return true, nil
	}
	testing.Main(match, tests, benchmarks, examples)
}`))
