package render

import (
	"net/http"
	"testing"

	"github.com/Lokiop/Bookings/internal/models"
)

func TestAddTemlateData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddTemplateData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found n session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Rendertemplate(&ww, r, "home.page.html", &models.TemplateData{})

	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Rendertemplate(&ww, r, "non-existent.page.html", &models.TemplateData{})

	if err == nil {
		t.Error("Rendered a template that does not exit")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewTemplate(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}
}
