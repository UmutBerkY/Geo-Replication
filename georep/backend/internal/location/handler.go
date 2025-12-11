package location

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler exposes HTTP endpoints for locations.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts the location routes on a Gin router.
func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.POST("/locations", h.createLocation)
		api.GET("/locations/master", h.listMaster)
		api.GET("/locations/closest", h.listClosest)
		api.GET("/locations/replica/:index", h.listReplica)
		api.GET("/health", h.health)
	}
}

func (h *Handler) createLocation(c *gin.Context) {
	var in CreateLocationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loc, err := h.svc.CreateLocation(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loc)
}

func (h *Handler) listMaster(c *gin.Context) {
	locs, err := h.svc.ListFromMaster(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locs)
}

func (h *Handler) listReplica(c *gin.Context) {
	indexStr := c.Param("index")
	i, err := strconv.Atoi(indexStr)
	if err != nil || i < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid replica index"})
		return
	}
	// path uses 1-based index, repository uses 0-based.
	locs, err := h.svc.ListFromReplica(c.Request.Context(), i-1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locs)
}

// listClosest returns locations from the node that would be considered
// "closest" to the user, based on a simple region hint.
// Regions:
//   - eu   -> master (single primary in EU)
//   - us   -> replica1 (US replica)
//   - apac -> replica2 (APAC replica)
func (h *Handler) listClosest(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		region = "eu"
	}

	var (
		locs   []Location
		err    error
		source string
	)

	switch region {
	case "eu":
		// EU users talk directly to the master (same region).
		locs, err = h.svc.ListFromMaster(c.Request.Context())
		source = "master"
	case "us":
		// US users read from the US replica (replica #1 in this demo).
		locs, err = h.svc.ListFromReplica(c.Request.Context(), 0)
		source = "replica1"
	case "apac":
		// APAC users read from the APAC replica (replica #2 in this demo).
		locs, err = h.svc.ListFromReplica(c.Request.Context(), 1)
		source = "replica2"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown region"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For the frontend we only return the list of stories here; which physical
	// node served the request is explained in the UI copy based on region.
	_ = source
	c.JSON(http.StatusOK, locs)
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}


