package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type module struct {
	*pgs.ModuleBase

	goContext pgsgo.Context
	tpl       *template.Template
}

func New() pgs.Module { return &module{ModuleBase: &pgs.ModuleBase{}} }

func (m *module) Name() string {
	return "fuzz"
}

func (m *module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.goContext = pgsgo.InitContext(c.Parameters())

	tpl, err := parseTemplates(map[string]interface{}{
		"Name":        m.goContext.Name,
		"OneOfOption": m.goContext.OneofOption,
		"ImportPath":  m.goContext.ImportPath,
	})
	if err != nil {
		m.Fail(err)
	}

	m.tpl = tpl
}

func (m *module) Execute(files map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	var p pgs.Package
	for _, f := range files {
		if len(f.AllMessages()) == 0 {
			continue
		}

		p = f.Package()

		fp := m.goContext.OutputPath(f)
		of := fp.Dir().Push("fuzz").Push(fp.BaseName()).SetExt(".fz.go").String()
		m.AddGeneratorTemplateFile(of, m.tpl.Lookup("header"), f)
		for _, msg := range f.AllMessages() {
			m.AddGeneratorTemplateAppend(of, m.tpl.Lookup("message"), msg)

			for _, o := range msg.OneOfFields() {
				m.AddGeneratorTemplateAppend(of, m.tpl.Lookup("oneof"), o)
			}
		}
	}

	if p != nil {
		m.AddGeneratorTemplateFile(m.goContext.OutputPath(p.Files()[0]).Dir().Push("fuzz").Push("fuzz").SetExt(".go").String(), m.tpl.Lookup("fuzz"), p)
	}

	return m.Artifacts()
}

//go:embed templates/*.tmpl
var tfs embed.FS

func parseTemplates(fm template.FuncMap) (*template.Template, error) {
	tpl := template.New("fuzz").Funcs(fm)

	if err := fs.WalkDir(tfs, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := tfs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		if _, err := tpl.New(strings.TrimSuffix(info.Name(), ".tmpl")).Parse(string(content)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tpl, nil
}
