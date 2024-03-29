package config

import (
	"html/template"
	"log"

	"github.com/Lokiop/Bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	Usecache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
