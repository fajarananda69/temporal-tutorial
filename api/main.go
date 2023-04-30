package main

import (
	"log"
	"net/http"

	"temporal-tutorial/api/handler"

	"github.com/go-chi/chi"
	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	handler.Temporal = c

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetProductsHandler)
	})

	router.Route("/cart", func(r chi.Router) {
		r.Post("/", handler.CreateCartHandler)
		r.Get("/{workflowID}", handler.GetCartHandler)
		r.Post("/{workflowID}/add", handler.AddToCartHandler)
		r.Post("/{workflowID}/remove", handler.RemoveFromCartHandler)
		r.Post("/{workflowID}/checkout", handler.CheckoutHandler)
		r.Post("/{workflowID}/payment", handler.PaymentHandler)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Println(method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicln(err)
	}

	server := http.Server{
		Addr:    ":1234",
		Handler: router,
	}

	log.Println("Temporal Tutorial API serving at", server.Addr)
	log.Fatal(server.ListenAndServe())
}
