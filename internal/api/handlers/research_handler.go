package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/utils"
)

type ResearchHandler struct {
	orchestrator agent.Orchestrator
}

type CreateResearchTaskRequest struct {
	Question        string `json:"question" binding:"required"`
	UseWebSearch    bool   `json:"use_web_search"`
	MaxSubQuestions int    `json:"max_sub_questions"`
}

func NewResearchHandler(orchestrator agent.Orchestrator) *ResearchHandler {
	return &ResearchHandler{orchestrator: orchestrator}
}

func (h *ResearchHandler) CreateTask(c *gin.Context) {
	var req CreateResearchTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.MaxSubQuestions <= 0 {
		req.MaxSubQuestions = 5
	}

	task, err := h.orchestrator.CreateTask(c.Request.Context(), agent.CreateTaskRequest{
		Question:        req.Question,
		UseWebSearch:    req.UseWebSearch,
		MaxSubQuestions: req.MaxSubQuestions,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create research task"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"task_id": task.ID,
		"status":  task.Status,
	})
}

func (h *ResearchHandler) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")
	task, err := h.orchestrator.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *ResearchHandler) GetReport(c *gin.Context) {
	taskID := c.Param("task_id")
	report, err := h.orchestrator.GetReport(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *ResearchHandler) Stream(c *gin.Context) {
	taskID := c.Param("task_id")
	writer := utils.NewSSEWriter(c.Writer)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// TODO: replace placeholder events with Redis/task-event backed streaming.
	_ = writer.Event("step", gin.H{
		"step":    "queued",
		"message": "research task stream is not implemented yet",
	})
	_ = writer.Event("done", gin.H{"task_id": taskID})
}
