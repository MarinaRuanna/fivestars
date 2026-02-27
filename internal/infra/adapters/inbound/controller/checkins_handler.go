package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
	"fivestars/internal/infra/auth"
)

type CheckinsHandler struct {
	createUC usecases.CreateCheckinUseCase
	listUC   usecases.ListCheckinsUseCase
}

func NewCheckinsHandler(createUC usecases.CreateCheckinUseCase, listUC usecases.ListCheckinsUseCase) *CheckinsHandler {
	return &CheckinsHandler{createUC: createUC, listUC: listUC}
}

func (h *CheckinsHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not authenticated"})
		return
	}

	var dto createCheckinDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
		return
	}
	if dto.Lat == nil || dto.Lng == nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "lat/lng required"})
		return
	}

	checkin, err := ToDomainCheckin(&dto, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res, err := h.createUC.Execute(r.Context(), *checkin)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "establishment not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if err.Error() == "check-in already performed today for this establishment" {
			w.WriteHeader(http.StatusConflict)
		} else if err.Error() == "user too far from establishment" || err.Error() == "establishment has no location" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	checkinDTO, err := ToCheckinDTO(res)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to serialize response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(checkinDTO)
}

func (h *CheckinsHandler) ListMyCheckins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not authenticated"})
		return
	}

	list, err := h.listUC.Execute(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}
