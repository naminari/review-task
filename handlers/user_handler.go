package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) SetUserActiveHandler(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	user, err := app.Service.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "user not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (app *App) GetUserReviewHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "user_id parameter is required",
			},
		})
		return
	}

	prs, err := app.Service.GetUserReviewPRs(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
