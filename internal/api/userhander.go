package api

import (
	"API_for_ping_servers/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

type ResponseApiKey struct {
	Status string `json:"status"`
	ApiKey string `json:"api_key"`
}

func GetApiKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := services.CreateApiKey()
	resp, err := json.Marshal(ResponseApiKey{
		Status: "created",
		ApiKey: key,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application-json")
	_, err = w.Write(resp)
	if err != nil {
		log.Fatal(err)
	}
}
