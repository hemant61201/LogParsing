package main

import (
	"LogParsing/internal/config"
	"LogParsing/internal/loading"
	"LogParsing/internal/parsing"
	"LogParsing/internal/storage/sqlite"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// load config
	slog.Info("Loading config...")

	config := config.MustLoad()

	slog.Info("Config loaded successfully.")

	// setup database

	slog.Info("Connecting to database...")

	storage, err := sqlite.NewSqlite(config)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database initialized", slog.String("env", config.Env), slog.String("version", "1.0.0"))

	//start reading log file

	logs := loading.LoadLogs{FileName: config.LogFilePath}

	slog.Info("Loading logs...", slog.String("logFilePath", config.LogFilePath))

	if err := logs.Load(); err != nil {
		slog.Error("Load failed", slog.String("err", err.Error()))
		return
	}

	slog.Info("Parsing logs...")

	parsing.Parse(storage)

	// setup router

	router := gin.New()

	// Get logs of particular type

	router.GET("/logs/:logType", func(context *gin.Context) {

		slog.Info("Getting logs from database",
			slog.String("logType", context.Param("logType")),
		)

		result, err := storage.GetLog(context.Param("logType"))

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})

	// Get All logs

	router.GET("/logs", func(context *gin.Context) {

		slog.Info("Getting all logs")

		result, err := storage.GetAllLog()

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})

	server := http.Server{
		Addr:    config.Addr,
		Handler: router.Handler(),
	}

	slog.Info("server started", slog.String("address", config.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
