package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/api/handlers"
	"eino-researcher/internal/config"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/store"
)

type Dependencies struct {
	Config       config.Config
	Documents    store.DocumentRepository
	RAG          rag.Service
	Orchestrator agent.Orchestrator
}

func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(RequestLogger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		documentHandler := handlers.NewDocumentHandler(deps.Documents)
		ragHandler := handlers.NewRAGHandler(deps.RAG)
		researchHandler := handlers.NewResearchHandler(deps.Orchestrator)

		v1.POST("/documents", documentHandler.Upload)
		v1.POST("/rag/query", ragHandler.Query)
		v1.POST("/research/tasks", researchHandler.CreateTask)
		v1.GET("/research/tasks/:task_id", researchHandler.GetTask)
		v1.GET("/research/tasks/:task_id/report", researchHandler.GetReport)
		v1.GET("/research/tasks/:task_id/stream", researchHandler.Stream)
	}

	return router
}
