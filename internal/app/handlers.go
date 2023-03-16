package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// indexHandler Handles the "/" route, returns 200 Status Code and simple string as a page content
func (a *App) indexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello from %v\n", "Gin")
}

// clearLimitsHandler Handles the "/clear" route.
// It cleans clientsRate field and returns 200 Status Code and simple string as a page content
func (a *App) clearLimitsHandler(c *gin.Context) {
	a.clientsRate = make(map[string]*client)
	c.String(http.StatusOK, "Limits are cleared\n")
}
