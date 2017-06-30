package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mdjunior/ct/logger"
	"github.com/mdjunior/modulus/models"
)

func main() {
	initTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0

	// Verificando se nome do jogador foi informado
	if len(os.Args) < 2 {
		fmt.Println("Você deve informar o nome do jogador rodando: " + os.Args[0] + " playerName")
		return
	}

	// Jogando infinitamente
	for {
		// Verificando status do serviço
		statusReq, _ := http.NewRequest(http.MethodGet, os.Getenv("API_ENDPOINT")+"/status", nil)
		statusRes, err := http.DefaultClient.Do(statusReq)
		if err != nil || statusRes.StatusCode != http.StatusOK {
			finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
			logger.Log(map[string]interface{}{
				"short_message": err.Error,
				"_action":       "api.check",
				"_result":       "err",
				"_duration":     finishTime - initTime,
			})
			fmt.Println("Cadê o arbitro?")
		}
		defer statusRes.Body.Close()

		statusResBody, _ := ioutil.ReadAll(statusRes.Body)
		if string(statusResBody) != "WORKING" {
			finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
			logger.Log(map[string]interface{}{
				"short_message": err.Error,
				"_action":       "api.check.body",
				"_result":       "err",
				"_duration":     finishTime - initTime,
			})
			fmt.Println("Deu erro no status")
		}

		// Verificando games disponíveis
		gamesReq, _ := http.NewRequest(http.MethodGet, os.Getenv("API_ENDPOINT")+"/games", nil)
		gamesRes, err := http.DefaultClient.Do(gamesReq)
		if err != nil || gamesRes.StatusCode != http.StatusOK {
			finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
			logger.Log(map[string]interface{}{
				"short_message": err.Error,
				"_action":       "api.games",
				"_result":       "err",
				"_duration":     finishTime - initTime,
			})
			fmt.Println("O arbitro não te respondeu quais jogos estão rolando...")
		}
		defer gamesRes.Body.Close()
		gamesResBody, _ := ioutil.ReadAll(gamesRes.Body)

		// Fazendo unmarshall do JSON
		var games []models.Game
		if err := json.Unmarshal(gamesResBody, &games); err != nil {
			finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
			logger.Log(map[string]interface{}{
				"short_message": err.Error,
				"_action":       "api.games.unmarshall",
				"_result":       "err",
				"_duration":     finishTime - initTime,
			})
		}

		// Verificando se tem algum jogo aberto
		var openGameID string
		for i := range games {
			if games[i].Status == "OPEN" {
				openGameID = games[i].ID
			}
		}
		if openGameID == "" {
			// vamos criar um jogo
			gamesNewReq, _ := http.NewRequest(http.MethodPost, os.Getenv("API_ENDPOINT")+"/games", bytes.NewBuffer([]byte("{}")))
			gamesNewReq.Header.Set("Content-Type", "application/json")
			gamesNewRes, err := http.DefaultClient.Do(gamesNewReq)
			if err != nil || gamesNewRes.StatusCode != http.StatusOK {
				finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
				logger.Log(map[string]interface{}{
					"short_message": err.Error,
					"_action":       "api.newgame",
					"_result":       "err",
					"_duration":     finishTime - initTime,
				})
				fmt.Println("O arbitro não deixou você iniciar um novo jogo...")
			}
			defer gamesNewRes.Body.Close()
			gamesNewResBody, _ := ioutil.ReadAll(gamesNewRes.Body)
			// Fazendo unmarshall do JSON
			var newGame models.Game
			if err := json.Unmarshal(gamesNewResBody, &newGame); err != nil {
				finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
				logger.Log(map[string]interface{}{
					"short_message": err.Error,
					"_action":       "api.newgame.unmarshall",
					"_result":       "err",
					"_duration":     finishTime - initTime,
				})
			}
			openGameID = newGame.ID
		}

		// Jogo criado, vamos solicitar ao usuário informar sua jogada
		fmt.Println("Você está participando do jogo:" + openGameID + ". Qual sua jogada?")
		var attempt string
		fmt.Scan(&attempt)

		// Verificando se string é int
		var try models.Try
		try.Value, err = strconv.Atoi(attempt)
		if err != nil {
			fmt.Println("Sua jogada tem que ser um número. Tente novamente:")
			fmt.Scan(&attempt)
		}
		try.Value, err = strconv.Atoi(attempt)
		if err != nil {
			fmt.Println("Parece que você não quer jogar... vou sair...")
			return
		}

		// Preparando request
		try.Name = os.Args[1]

		// Vamos fazer a jogada
		tryJSON, err := json.Marshal(try)
		if err != nil {
			fmt.Println("Não conseguimos preparar a jogada, tente de novo...")
			return
		}
		tryReq, _ := http.NewRequest(http.MethodPost, os.Getenv("API_ENDPOINT")+"/games/"+openGameID+"/join", bytes.NewBuffer(tryJSON))
		tryReq.Header.Set("Content-Type", "application/json")
		tryRes, err := http.DefaultClient.Do(tryReq)
		if err != nil || tryRes.StatusCode != http.StatusOK {
			fmt.Println("O arbitro não aceitou sua jogada...")
			break
		}
		defer tryRes.Body.Close()
		tryResBody, _ := ioutil.ReadAll(tryRes.Body)
		// Fazendo unmarshall do JSON
		var newGameAfterTry models.Game
		if err := json.Unmarshal(tryResBody, &newGameAfterTry); err != nil {
			finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
			logger.Log(map[string]interface{}{
				"short_message": err.Error,
				"_action":       "api.newgameaftertry.unmarshall",
				"_result":       "err",
				"_duration":     finishTime - initTime,
			})
		}

		// Verificando se ha um vencedor
		if len(newGameAfterTry.Winner) == 1 {
			fmt.Println(newGameAfterTry.Winner)
		} else {
			// Verificando até o jogo ter encerrado
			retry := true
			for retry {
				gamesRetryReq, _ := http.NewRequest(http.MethodGet, os.Getenv("API_ENDPOINT")+"/games/"+openGameID, nil)
				gamesRetryRes, err := http.DefaultClient.Do(gamesRetryReq)
				if err != nil || gamesRetryRes.StatusCode != http.StatusOK {
					finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
					logger.Log(map[string]interface{}{
						"short_message": err.Error,
						"_action":       "api.gamesretry",
						"_result":       "err",
						"_duration":     finishTime - initTime,
					})
					fmt.Println("O arbitro esqueceu esse jogo...")
				}
				defer gamesRetryRes.Body.Close()
				// Fazendo unmarshall do JSON
				gamesRetryResBody, _ := ioutil.ReadAll(gamesRetryRes.Body)
				var retryGame models.Game
				if err := json.Unmarshal(gamesRetryResBody, &retryGame); err != nil {
					finishTime := float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0
					logger.Log(map[string]interface{}{
						"short_message": err.Error,
						"_action":       "api.gameretry.unmarshall",
						"_result":       "err",
						"_duration":     finishTime - initTime,
					})
				}
				if len(retryGame.Winner) == 1 {
					for key := range retryGame.Winner {
						fmt.Println("O vencedor foi:" + key + "!")
					}
					retry = false
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}
