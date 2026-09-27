package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
)

func ChechHealth(w http.ResponseWriter, r *http.Request) {

	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: http.StatusOK, Message: "healthy", Version: "v1"})
}
