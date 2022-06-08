package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

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
	fmt.Printf("Running Read, PrevID:%d,Size:%d\n", req.PrevID, req.Size)
	TODOS, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	fmt.Printf("%+v\n", TODOS)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: TODOS}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	return &model.DeleteTODOResponse{}, err
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
		values := r.URL.Query()
		fmt.Printf("values:%+v\n", values)
		prevID, err := strconv.ParseInt(values.Get("prev_id"), 10, 64)
		if err != nil {
			prevID = 0
		}
		size, err := strconv.ParseInt(values.Get("size"), 10, 64)
		if err != nil {
			size = 0
		}
		req = model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
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
			w.WriteHeader(http.StatusBadRequest)
			log.Println("[ERROR]", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		var req model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}
		if reflect.DeepEqual([]int64{}, req.IDs) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
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
