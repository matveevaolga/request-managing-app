package handler

import (
	"net/http"

	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
)

type ProjectTypeHandler struct {
	service *service.ProjectTypeService
}

func NewProjectTypeHandler(svc *service.ProjectTypeService) *ProjectTypeHandler {
	return &ProjectTypeHandler{service: svc}
}

func (h *ProjectTypeHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetAllProjects(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get project types", err)
		return
	}

	resp := make([]dto.ProjectTypeResponse, len(types))
	for i, pt := range types {
		resp[i] = dto.ProjectTypeResponse{
			ID:   pt.ID,
			Name: pt.Name,
		}
	}

	RespondWithJSON(w, http.StatusOK, resp)
}
