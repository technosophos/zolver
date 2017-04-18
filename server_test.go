package main

import (
	"testing"

	"github.com/Masterminds/cookoo"
)

func TestParseYaml(t *testing.T) {
	vals := map[string]interface{}{
		"file": "testdata/zolver.yaml",
	}
	c := cookoo.NewContext()
	ps := cookoo.NewParamsWithValues(vals)

	confIface, err := ParseYaml(c, ps)
	if err != nil {
		t.Fatal("Failed to parse YAML")
	}

	conf, ok := confIface.(ZolverYaml)
	if !ok {
		t.Fatal("Failed to get interface as ZolverYaml")
	}
	r, ok := conf["gh"]
	if !ok {
		t.Error("Could not find route for gh")
	}
	if r.To != "https://github.com" {
		t.Errorf("Expected https://github.com, got %q", r.To)
	}
}

func TestBuildTemplates(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()
	reg.Route("test", "Test").
		Does(ParseYaml, "conf").Using("file").From("testadata/zolver.yaml").
		Does(BuildTemplates, "tpl").
		Using("config").From("cxt:conf")

	if err := router.HandleRequest("test", cxt, false); err != nil {
		t.Fatal(err)
	}
}
