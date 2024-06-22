package http

import (
	"github.com/MangoSociety/know_api/internal/user/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service *service.UserSelectionService
}

func NewUserHandler(service *service.UserSelectionService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AddCategoryToUser(c *gin.Context) {
	var req struct {
		ChatID     int    `json:"chat_id"`
		SphereID   string `json:"sphere_id"`
		CategoryID string `json:"category_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//sphereID, err := primitive.ObjectIDFromHex(req.SphereID)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sphere_id"})
	//	return
	//}
	//
	//categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
	//	return
	//}

	//if err := h.service.AddCategoryToUser(req.ChatID, sphereID, categoryID); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"status": "category added"})
}
