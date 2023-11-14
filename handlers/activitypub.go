package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"necroteuch/common"
	"net/http"
	"strings"
)

func (app *Application) WebFinger(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimSpace(r.FormValue("resource"))
	log.Info().Str(common.UniqueCode, "1c778539").Str("resource", resource).Str("ip", GetIP(r)).Msg("WebFinger")

	if resource == "" {
		log.Error().Str(common.UniqueCode, "0e2fbbfc").Str("ip", GetIP(r)).Msg("WebFinger Invalid Resource")
		http.Error(w, "Invalid Resource", 400)
		return
	}

	if !strings.HasPrefix(resource, "acct:") {
		log.Error().Str(common.UniqueCode, "8f372bdd").Str("ip", GetIP(r)).Msg("WebFinger Invalid Resource")
		http.Error(w, "Invalid Resource", 400)
		return
	}

	resource = resource[5:]
	s := strings.Split(resource, "@")
	if len(s) != 2 {
		log.Error().Str(common.UniqueCode, "25b7930b").Str("ip", GetIP(r)).Msg("WebFinger Invalid Resource")
		http.Error(w, "Invalid Resource", 400)
		return
	}

	resp := common.WebFingerResponse{
		Subject: "acct:" + resource,
		Links: []common.WebFingerLinks{
			{
				Rel:  "self",
				Type: "application/activity+json",
				//Href: fmt.Sprintf("%su/%s", app.Environment.BaseUrl, s[0]),
				Href: fmt.Sprintf("%su/%s", "REALURLHERE", s[0]),
			},
		},
	}

	b, err := json.MarshalIndent(resp, "", jsonIndent)
	if err != nil {
		log.Error().Str(common.UniqueCode, "857f707b").Str("ip", GetIP(r)).Msg("Error in JSON marshal")
		http.Error(w, "Invalid Resource", 500)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	_, _ = fmt.Fprint(w, string(b))
}
