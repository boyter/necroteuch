package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"necroteuch/common"
	"net/http"
	"time"
)

const (
	jsonContentType = "application/json; charset=utf-8"
	jsonIndent      = "    "
)

type Timing struct {
	TimeMillis int64  `json:"timeMillis"`
	Source     string `json:"source"`
}

type HealthCheckResult struct {
	Success  bool                `json:"success"`
	Messages []string            `json:"messages"`
	Time     time.Time           `json:"time"`
	Timing   []Timing            `json:"timing"`
	Response HealthCheckResponse `json:"response"`
}

type HealthCheckResponse struct {
	IPAddress string `json:"ipAddress"`
	MemUsage  string `json:"memUsage"`
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Info().Str(common.UniqueCode, "b10037c0").Str("ip", GetIP(r)).Msg("HealthCheck")
	start := time.Now().UnixMilli()

	t, err := json.MarshalIndent(HealthCheckResult{
		Success:  true,
		Messages: []string{},
		Time:     time.Now().UTC(),
		Timing: []Timing{
			{
				Source:     "HealthCheck",
				TimeMillis: time.Now().UnixMilli() - start,
			},
		},
		Response: HealthCheckResponse{
			IPAddress: GetIP(r),
			MemUsage:  common.MemUsage(),
		},
	}, "", jsonIndent)
	if err != nil {
		log.Error().Str(common.UniqueCode, "5e192b2e").Str("ip", GetIP(r)).Msg("Error in JSON marshal")
		http.Error(w, "Invalid Resource", 500)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	_, _ = fmt.Fprint(w, string(t))
}
