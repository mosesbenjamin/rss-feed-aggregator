package handlers

import (
	"net/http"

	json_resp "github.com/mosesbenjamin/rss-feed-aggregator/helpers"
)

func Readiness(w http.ResponseWriter, r *http.Request) {
	json_resp.RespondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func Error(w http.ResponseWriter, r *http.Request) {
	json_resp.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
