package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lauz1n/go-gateway/internal/domain"
	"github.com/lauz1n/go-gateway/internal/dto"
	"github.com/lauz1n/go-gateway/internal/service"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateInvoiceInput
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.APIKey = r.Header.Get("X-API-KEY")

	output, err := h.invoiceService.Create(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *InvoiceHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(w, "X-API-KEY is required", http.StatusBadRequest)
		return
	}

	output, err := h.invoiceService.GetById(id, apiKey)

	if err != nil {
		switch err {
		case domain.ErrInvoiceNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case domain.ErrUnauthorizedAccess:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *InvoiceHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(w, "X-API-KEY is required", http.StatusUnauthorized)
		return
	}

	accountID := chi.URLParam(r, "account_id")
	if accountID == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	output, err := h.invoiceService.ListByAccount(accountID)
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
