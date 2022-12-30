package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lokiop/Bookings/internal/config"
	"github.com/Lokiop/Bookings/internal/forms"
	"github.com/Lokiop/Bookings/internal/helpers"
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
	render.Rendertemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "about.page.html", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays a form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Rendertemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PoatReservation handles the posting of the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 4)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Rendertemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

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
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the Contact availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Rendertemplate(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Cannot get error from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Rendertemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
