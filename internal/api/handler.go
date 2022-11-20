package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"xm/internal/model"
)

type (
	service interface {
		New(ctx context.Context, company *model.Company) (uuid.UUID, error)
		Delete(ctx context.Context, companyID uuid.UUID) error
		Update(ctx context.Context, companyID uuid.UUID, company *model.Company) error
		Get(ctx context.Context, companyID uuid.UUID) (model.Company, error)
	}

	newResponse struct {
		ID uuid.UUID
	}

	handler struct {
		logger  *zap.Logger
		service service
	}
)

func NewHandler(
	logger *zap.Logger,
	service service,
) *handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) RegisterRoutes(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/").HandlerFunc(h.new)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.get)
	r.Methods(http.MethodPatch).Path("/{id}").HandlerFunc(h.patch)
	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(h.delete)
}

func (h *handler) get(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	company, err := h.service.Get(request.Context(), uuid)
	if err != nil {
		h.handleError(err, writer)
		return
	}

	if err := json.NewEncoder(writer).Encode(company); err != nil {
		h.handleError(err, writer)
	}
}

func (h *handler) new(writer http.ResponseWriter, request *http.Request) {
	var company model.Company
	err := json.NewDecoder(request.Body).Decode(&company)
	if err != nil {
		h.logger.Error("decode error", zap.Error(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := company.Validate(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	
	uuid, err := h.service.New(request.Context(), &company)
	if err != nil {
		h.handleError(err, writer)
		return
	}

	if err := json.NewEncoder(writer).Encode(newResponse{ID: uuid}); err != nil {
		h.handleError(err, writer)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (h *handler) patch(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	var company model.Company
	if err := json.NewDecoder(request.Body).Decode(&company); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := company.Validate(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Update(request.Context(), uuid, &company); err != nil {
		h.handleError(err, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *handler) delete(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err := h.service.Delete(request.Context(), uuid); err != nil {
		h.handleError(err, writer)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h *handler) handleError(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	if errors.Is(err, model.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Error("failed to process request", zap.Error(err))
	http.Error(w, "something went wrong", http.StatusInternalServerError)
}
