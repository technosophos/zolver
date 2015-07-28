package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/web"
	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

func main() {
	reg, router, cxt := cookoo.Cookoo()

	cxt.Put("server.Address", ":80")
	//cxt.Put("server.Address", ":8080")

	if len(os.Args) > 1 {
		first := os.Args[1]
		cxt.Put("file", first)
	}

	reg.Route("@startup", "Start app").
		Does(ParseYaml, "conf").
		Using("file").From("cxt:file").
		Does(BuildTemplates, "tpl").
		Using("config").From("cxt:conf")

	reg.Route("GET /**", "Handle inbound requests").
		Does(Resolve, "res").
		Using("config").From("cxt:conf").
		Using("tpl").From("cxt:tpl")

	router.HandleRequest("@startup", cxt, false)

	// Periodically re-read the file.
	go func() {
		t := time.NewTicker(5 * time.Minute)
		for range t.C {
			router.HandleRequest("@startup", cxt, false)
		}
	}()

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

	//fmt.Printf("Zolver: %v\n", zconf)
	domains := make([]string, len(zconf))
	i := 0
	for d, _ := range zconf {
		domains[i] = d
		i++
	}
	c.Logf("info", "Your /etc/hosts file should have a line like this:")
	c.Logf("info", "\t127.0.0.1\tlocahost %s", strings.Join(domains, " "))
	return zconf, nil
}

func BuildTemplates(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	domains := p.Get("config", nil).(ZolverYaml)

	// We could probably refactor this to use one master template with
	// several sub-templates.
	templates := make(map[string]*template.Template)
	extras := sprig.TxtFuncMap()

	for dom, route := range domains {
		if len(route.Tpl) > 0 {
			c.Logf("info", "Compiling template: '%s'", route.Tpl)
			t, err := template.New(dom).Funcs(extras).Parse(route.Tpl)
			if err != nil {
				c.Logf("error", "Error compiling template: %s", err)
				return templates, err
			}
			templates[route.Tpl] = t
		}
	}
	return templates, nil
}

func Resolve(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {

	req := c.Get("http.Request", nil).(*http.Request)
	res := c.Get("http.ResponseWriter", nil).(http.ResponseWriter)
	cfg := p.Get("config", nil).(ZolverYaml)
	tpls := p.Get("tpl", nil).(map[string]*template.Template)

	host := strings.SplitN(req.Host, ":", 2)

	c.Logf("info", "Requested host: %s", host[0])
	route, ok := cfg[host[0]]
	if !ok {
		http.NotFound(res, req)
		return nil, nil
	}

	newurl, err := destination(req, &route, tpls)
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
func destination(req *http.Request, route *ZolverRoute, tpls map[string]*template.Template) (string, error) {
	oldurl := &URL{req.URL}

	if val, ok := route.Short[oldurl.Path[1:]]; ok {
		return val, nil
	}

	if len(route.Tpl) > 0 {
		var err error
		tpl := tpls[route.Tpl]
		newurl, err := doTemplate(tpl, oldurl)
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

type URL struct {
	*url.URL
}

func (u *URL) Part(index int) string {
	parts := strings.Split(u.Path, "/")
	if len(parts) < index+1 {
		return ""
	}
	return parts[index]
}

func doTemplate(t *template.Template, oldurl *URL) (string, error) {
	var b bytes.Buffer
	err := t.Execute(&b, oldurl)

	return b.String(), err
}

type ZolverYaml map[string]ZolverRoute

type ZolverRoute struct {
	Default bool
	To      string
	Tpl     string
	Short   map[string]string
}
