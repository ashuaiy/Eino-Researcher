package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/rag"
)

type RAGHandler struct {
	service rag.Service
}

type RAGQueryRequest struct {
	Question string `json:"question" binding:"required"`
	TopK     int    `json:"top_k"`
	Stream   bool   `json:"stream"`
}

func NewRAGHandler(service rag.Service) *RAGHandler {
	return &RAGHandler{service: service}
}

func (h *RAGHandler) Query(c *gin.Context) {
	var req RAGQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TopK <= 0 {
		req.TopK = 5
	}

	answer, err := h.service.Query(c.Request.Context(), rag.QueryRequest{
		Question: req.Question,
		TopK:     req.TopK,
		Stream:   req.Stream,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query rag service"})
		return
	}

	c.JSON(http.StatusOK, answer)
}
