package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mdjunior/ct/logger"
	"github.com/mdjunior/modulus/models"
	uuid "github.com/satori/go.uuid"
)

// GamesList is a method that list all Games
func (h DBHandler) GamesList(c echo.Context) error {
	// Capturando tempo para metricas
	initTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0

	games := h.cache.Items()
	var gamesResult []*models.Game
	for _, v := range games {
		game := v.Object.(*models.Game)

		if game.Status == "OPEN" {
			var gameClear models.Game
			gameClear.ID = game.ID
			gameClear.MaxUsers = game.MaxUsers
			gameClear.Name = game.Name
			gameClear.Status = game.Status
			gamesResult = append(gamesResult, &gameClear)
		} else {
			gamesResult = append(gamesResult, game)
		}
	}

	finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
	logger.Log(map[string]interface{}{
		"_action":   "GamesList.Find",
		"_result":   "success",
		"_duration": finishTime - initTime,
	})

	logger.Log(map[string]interface{}{
		"_action":       "GamesList",
		"_rid":          c.Get(echo.HeaderXRequestID),
		"_real-ip":      c.RealIP,
		"_duration":     finishTime - initTime,
		"_result":       "success",
		"short_message": "List Games",
	})

	if len(gamesResult) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	return c.JSON(http.StatusOK, gamesResult)
}

// GamesCreate is a method that create a game
func (h DBHandler) GamesCreate(c echo.Context) error {
	// Capturando tempo para metricas
	initTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0

	// Importando na struct
	game := new(models.Game)
	if err := c.Bind(game); err != nil {
		finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
		logger.Log(map[string]interface{}{
			"_action":       "GamesCreate.Bind",
			"_result":       "fail",
			"_duration":     finishTime - initTime,
			"short_message": err.Error,
		})
		return c.JSON(http.StatusBadRequest,
			map[string]string{"result": "fail", "details": err.Error()})
	}

	if game.ID == "" {
		game.ID = uuid.NewV4().String()
	}
	if game.Name == "" {
		game.Name = uuid.NewV4().String()
	}
	if game.Status == "" {
		game.Status = "OPEN"
	}
	if game.MaxUsers == 0 {
		game.MaxUsers = 2
	}
	game.Users = []map[string]int{}
	h.cache.Set(game.ID, game, 5*time.Minute)

	finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
	logger.Log(map[string]interface{}{
		"_action":       "GamesCreate",
		"_rid":          c.Get(echo.HeaderXRequestID),
		"_real-ip":      c.RealIP,
		"_duration":     finishTime - initTime,
		"_result":       "success",
		"short_message": "Created Games",
	})

	return c.JSON(http.StatusOK, game)
}

// GamesGET is a method that list a game by id
func (h DBHandler) GamesGET(c echo.Context) error {
	// Capturando tempo para metricas
	initTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0

	games, err := h.cache.Get(c.Param("id"))
	if err == false {
		return c.JSON(http.StatusNotFound, []string{})
	}

	finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
	logger.Log(map[string]interface{}{
		"_action":   "GamesGET.Find",
		"_result":   "success",
		"_duration": finishTime - initTime,
	})

	logger.Log(map[string]interface{}{
		"_action":       "GamesGET",
		"_rid":          c.Get(echo.HeaderXRequestID),
		"_real-ip":      c.RealIP,
		"_duration":     finishTime - initTime,
		"_result":       "success",
		"short_message": "Get Games",
	})

	return c.JSON(http.StatusOK, games.(*models.Game))
}

// GamesJoin is a method that insert a player on game
func (h DBHandler) GamesJoin(c echo.Context) error {
	// Capturando tempo para metricas
	initTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0

	games, err := h.cache.Get(c.Param("id"))
	if err == false {
		return c.JSON(http.StatusNotFound, []string{})
	}
	game := games.(*models.Game)

	// Importando na struct
	try := new(models.Try)
	if err := c.Bind(try); err != nil {
		finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
		logger.Log(map[string]interface{}{
			"_action":       "GamesJoin.Bind",
			"_result":       "fail",
			"_duration":     finishTime - initTime,
			"short_message": err.Error,
		})
		return c.JSON(http.StatusBadRequest,
			map[string]string{"result": "fail", "details": err.Error()})
	}

	// Verificando status do jogo
	if game.Status != "OPEN" {
		return c.JSON(http.StatusBadRequest,
			map[string]string{"result": "fail", "details": "Jogo não está aberto a novos participantes"})
	}

	// Verificando numero máximo de participantes
	if len(game.Users) >= game.MaxUsers {
		game.Status = "CLOSED"
		h.cache.Set(game.ID, game, 5*time.Minute)
		return c.JSON(http.StatusBadRequest,
			map[string]string{"result": "fail", "details": "Jogo não está aberto a novos participantes"})
	}

	// Adicionando participante
	game.Users = append(game.Users, map[string]int{try.Name: try.Value})

	// Verificando se é o ultimo
	if game.MaxUsers == len(game.Users) {
		game.Status = "CLOSED"

		// Obtendo o vencedor
		var sum int
		for i := range game.Users {
			for k := range game.Users[i] {
				sum += game.Users[i][k]
			}
		}
		winner := sum % game.MaxUsers
		game.Winner = game.Users[winner]
	}

	// Salvando
	h.cache.Set(game.ID, game, 5*time.Minute)

	finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
	logger.Log(map[string]interface{}{
		"_action":   "GamesJoin.Find",
		"_result":   "success",
		"_duration": finishTime - initTime,
	})

	logger.Log(map[string]interface{}{
		"_action":       "GamesJoin",
		"_rid":          c.Get(echo.HeaderXRequestID),
		"_real-ip":      c.RealIP,
		"_duration":     finishTime - initTime,
		"_result":       "success",
		"short_message": "Join Games",
	})

	return c.JSON(http.StatusOK, game)
}
