package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"git.licolas.net/delegit/delegit/database"
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

	router := gin.Default()
	router.GET("/feedback", func(ctx *gin.Context) {
		feedback, err := db.GetAllFeedback()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		}
		ctx.JSON(http.StatusOK, feedback)
	})

	router.POST("/feedback", func(ctx *gin.Context) {
		var feedback database.Feedback
		if err := ctx.ShouldBind(&feedback); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		f, err := db.AddFeedback(&feedback)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, f)
	})

	router.GET("/feedback/:id", func(ctx *gin.Context) {
		_id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		id := uint(_id)
		logger.Debug().Uint("feedback", id).Msg("looking for feedback by ID")

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		feedback, err := db.GetFeedback(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		ctx.JSON(http.StatusOK, feedback)
	})

	router.PUT("/feedback", func(ctx *gin.Context) {
		var feedback *database.Feedback
		if err := ctx.ShouldBind(feedback); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		feedback, err := db.UpdateFeedback(feedback)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		ctx.JSON(http.StatusOK, feedback)
	})

	router.DELETE("/feedback", func(ctx *gin.Context) {
		var feedback *database.Feedback
		if err := ctx.ShouldBind(feedback); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := db.DeleteFeedback(feedback); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		feedback.ID = 0
		ctx.JSON(http.StatusOK, feedback)
	})

	router.OPTIONS("/feedback", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Autorization")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "300")
		ctx.AbortWithStatus(http.StatusNoContent)
	})

	err = http.ListenAndServe(
		fmt.Sprintf("%s:%d", host, port),
		router,
	)

	fmt.Printf("%e\n", err)
}
