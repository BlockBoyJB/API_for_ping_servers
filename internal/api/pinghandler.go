package api

import (
	"API_for_ping_servers/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	AddPing struct {
		Url    string `json:"url"`
		Email  string `json:"email"`
		ApiKey string `json:"api_key"`
	}
	DeletePing struct {
		ApiKey string `json:"api_key"`
		Url    string `json:"url"`
	}
	AllPingsResponse struct {
		Urls []string `json:"urls"`
	}
)

type PingHandler struct {
	PingService services.PingService
}

func (h *PingHandler) AddPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var data AddPing
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(req, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !services.VerifyApiKey(data.ApiKey) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err = h.PingService.CreateNewJob(data.ApiKey, data.Email, data.Url); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *PingHandler) DeletePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var data DeletePing
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(req, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !services.VerifyApiKey(data.ApiKey) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err = h.PingService.DeleteJob(data.ApiKey, data.Url); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PingHandler) GetAllPigs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	apiKey := r.URL.Query().Get("api_key")
	if len(apiKey) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !services.VerifyApiKey(apiKey) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	data, err := h.PingService.GetAllJobs(apiKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(AllPingsResponse{Urls: data})
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
