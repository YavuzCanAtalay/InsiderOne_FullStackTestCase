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
	db := connectDB()
	defer db.Close()

	// Repositories
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	// Simulator
	sim := simulator.NewMatchSimulator()

	// Services
	leagueSvc := service.NewLeagueService(teamRepo, matchRepo, sim)

	// Prediction engine
	predEngine := prediction.NewPredictionEngine(teamRepo, matchRepo, sim)

	// Handlers
	teamHandler := handler.NewTeamHandler(teamRepo)
	matchHandler := handler.NewMatchHandler(matchRepo)
	standingsHandler := handler.NewStandingsHandler(leagueSvc)
	weekHandler := handler.NewWeekHandler(leagueSvc)
	leagueHandler := handler.NewLeagueHandler(leagueSvc)
	predictionHandler := handler.NewPredictionHandler(predEngine, matchRepo)

	// Routes
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

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func connectDB() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

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
