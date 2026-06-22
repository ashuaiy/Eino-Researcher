package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/model"
	"eino-researcher/internal/store"
)

type DocumentHandler struct {
	repo store.DocumentRepository
}

func NewDocumentHandler(repo store.DocumentRepository) *DocumentHandler {
	return &DocumentHandler{repo: repo}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		title = file.Filename
	}

	fileType := filepath.Ext(file.Filename)
	doc := model.NewDocument(title, file.Filename, fileType)

	// TODO: persist uploaded file content, parse Markdown/txt, chunk, embed, and index.
	if err := h.repo.Create(c.Request.Context(), doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"document_id": doc.ID,
		"title":       doc.Title,
		"status":      "accepted",
	})
}
