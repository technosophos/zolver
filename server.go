package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/web"
	"gopkg.in/yaml.v2"
)

func main() {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("@startup", "Start app").
		Does(ParseYaml, "conf")

	reg.Route("GET **", "Handle inbound requests").
		Does(Resolve, "res")

	router.HandleRequest("@startup", cxt, false)
	web.Serve(reg, router, cxt)

	// Listen
}

func ParseYaml(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	var zconf ZolverYaml
	zconf = make(map[string]ZolverRoute)

	f := p.Get("file", "./zolver.yaml").(string)
	// Read YAML file
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return zconf, fmt.Errorf("Failed to open zolver.yaml: %s\n", err)
	}

	if err := yaml.Unmarshal(data, zconf); err != nil {
		return zconf, fmt.Errorf("Invalid YAML: %s\n", err)
	}

	fmt.Printf("Zolver: %v\n", zconf)
	return zconf, nil
}

func Resolve(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	return nil, nil
}

type ZolverYaml map[string]ZolverRoute

type ZolverRoute struct {
	Default bool
	To      string
}
