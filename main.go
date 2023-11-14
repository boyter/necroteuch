package main

import (
	"context"
	"embed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"necroteuch/common"
	"necroteuch/data"
	"necroteuch/handlers"
	"necroteuch/service"
	"net/http"
	"strconv"
	"time"
)

//go:embed assets/static
var staticFiles embed.FS

//go:embed schema.sql
var schema string

func main() {
	config := common.NewConfig()
	log.Info().Str(common.UniqueCode, "22840cb2").Interface("config", config).Msg("config")

	db, err := data.ConnectDB(config)
	if err != nil {
		log.Error().Str(common.UniqueCode, "99d9ab69").Err(err).Msg("error creating database")
		return
	}

	q := data.New(db)
	err = q.CreateSchema(context.TODO(), schema)
	if err != nil {
		log.Error().Str(common.UniqueCode, "a26bb9c7").Err(err).Msg("error creating schema")
		return
	}

	ser := service.NewService(config)

	switch config.LogLevel {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	app, err := handlers.NewApplication(ser, staticFiles)
	if err != nil {
		log.Error().Str(common.UniqueCode, "7fcf2c4c").Err(err).Msg("error creating application")
		return
	}

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(config.HTTPPort),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           app.Routes(),
	}

	log.Log().Str(common.UniqueCode, "ac39dae8").Msg("starting server on :" + strconv.Itoa(config.HTTPPort))
	err = srv.ListenAndServe()
	log.Error().Str(common.UniqueCode, "2c6e5b79").Err(err).Msg("exiting server")
}
