package main

import (
	"net/http"
	"os"

	"github.com/adamfdl/tenx/database"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// =========================================================================
	// Zap Logger Support
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := loggerConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	defer logger.Sync()

	// =========================================================================
	// Start Database
	db, err := database.Open()
	if err != nil {
		logger.Fatal(err.Error(), zap.String("op", "database.open"))
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		logger.Fatal(err.Error(), zap.String("op", "database.ping"))
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing database connection gracefully", zap.String("op", "database.close"))
		db.Close()
	}()

	// =========================================================================
	// Initialize Handlers
	server := http.Server{
		Addr: "localhost:3000",
		Handler: func() http.Handler {
			r := mux.NewRouter()
			r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			})

			ctrl := newController(logger, db)
			r.HandleFunc("/auth/jira", ctrl.hAuthJira).Methods("GET")
			r.HandleFunc("/auth/jira/callback", ctrl.hAuthJiraCallback).Methods("GET")
			r.HandleFunc("/auth/harvest", ctrl.hAuthHarvest).Methods("GET")
			r.HandleFunc("/auth/harvest/callback", ctrl.hAuthHarvestCallback).Methods("GET")

			r.HandleFunc("/user", ctrl.hUserSignUp).Methods("POST")

			//r.HandleFunc("/sprint/current", ctrl.hSprintCurrent)
			return r
		}(),
	}

	logger.Info("server started", zap.Int("port", 3000))
	server.ListenAndServe()
}
