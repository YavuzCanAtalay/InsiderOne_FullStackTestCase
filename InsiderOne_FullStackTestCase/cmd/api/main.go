package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/handler"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/prediction"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/simulator"
)

func main() {
	db := connectDB() // Connection to database
	defer db.Close()  // close database when main is terminated

	// Repositories
	teamRepo := repository.NewTeamRepository(db)   // Object creation for team repository
	matchRepo := repository.NewMatchRepository(db) // Object creation for match repository

	// Simulator
	sim := simulator.NewMatchSimulator() // Object creation for match simulator

	// Services
	leagueSvc := service.NewLeagueService(teamRepo, matchRepo, sim) // Object creation for league service

	// Prediction engine
	predEngine := prediction.NewPredictionEngine(teamRepo, matchRepo, sim) // Object creation for prediction engine

	// Handlers, create HTTP-facing object that will receive requests
	teamHandler := handler.NewTeamHandler(teamRepo)                          // Object creation for team handler
	matchHandler := handler.NewMatchHandler(matchRepo)                       // Object creation for match handler
	standingsHandler := handler.NewStandingsHandler(leagueSvc)               // Object creation for standings handler
	weekHandler := handler.NewWeekHandler(leagueSvc, predEngine)              // Object creation for week handler
	leagueHandler := handler.NewLeagueHandler(leagueSvc)                     // Object creation for league handler
	predictionHandler := handler.NewPredictionHandler(predEngine, matchRepo) // Object creation for prediction handler

	// Routes , wire dependencies together and start server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("/teams", teamHandler.GetTeams)
	mux.HandleFunc("/matches", matchHandler.GetMatches)
	mux.HandleFunc("/matches/week/", matchHandler.GetMatchesByWeek)
	mux.HandleFunc("/matches/", matchHandler.UpdateMatch)
	mux.HandleFunc("/standings", standingsHandler.GetStandings)
	mux.HandleFunc("/weeks/next", weekHandler.PlayNextWeek)
	mux.HandleFunc("/league/play-all", leagueHandler.PlayAll)
	mux.HandleFunc("/predictions/current", predictionHandler.GetPredictions)

	port := os.Getenv("API_PORT") // read environment variable for port
	if port == "" {
		port = "8080"
	}

	// Start HTTP server, waiting for requests
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}
