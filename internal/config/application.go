package config

import (
	"githun.com/farmani/snippetbox/internal/models"
	"log"
)

// Application Define an application struct to hold the application-wide dependencies for the
// web application. For now, we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Snippets *models.SnippetModel
}
