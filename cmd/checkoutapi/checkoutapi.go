package main

import (
	"log"
	"net/http"

	pricingSvc "github.com/TheFodfather/checkoutapi/pricing/service"

	"github.com/TheFodfather/checkoutapi/checkout/handler"
	"github.com/TheFodfather/checkoutapi/checkout/repository"
)

func main() {
	pricer, err := pricingSvc.New("./cmd/configs/pricing.json")
	if err != nil {
		log.Fatalf("❌ Could not start pricing service - err=%q", err)
	}

	repo := repository.NewInMemoryRepository()
	httpHandler := handler.New(repo, pricer)

	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	port := "8080"
	log.Printf("🚀 Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("❌ Could not start server - err=%q", err)
	}
}
