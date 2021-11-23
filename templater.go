package templating

import (
	"html/template"
	"net/http"
)

type Templater interface {
	GetTemplate(path string) *template.Template
	GetRefs(path string) Refs
	GetPageReqs(path string) ??
}

type Page struct {
	Template *template.Template
	Data     interface{}
}

type Renderer struct {
	Templates  Templater
	Err500Page *Page
}

// TODO: need a convenient way of tracking which templates rely on other
// templates so we can make sure they all get rendered.
// Templates can be associated with each other using AddParseTree, passing
// otherTemplate.Tree. So we can do something like
/*
t := template.New("base")
for _, tmpl := range templates {
	t, err = t.AddParseTree(tmpl.Name(), tmpl.Tree)
	if err != nil {
		// do something
	}
}
*/
// then render t

func (rend Renderer) Render(w http.ResponseWriter, r *http.Request, status int, template string, data interface{}) {
	tmpl := rend.Templates.Get(template)
	if tmpl == nil {
		// TODO: log this error
		// TODO: set tmpl to default error template
		data = nil
		status = http.StatusInternalServerError
		if rend.Err500Page != nil {
			tmpl = rend.Err500Page.Template
			data = rend.Err500Page.Data
		}
	}
	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		// TODO: log this error, do something?
	}

}
