package handlers

import (
	"log"
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
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	log.Printf("Creating PR: title=%s, author_id=%d", req.Title, req.AuthorID)

	author, err := app.Repo.GetUserByID(req.AuthorID)
	if err != nil {
		log.Printf("Author not found: %v", err)
		c.JSON(400, gin.H{"error": "author not found"})
		return
	}
	log.Printf("Author found: %s (team: %d)", author.Username, author.TeamID)

	reviewers, err := app.Service.AssignReviewers(author.ID, author.TeamID)
	if err != nil {
		log.Printf("Failed to assign reviewers: %v", err)
		c.JSON(500, gin.H{"error": "cant assign reviewers"})
		return
	}
	log.Printf("Assigned reviewers: %v", reviewers)

	pr := &models.PullRequest{
		Title:     req.Title,
		AuthorID:  author.ID,
		Status:    models.StatusOpen,
		Reviewers: reviewers,
	}

	if err := app.Repo.CreatePR(pr); err != nil {
		log.Printf("Failed to create PR: %v", err)
		c.JSON(500, gin.H{"error": "cant create pr"})
		return
	}

	log.Printf("PR created successfully: ID=%d", pr.ID)
	c.JSON(200, pr)
}

func (app *App) GetPRHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	pr, err := app.Repo.GetPRByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "pr not found"})
		return
	}

	c.JSON(200, pr)
}

func (app *App) MergePRHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	pr, err := app.Repo.GetPRByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "pr not found"})
		return
	}

	if pr.Status == models.StatusMerged {
		c.JSON(200, pr)
		return
	}

	pr.Status = models.StatusMerged
	if err := app.Repo.UpdatePR(pr); err != nil {
		c.JSON(500, gin.H{"error": "cant merge"})
		return
	}

	c.JSON(200, pr)
}
