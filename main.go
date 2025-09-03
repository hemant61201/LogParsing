package main

import (
	"LogParsing/internal/config"
	"LogParsing/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"

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

	//logs := loading.LoadLogs{FileName: config.LogFilePath}
	//
	//slog.Info("Loading logs...", slog.String("logFilePath", config.LogFilePath))
	//
	//if err := logs.Load(); err != nil {
	//	slog.Error("Load failed", slog.String("err", err.Error()))
	//	return
	//}
	//
	//slog.Info("Parsing logs...")
	//
	//parsing.Parse(storage)

	// setup router

	router := gin.New()

	router.GET("/logs/{logType}", func(context *gin.Context) {

		//slog.Info("Getting logs from database",
		//	slog.String("logType", context.Param("logType")),
		//)

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

	router.Run(config.Addr)
}
