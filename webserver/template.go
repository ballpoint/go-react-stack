package webserver

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type reactComponent struct {
	Name        string
	Props       []byte
	Prerendered template.HTML
}

func prerenderedComponent(rp reactComponent) template.HTML {
	return rp.Prerendered
}

func NewReactComponent(component string, props map[string]interface{}) reactComponent {
	var (
		client http.Client
	)
	propsJson, jsonErr := json.Marshal(props)

	if jsonErr != nil {
		return reactComponent{}
	}

	root := html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			html.Attribute{Key: "data-react-component", Val: component},
			html.Attribute{Key: "data-react-props", Val: string(propsJson)},
		},
	}

	resp, renderErr := client.PostForm("http://127.0.0.1:8111", url.Values{
		"component": []string{component},
		"props":     []string{string(propsJson)},
	})

	var rootWriter bytes.Buffer

	if renderErr == nil && resp != nil && resp.StatusCode == 200 {

		defer resp.Body.Close()
		renderedNodes, componentErr := html.ParseFragment(resp.Body, &html.Node{
			Type:     html.ElementNode,
			Data:     "body",
			DataAtom: atom.Body,
		})

		if len(renderedNodes) == 1 && componentErr == nil {
			log.Println(renderedNodes)
			root.FirstChild = renderedNodes[0]
		}

	}
	html.Render(&rootWriter, &root)

	return reactComponent{
		Name:        component,
		Props:       propsJson,
		Prerendered: template.HTML(rootWriter.String()),
	}
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"react": prerenderedComponent,
	}
}
