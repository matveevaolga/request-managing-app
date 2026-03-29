package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/matveevaolga/request-managing-app/internal/transport/middleware"
)

type ApplicationHandler struct {
	service  *service.ApplicationService
	validate *validator.Validate
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		service:  svc,
		validate: validator.New(),
	}
}

func (h *ApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
		return
	}

	id, err := h.service.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProjectTypeNotFound):
			RespondWithError(w, http.StatusBadRequest, "Invalid project type", err)
		case errors.Is(err, domain.ErrApplicationAlreadyExists):
			RespondWithError(w, http.StatusBadRequest, "Application with this project name and email already exists", err)
		case strings.Contains(err.Error(), "invalid phone format"):
			RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		default:
			RespondWithError(w, http.StatusInternalServerError, "Failed to create application", err)
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, id)
}

func (h *ApplicationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.getIDFromPath(r.URL.Path, "")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	app, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrApplicationNotFound:
			RespondWithError(w, http.StatusNotFound, "Application not found", err)
		default:
			RespondWithError(w, http.StatusInternalServerError, "Failed to get application", err)
		}
		return
	}

	pt, err := h.service.GetProjectTypeByID(r.Context(), app.TypeID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get project type", err)
		return
	}

	resp := dto.ApplicationResponse{
		ApplicationID:         app.ID,
		FullName:              app.FullName,
		Email:                 app.Email,
		Phone:                 app.Phone,
		OrganisationName:      app.OrganisationName,
		OrganisationURL:       app.OrganisationURL,
		ProjectName:           app.ProjectName,
		TypeName:              pt.Name,
		ExpectedResults:       app.ExpectedResults,
		IsPayed:               app.IsPayed,
		AdditionalInformation: app.AdditionalInformation,
		Status:                string(app.Status),
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *ApplicationHandler) GetAllFiltered(w http.ResponseWriter, r *http.Request) {
	params := h.createFilterParams(r)

	apps, total, err := h.service.GetAllFiltered(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get applications", err)
		return
	}

	resp := dto.ApplicationListResponse{
		Count:        total,
		Applications: make([]dto.ApplicationPreviewResponse, len(apps)),
	}

	for i, app := range apps {
		resp.Applications[i] = dto.ApplicationPreviewResponse{
			ExternalApplicationID: app.ID,
			ProjectName:           app.ProjectName,
			TypeName:              app.TypeName,
			Initiator:             app.Initiator,
			OrganisationName:      app.OrganisationName,
			DateUpdated:           app.DateUpdated,
			Status:                string(app.Status),
			RejectionMessage:      app.RejectionMessage,
		}
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *ApplicationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	id, err := h.getIDFromPath(r.URL.Path, "/accept")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	reviewerID := r.Context().Value(middleware.UserIDKey).(int64)

	err = h.service.Accept(r.Context(), id, reviewerID)
	if err != nil {
		switch err {
		case domain.ErrApplicationNotFound:
			RespondWithError(w, http.StatusNotFound, "Application not found", err)
		case domain.ErrApplicationNotPending:
			RespondWithError(w, http.StatusBadRequest, "Application is not in pending status", err)
		default:
			RespondWithError(w, http.StatusInternalServerError, "Failed to accept application", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ApplicationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := h.getIDFromPath(r.URL.Path, "/reject")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	var req dto.RejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
		return
	}

	reviewerID := r.Context().Value(middleware.UserIDKey).(int64)

	err = h.service.Reject(r.Context(), id, reviewerID, req.Reason)
	if err != nil {
		switch err {
		case domain.ErrApplicationNotFound:
			RespondWithError(w, http.StatusNotFound, "Application not found", err)
		case domain.ErrApplicationNotPending:
			RespondWithError(w, http.StatusBadRequest, "Application is not in pending status", err)
		default:
			RespondWithError(w, http.StatusInternalServerError, "Failed to reject application", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ApplicationHandler) getIDFromPath(path, suffix string) (int64, error) {
	idStr := strings.TrimPrefix(path, "/project/application/external/")
	idStr = strings.TrimSuffix(idStr, suffix)
	if idStr == "" {
		return 0, nil
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func (h *ApplicationHandler) createFilterParams(r *http.Request) repository.ApplicationFilterParameters {
	params := repository.ApplicationFilterParameters{
		Limit:  20,
		Offset: 0,
	}

	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		active := activeStr == "true"
		params.Active = &active
	}

	if search := r.URL.Query().Get("search"); search != "" {
		params.Search = &search
	}

	if typeIDStr := r.URL.Query().Get("projectTypeId"); typeIDStr != "" {
		if typeID, err := strconv.ParseInt(typeIDStr, 10, 64); err == nil {
			params.ProjectTypeID = &typeID
		}
	}

	if sortBy := r.URL.Query().Get("sortByDateUpdated"); sortBy != "" {
		params.SortByDateUpdated = sortBy
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	return params
}
