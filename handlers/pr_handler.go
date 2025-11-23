package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) CreatePRHandler(c *gin.Context) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
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

	log.Printf("Creating PR: id=%s, name=%s, author_id=%s",
		req.PullRequestID, req.PullRequestName, req.AuthorID)

	pr, err := app.Service.CreatePRWithReviewers(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch err.Error() {
		case "PR id already exists":
			c.JSON(http.StatusConflict, gin.H{
				"error": map[string]interface{}{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
		case "author not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": map[string]interface{}{
					"code":    "NOT_FOUND",
					"message": "author/team not found",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		}
		return
	}

	log.Printf("PR created successfully: ID=%s", pr.PullRequestID)
	c.JSON(http.StatusCreated, gin.H{
		"pr": pr,
	})
}

func (app *App) MergePRHandler(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
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

	pr, err := app.Service.MergePR(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "PR not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": pr,
	})
}

func (app *App) ReassignReviewerHandler(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
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

	newReviewer, err := app.Service.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err.Error() {
		case "cannot reassign on merged PR":
			c.JSON(http.StatusConflict, gin.H{
				"error": map[string]interface{}{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		case "reviewer is not assigned to this PR":
			c.JSON(http.StatusConflict, gin.H{
				"error": map[string]interface{}{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
		case "no active replacement candidate in team":
			c.JSON(http.StatusConflict, gin.H{
				"error": map[string]interface{}{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
				},
			})
		case "PR not found", "old reviewer not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": map[string]interface{}{
					"code":    "NOT_FOUND",
					"message": err.Error(),
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		}
		return
	}

	// Получаем обновленный PR
	pr, err := app.Repo.GetPR(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "failed to get updated PR",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}
