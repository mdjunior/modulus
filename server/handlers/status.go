package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdjunior/ct/logger"
)

// Status is a method that respond WORKING
func Status(echoContext echo.Context) error {
	logger.Log(map[string]interface{}{
		"_action":       "Status",
		"_rid":          echoContext.Get(echo.HeaderXRequestID),
		"_real-ip":      echoContext.RealIP,
		"_result":       "success",
		"short_message": "Verify Healthcheck",
	})

	return echoContext.String(http.StatusOK, "WORKING")
}
