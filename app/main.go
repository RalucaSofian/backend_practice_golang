package main

import (
	"app/controllers"
	"app/utils/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TODO: Env variable
const serverPort = ":3000"

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

	fmt.Println("[server] API Server Running on Port", serverPort)
	http.ListenAndServe(serverPort, router)

	fmt.Println("[server] API Server Exiting")
}
