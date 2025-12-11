package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"geo-repl-demo/internal/model"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.GET("/articles", h.list)
		api.POST("/articles", h.create)
		api.DELETE("/articles/:id", h.delete)
		api.GET("/replication-status", h.status)
	}
}

func (h *Handler) list(c *gin.Context) {
	// Önce query parameter'dan region'ı kontrol et (test için)
	regionStr := c.Query("region")
	
	// Query parameter yoksa middleware'den gelen region'ı kullan
	if regionStr == "" {
		region, _ := c.Get("region")
		if val, ok := region.(string); ok {
			regionStr = val
		} else {
			regionStr = "eu" // Varsayılan
		}
	}

	arts, err := h.svc.ListByRegion(c.Request.Context(), regionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, arts)
}

func (h *Handler) create(c *gin.Context) {
	var in model.CreateArticleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) status(c *gin.Context) {
	status, err := h.svc.ReplicationStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}
