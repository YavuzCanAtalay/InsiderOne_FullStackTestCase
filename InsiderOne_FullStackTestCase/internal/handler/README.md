# internal/handler

## Purpose
HTTP layer. Maps incoming requests to service calls and formats JSON responses.

## What goes here
- `team_handler.go` — handles `GET /teams`
- `match_handler.go` — handles `GET /matches`, `GET /matches/week/{week}`, `PUT /matches/{id}`
- `standings_handler.go` — handles `GET /standings`
- `week_handler.go` — handles `POST /weeks/next`
- `league_handler.go` — handles `POST /league/play-all`
- `prediction_handler.go` — handles `GET /predictions/current`

## What needs to be done
- Parse path parameters and request bodies
- Call the appropriate service method
- Return consistent JSON with proper HTTP status codes (200, 201, 400, 404, 500)
- No business logic here — handlers only translate HTTP ↔ service
