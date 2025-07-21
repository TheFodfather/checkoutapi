package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheFodfather/checkoutapi/checkout"
	"github.com/TheFodfather/checkoutapi/checkout/repository"
	"github.com/TheFodfather/checkoutapi/domain"
)

type PricingService interface {
	GetRules() map[string]domain.PricingRule
}

type HTTPHandler struct {
	repo   repository.SessionRepository
	pricer PricingService
}

func New(repo repository.SessionRepository, pricer PricingService) *HTTPHandler {
	return &HTTPHandler{repo: repo, pricer: pricer}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /checkouts", h.handleCreateCheckout)
	mux.HandleFunc("GET /checkouts/{checkoutID}", h.handleGetTotalPrice)
	mux.HandleFunc("POST /checkouts/{checkoutID}/scan", h.handleScanItem)
}

func (h *HTTPHandler) handleCreateCheckout(w http.ResponseWriter, r *http.Request) {
	session := checkout.New(h.pricer)
	if err := h.repo.Save(session); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not save session")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"checkoutId": session.GetID()})
}

func (h *HTTPHandler) handleScanItem(w http.ResponseWriter, r *http.Request) {
	checkoutID := r.PathValue("checkoutID")

	session, err := h.repo.Get(checkoutID)
	if err != nil {
		log.Printf("INFO: Session not found for checkoutID=%q err=%q", checkoutID, err)
		respondWithError(w, http.StatusNotFound, "session not found")
		return
	}

	var reqBody struct {
		SKU string `json:"sku"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("WARN: Failed to decode request body for checkoutID=%q err=%q", checkoutID, err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := session.Scan(reqBody.SKU); err != nil {
		log.Printf("WARN: Invalid SKU scan for checkoutID=%q: err=%q", checkoutID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Save(session); err != nil {
		log.Printf("ERROR: Failed to save session after scan for checkoutID %q: %v", checkoutID, err)
		respondWithError(w, http.StatusInternalServerError, "could not save session")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *HTTPHandler) handleGetTotalPrice(w http.ResponseWriter, r *http.Request) {
	checkoutID := r.PathValue("checkoutID")

	session, err := h.repo.Get(checkoutID)
	if err != nil {
		log.Printf("INFO: Session not found for checkoutID %q err=%q", checkoutID, err)
		respondWithError(w, http.StatusNotFound, "session not found")
		return
	}

	totalPrice, _ := session.GetTotalPrice()

	response := struct {
		CheckoutID string `json:"checkoutId"`
		TotalPrice int    `json:"totalPrice"`
	}{
		CheckoutID: session.GetID(),
		TotalPrice: totalPrice,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: Failed to marshal JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
