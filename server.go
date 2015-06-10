package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/web"
	"gopkg.in/yaml.v2"
)

func main() {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("@startup", "Start app").
		Does(ParseYaml, "conf")

	reg.Route("GET /**", "Handle inbound requests").
		Does(Resolve, "res").Using("config").From("cxt:conf")

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

	req := c.Get("http.Request", nil).(*http.Request)
	res := c.Get("http.ResponseWriter", nil).(http.ResponseWriter)
	cfg := p.Get("config", nil).(ZolverYaml)

	host := strings.SplitN(req.Host, ":", 2)

	c.Logf("info", "Requested host: %s", host[0])
	route, ok := cfg[host[0]]
	if !ok {
		http.NotFound(res, req)
		return nil, nil
	}
	//res.Header().Add("location", route.To)
	http.Redirect(res, req, route.To, 302)
	return nil, nil
}

type ZolverYaml map[string]ZolverRoute

type ZolverRoute struct {
	Default bool
	To      string
}
