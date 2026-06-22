package handlers

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
	"eino-researcher/internal/rag"
)

type DocumentHandler struct {
	indexer  rag.Indexer
	maxBytes int64
}

func NewDocumentHandler(indexer rag.Indexer, maxBytes int64) *DocumentHandler {
	if maxBytes <= 0 {
		maxBytes = 2 * 1024 * 1024
	}
	return &DocumentHandler{indexer: indexer, maxBytes: maxBytes}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	if file.Size > h.maxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file exceeds maximum upload size"})
		return
	}

	fileType := strings.ToLower(filepath.Ext(file.Filename))
	if fileType != ".md" && fileType != ".txt" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .md and .txt files are supported"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read uploaded file"})
		return
	}
	defer opened.Close()

	content, err := io.ReadAll(io.LimitReader(opened, h.maxBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read uploaded file"})
		return
	}
	if int64(len(content)) > h.maxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file exceeds maximum upload size"})
		return
	}
	if !utf8.Valid(content) || strings.TrimSpace(string(content)) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file must contain non-empty UTF-8 text"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		title = file.Filename
	}

	doc := model.NewDocument(title, file.Filename, fileType)
	if err := h.indexer.Index(c.Request.Context(), doc, string(content)); err != nil {
		if errors.Is(err, llm.ErrProvider) {
			c.JSON(http.StatusBadGateway, gin.H{"error": "model provider request failed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to index document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"document_id": doc.ID,
		"title":       doc.Title,
		"status":      "indexed",
	})
}
