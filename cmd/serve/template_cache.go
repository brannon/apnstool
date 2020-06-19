// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"fmt"
	"text/template"
)

type TemplateCache struct {
	templates map[string]*template.Template
	funcs     template.FuncMap
}

func (cache *TemplateCache) Lookup(name string) *template.Template {
	if cache.templates != nil {
		tmpl, found := cache.templates[name]
		if found {
			return tmpl
		}
	}
	return nil
}

func (cache *TemplateCache) MustParse(name string, html string) *template.Template {
	if cache.templates == nil {
		cache.templates = make(map[string]*template.Template)
	}

	tmpl := template.New(name)

	if cache.funcs != nil {
		tmpl.Funcs(cache.funcs)
	}

	_, err := tmpl.Parse(html)
	if err != nil {
		panic(fmt.Errorf("error parsing template: %s", err))
	}

	cache.templates[name] = tmpl
	return tmpl
}

func (cache *TemplateCache) RegisterFunc(name string, fn interface{}) {
	if cache.funcs == nil {
		cache.funcs = make(template.FuncMap)
	}

	cache.funcs[name] = fn
}
