package web

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
)

// Simple wrapper for the content that is used for fragments
type container struct {
	Content template.HTML
}

//go:embed templates
var tmplFS embed.FS

var tmpl = template.Must(
	template.New("").ParseFS(tmplFS,
		"templates/*.gohtml",
		"templates/*/*.gohtml",
	),
)

// fragment renders a named template as-is (for HTMX responses).
func fragment(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// page renders a named template wrapped in the index shell (for full loads).
func page(w http.ResponseWriter, name string, data any) {
	var inner bytes.Buffer
	if err := tmpl.ExecuteTemplate(&inner, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var out bytes.Buffer
	err := tmpl.ExecuteTemplate(&out, "index", container{
		Content: template.HTML(inner.String()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	out.WriteTo(w)
}

type pageData struct {
	Title   string
	Content template.HTML
}

// view renders `name` as a fragment for htmx requests, or wrapped in the
// index shell for normal browser loads (refresh, deep links, first visit).
func view(w http.ResponseWriter, r *http.Request, name string, data any) {
	w.Header().Add("Vary", "HX-Request")
	if r.Header.Get("HX-Request") == "true" {
		fragment(w, name, data)
		return
	}
	page(w, name, data)
}
