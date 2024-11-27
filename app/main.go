package main

import (
	"app/controllers"
	"app/utils"
	"app/utils/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if len(os.Args) > 1 {
		handleCLI()
		return
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	// Health Check
	router.Get("/", controllers.RootHandler)

	// --- USERS ---
	// Register an User
	router.Post("/auth/register", controllers.RegisterHandler)
	// Login as an User
	router.Post("/auth/login", controllers.LoginHandler)
	// Routes that need authentication
	router.Group(func(rtr chi.Router) {
		rtr.Use(middlewares.AuthMiddleware)
		// Get an User
		rtr.Get("/users/{user_id}", controllers.GetUserByIdHandler)
		// Get all Users (+ Query)
		rtr.Get("/users", controllers.GetAllUsersHandler)
		// Patch an User
		rtr.Patch("/users/{user_id}", controllers.UpdateUserHandler)
		// Delete an User
		rtr.Delete("/users/{user_id}", controllers.DeleteUserHandler)
	})

	// --- PETS ---
	// Routes that need authentication
	router.Group(func(pRtr chi.Router) {
		pRtr.Use(middlewares.AuthMiddleware)
		// Create a Pet
		pRtr.Post("/pets", controllers.CreatePetHandler)
		// Get a Pet
		pRtr.Get("/pets/{pet_id}", controllers.GetPetByIdHandler)
		// Get all Pets (+ Query)
		pRtr.Get("/pets", controllers.GetAllPetsHandler)
		// Patch a Pet
		pRtr.Patch("/pets/{pet_id}", controllers.UpdatePetHandler)
		// Delete a Pet
		pRtr.Delete("/pets/{pet_id}", controllers.DeletePetHandler)
	})

	// --- CLIENTS ---
	// Routes that need authentication
	router.Group(func(cRtr chi.Router) {
		cRtr.Use(middlewares.AuthMiddleware)
		// Create a Client
		cRtr.Post("/clients", controllers.CreateClientHandler)
		// Get a Client
		cRtr.Get("/clients/{client_id}", controllers.GetClientByIdHandler)
		// Get all Clients (+ Query)
		cRtr.Get("/clients", controllers.GetAllClientsHandler)
		// Patch a Client
		cRtr.Patch("/clients/{client_id}", controllers.UpdateClientHandler)
		// Delete a Client
		cRtr.Delete("/clients/{client_id}", controllers.DeleteClientHandler)
	})

	fmt.Println("[server] API Server Running on Port", utils.SERVER_PORT)
	serverAddr := fmt.Sprintf("0.0.0.0:%d", utils.SERVER_PORT)
	http.ListenAndServe(serverAddr, router)

	fmt.Println("[server] API Server Exiting")
}
