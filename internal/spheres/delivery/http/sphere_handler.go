package http

import (
	"github.com/MangoSociety/know_api/internal/spheres/domain"
	"github.com/MangoSociety/know_api/internal/spheres/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type SphereHandler struct {
	sphereService service.SphereService
}

func NewSphereHandler(sphereService service.SphereService) *SphereHandler {
	return &SphereHandler{sphereService: sphereService}
}

func (h *SphereHandler) GetSpheres(c *gin.Context) {
	spheres, err := h.sphereService.GetSpheres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spheres)
}

func (h *SphereHandler) CreateSphere(c *gin.Context) {
	var sphere domain.Sphere
	if err := c.ShouldBindJSON(&sphere); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sphere.ID = primitive.NewObjectID()
	if err := h.sphereService.CreateSphere(sphere); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sphere)
}

func (h *SphereHandler) UpdateSphere(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var sphere domain.Sphere
	if err := c.ShouldBindJSON(&sphere); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sphere.ID = id
	if err := h.sphereService.UpdateSphere(sphere); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sphere)
}

func (h *SphereHandler) DeleteSphere(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.sphereService.DeleteSphere(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sphere deleted"})
}
