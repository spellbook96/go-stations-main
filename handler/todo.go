package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	fmt.Printf("Running Create\n")
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode the request body: %v", err)
			http.Error(w, "Failed to decode the request body", http.StatusBadRequest)
			return
		}

		resp, err := h.Create(r.Context(), &req)
		if err != nil {
			log.Printf("Failed to create the TODO: %v", err)
			http.Error(w, "Failed to create the TODO", http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Failed to encode the response: %v", err)
			http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
			return
		}
		// fmt.Printf("%s\n", req.TODO.Subject)
	case http.MethodGet:
		var req model.ReadTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		resp, err := h.Read(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
	case http.MethodPut:
		var req model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		resp, err := h.Update(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
	case http.MethodDelete:
		var req model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		resp, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
