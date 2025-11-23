package handlers

import (
	"net/http"
	"strconv"

	"reviewtask/models"

	"github.com/gin-gonic/gin"
)

func (app *App) CreatePRHandler(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		AuthorID int    `json:"author_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Получаем автора
	author, err := app.Repo.GetUserByID(req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
		return
	}

	// Назначаем ревьюеров
	reviewers, err := app.Service.AssignReviewers(author.ID, author.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign reviewers"})
		return
	}

	// Создаем PR
	pr := &models.PullRequest{
		Title:     req.Title,
		AuthorID:  author.ID,
		Status:    models.StatusOpen,
		Reviewers: reviewers,
	}

	if err := app.Repo.CreatePR(pr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PR"})
		return
	}

	c.JSON(http.StatusCreated, pr)
}

func (app *App) GetPRHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	pr, err := app.Repo.GetPRByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PR not found"})
		return
	}

	c.JSON(http.StatusOK, pr)
}

func (app *App) MergePRHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	pr, err := app.Repo.GetPRByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PR not found"})
		return
	}

	// Идемпотентность - если уже мерджнут, возвращаем успех
	if pr.Status == models.StatusMerged {
		c.JSON(http.StatusOK, pr)
		return
	}

	// Мерджим
	pr.Status = models.StatusMerged
	if err := app.Repo.UpdatePR(pr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to merge PR"})
		return
	}

	c.JSON(http.StatusOK, pr)
}
