package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sakajunquality/clouddeploy-functions/slackbot"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"

	_ "github.com/GoogleCloudPlatform/berglas/pkg/auto"
)

var (
	logger zerolog.Logger
	client *slackbot.Slackbot
)

func init() {
	// @todo refactor
	client = slackbot.NewSlackbot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL"))
	client.SetStateBucket(os.Getenv("SLACK_BOT_STATE_BUCKET"))
}

func main() {
	logger = httplog.NewLogger("clouddeploy-notification", httplog.Options{
		JSON: true,
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Timeout(360 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Msg("healthy")
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/pubsub/push/operations", HandleOperationsTopic)
	r.Post("/pubsub/push/approvals", HandleApprovalsTopic)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}
