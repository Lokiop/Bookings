package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lokiop/Bookings/internal/config"
	"github.com/Lokiop/Bookings/internal/models"
	"github.com/Lokiop/Bookings/internal/render"
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

	render.Rendertemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	RemoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = RemoteIp

	render.Rendertemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays a form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "make-reservation.page.html", &models.TemplateData{})
}

// Generals renders the room page for generals quarters rooms
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "generals.page.html", &models.TemplateData{})
}

// Colonels renders the room page for the Colonels suite
func (m *Repository) Colonels(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "colonels.page.html", &models.TemplateData{})
}

// Availability renders the seacrh-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

// Availability renders the seacrh-availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is : %s & End date is : %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles the request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      false,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the Contact availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "contact.page.html", &models.TemplateData{})
}
