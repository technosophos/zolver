package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/web"
	"github.com/Masterminds/sprig"
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

	newurl, err := destination(req, &route)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return nil, nil
	}

	c.Logf("info", "Rerouting %s to %s", req.URL.String(), newurl)

	//res.Header().Add("location", route.To)
	http.Redirect(res, req, newurl, 302)
	return nil, nil
}

// destination builds the redirect URL.
func destination(req *http.Request, route *ZolverRoute) (string, error) {
	oldurl := req.URL

	if val, ok := route.Short[oldurl.Path[1:]]; ok {
		return val, nil
	}

	if len(route.Tpl) > 0 {
		var err error
		newurl, err := doTemplate(route.Tpl, oldurl)
		if err != nil {
			return oldurl.String(), err
		}
		fmt.Printf("Rendered TPL as %s", route.To)
		return newurl, nil
	}

	newurl, err := url.Parse(route.To)
	if err != nil {
		return newurl.String(), err
	}

	newurl.Path = oldurl.Path
	newurl.RawQuery = oldurl.RawQuery
	newurl.Fragment = oldurl.Fragment

	return newurl.String(), nil
}

func doTemplate(tpl string, oldurl *url.URL) (string, error) {
	t, err := template.New("url").Parse(tpl)
	if err != nil {
		return tpl, err
	}
	var b bytes.Buffer
	err = t.Execute(&b, oldurl)

	return b.String(), err
}

type ZolverYaml map[string]ZolverRoute

type ZolverRoute struct {
	Default bool
	To      string
	Tpl     string
	Short   map[string]string
}
