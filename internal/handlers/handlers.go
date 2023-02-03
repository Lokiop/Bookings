package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Lokiop/Bookings/internal/config"
	"github.com/Lokiop/Bookings/internal/drivers"
	"github.com/Lokiop/Bookings/internal/forms"
	"github.com/Lokiop/Bookings/internal/models"
	"github.com/Lokiop/Bookings/internal/render"
	"github.com/Lokiop/Bookings/internal/repository"
	"github.com/Lokiop/Bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

var Repo *Repository

// Reposirory type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new Repository
func NewRepo(a *config.AppConfig, db *drivers.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new Repo for testing
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// New Handlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays a form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Issue"
		stringMap["info"] = "Cannot get Reservation details from the Session"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Cannot find Room"
		stringMap["info"] = "Cannot get Rooms by Id"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("02-01-2006")
	ed := res.EndDate.Format("02-01-2006")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PoatReservation handles the posting of the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Error"
		stringMap["info"] = "Cannot get reservation form the Session"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()

	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the form. Please refresh the site."

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Error"
		stringMap["info"] = "Unable to Insert Reservation in the Reservation's table"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Error"
		stringMap["info"] = "Unable to insert Restrictions in the Restricion's table"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page for generals quarters rooms
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

// Colonels renders the room page for the Colonels suite
func (m *Repository) Colonels(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "colonels.page.html", &models.TemplateData{})
}

// Availability renders the seacrh-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// Availability renders the seacrh-availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "02-01-2006"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the Start Date. Please Refresh."

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the End Date. Please Refresh."

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Error"
		stringMap["info"] = "Unable to retrieve information. Please refresh."

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		//no availability
		var stringMap = make(map[string]string)
		stringMap["heading"] = "No Rooms Available"
		stringMap["info"] = "No Rooms available fot the given range of dates."

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles the request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal Server Error",
		}

		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "02-01-2006"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Error Connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the Contact availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Make a reservation"
		stringMap["info"] = "Cannot get to this page without making a reservation. Please go to reservation page"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("02-01-2006")
	ed := reservation.EndDate.Format("02-01-2006")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays the list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the Room Id"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Issue"
		stringMap["info"] = "Unable to get reservation from the session"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes url parameters, bulids a sessional variable and redirects user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the Room Id"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Internal Server Issue"
		stringMap["info"] = "Unable to get room by id"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the Start Date"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		var stringMap = make(map[string]string)
		stringMap["heading"] = "Parsing Issue"
		stringMap["info"] = "Unable to Parse the End Date"

		errorInfo := models.ErrorPage{
			ErrorInfo: stringMap,
		}

		m.App.Session.Put(r.Context(), "error", errorInfo)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ErrorRedirect(w http.ResponseWriter, r *http.Request) {
	stringMap, ok := m.App.Session.Get(r.Context(), "error").(models.ErrorPage)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	m.App.Session.Remove(r.Context(), "error")
	render.Template(w, r, "error.page.html", &models.TemplateData{
		StringMap: stringMap.ErrorInfo,
	})
}
