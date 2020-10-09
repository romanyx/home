package httprouter

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/romanyx/home/internal/medium"
)

const (
	internalServerErrorMessage = "Internal Server Error"
	indexMetaTitle             = "Roman Budnikov"
	indexMetaDescription       = "Golang developer"
	indexMetaType              = "cv"
)

// Handler holds necessary data for all routes.
type Handler struct {
	storiesFunc func() ([]medium.Story, error)
	logFunc     func(error)
	t           *template.Template
	cv          []byte
}

// NewHandler returns initialized handler.
func NewHandler(
	storiesFunc func() ([]medium.Story, error),
	logFunc func(error),
	t *template.Template,
	cv []byte,
) *Handler {
	h := Handler{
		storiesFunc: storiesFunc,
		logFunc:     logFunc,
		t:           t,
		cv:          cv,
	}

	return &h
}

// GetIndex /.
func (h *Handler) GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stories, err := h.storiesFunc()
	if err != nil {
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		h.logFunc(errors.Wrap(err, "get stories"))
		return
	}

	view := indexView{
		t:       h.t,
		Stories: stories,
		Meta: metaView{
			t:           h.t,
			Title:       indexMetaTitle,
			Description: indexMetaDescription,
			Type:        indexMetaType,
		},
	}

	if err := h.t.ExecuteTemplate(w, "layout.html", view); err != nil {
		h.logFunc(errors.Wrap(err, "execute index template"))
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
	}
}

// GetCV /.
func (h *Handler) GetCV(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/pdf")
	if _, err := w.Write(h.cv); err != nil {
		h.logFunc(errors.Wrap(err, "write cv"))
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
	}
}

type indexView struct {
	t       *template.Template
	Stories []medium.Story
	Meta    metaView
}

// Content returns html index template as html.
func (i indexView) Content() (template.HTML, error) {
	var c template.HTML

	buf := new(bytes.Buffer)
	if err := i.t.ExecuteTemplate(buf, "index.html", i); err != nil {
		return c, errors.Wrap(err, "execute index template")
	}
	c = template.HTML(buf.String())

	return c, nil
}

type metaView struct {
	t                        *template.Template
	Title, Description, Type string
	HideOG                   bool
}

// Content returns html meta tags using template.HTML.
func (m metaView) Content() (template.HTML, error) {
	var c template.HTML

	buf := new(bytes.Buffer)
	if err := m.t.ExecuteTemplate(buf, "meta.html", m); err != nil {
		return c, errors.Wrap(err, "execute meta template")
	}
	if err := m.t.ExecuteTemplate(buf, "og.html", m); err != nil {
		return c, errors.Wrap(err, "execute og template")
	}
	c = template.HTML(buf.String())

	return c, nil
}
