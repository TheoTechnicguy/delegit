package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"git.licolas.net/delegit/delegit/database"
	"git.licolas.net/delegit/delegit/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	host string = "0.0.0.0"
	port uint   = 41990
)

var (
	globalLogger zerolog.Logger = initLogger()
	logger       zerolog.Logger = GetLogger("main")
)

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	consoleLogger := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(consoleLogger).With().Timestamp().Logger()
}

func GetLogger(module string) zerolog.Logger {
	return globalLogger.With().Str("module", module).Logger()
}

func main() {
	logger.Info().Str("host", host).Uint("port", port).Msg("starting server")
	db, err := database.NewDatabase("sqlite", "feedback.db")
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to get database")
	}

	r := gin.Default()
	routes.RegisterFeedbackEndpoints(db, r)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r)

	fmt.Printf("%e\n", err)
}
