package v1

import (
	"encoding/json"
	"github.com/nordew/scope_test/internal/model"
	"github.com/nordew/scope_test/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessMsg = "Success"
)

type Handler struct {
	workerService service.WorkerService
}

func NewHandler(workerService service.WorkerService) *Handler {
	return &Handler{
		workerService: workerService,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	app := router.Group("/api")
	{
		app.POST("/", h.handlePostRequest)
	}

	return router
}

func (h *Handler) handlePostRequest(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	dataBytes, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	job := model.Job{
		Data: dataBytes,
	}
	h.workerService.Submit(job)

	c.JSON(http.StatusOK, gin.H{"message": SuccessMsg})
}
