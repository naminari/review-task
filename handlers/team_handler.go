package handlers

import (
	"net/http"
	"reviewtask/models"

	"github.com/gin-gonic/gin"
)

func (app *App) CreateTeamHandler(c *gin.Context) {
	var team models.Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	if err := app.Service.CreateTeam(&team); err != nil {
		if err.Error() == "team_name already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": map[string]interface{}{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": team,
	})
}

func (app *App) GetTeamHandler(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "team_name parameter is required",
			},
		})
		return
	}

	team, err := app.Repo.GetTeam(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "team not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, team)
}
