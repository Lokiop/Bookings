package handlers

import (
	"net/http"

	"github.com/Lokiop/Bookings/package/config"
	"github.com/Lokiop/Bookings/package/models"
	"github.com/Lokiop/Bookings/package/render"
)

var Repo *Repository

// Reposirory type
type Repository struct {
	App *config.AppConfig
}

// Creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// New Handlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Rendertemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	RemoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = RemoteIp

	render.Rendertemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
