package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Region   string `json:"region"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Region   string `json:"region"`
	Token    string `json:"token"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.POST("/login", h.login)
	}
}

// Fake login: just echoes back role and region, generates a dummy token.
func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))
	req.Region = strings.ToLower(strings.TrimSpace(req.Region))
	if req.Username == "" || req.Role == "" || req.Region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, role, region required"})
		return
	}
	resp := LoginResponse{
		Username: req.Username,
		Role:     req.Role,
		Region:   req.Region,
		Token:    "dummy-token-" + req.Username,
	}
	c.JSON(http.StatusOK, resp)
}


