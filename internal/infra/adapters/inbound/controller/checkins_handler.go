package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
)

const checkinBodyMaxBytes int64 = 16 << 10 // 16KB

type CheckinsHandler struct {
	createUC usecases.CreateCheckinUseCase
	listUC   usecases.ListCheckinsUseCase
}

func NewCheckinsHandler(createUC usecases.CreateCheckinUseCase, listUC usecases.ListCheckinsUseCase) *CheckinsHandler {
	return &CheckinsHandler{createUC: createUC, listUC: listUC}
}

func (h *CheckinsHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	var dto createCheckinDTO
	if err := decodeStrictJSONBody(w, r, &dto, checkinBodyMaxBytes); err != nil {
		return err
	}
	if dto.Lat == nil || dto.Lng == nil {
		return customerror.NewValidationError("lat/lng required")
	}

	checkin, err := ToDomainCheckin(&dto, userID)
	if err != nil {
		return err
	}

	res, err := h.createUC.Execute(r.Context(), *checkin)
	if err != nil {
		return err
	}

	checkinDTO, err := ToCheckinDTO(res)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(checkinDTO); err != nil {
		return err
	}

	return nil
}

func (h *CheckinsHandler) ListMyCheckins(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	list, err := h.listUC.Execute(r.Context(), userID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		return err
	}

	return nil
}
