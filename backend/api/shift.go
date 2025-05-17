package api

import (
	"math/rand"
	"net/http"
	"payd-mini-project/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShiftAPI struct {
	db *gorm.DB
}

func NewShiftAPI(db *gorm.DB) *ShiftAPI {
	return &ShiftAPI{db: db}
}

func (h *ShiftAPI) RegisterRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	{
		// CRUD for shift
		api.POST("/shift", h.CreateShift)
		api.GET("/shifts", h.GetAllShifts)
		api.GET("/shift/:id", h.GetShift)
		api.PUT("/shift/:id", h.UpdateShift)
		api.DELETE("/shift/:id", h.DeleteShift)

		// Worker endpoint
		api.POST("/shift/request", h.RequestShift)

		// Admin endpoint
		api.PUT("/shift/:id/status", h.UpdateShiftStatus)
	}
}

func (h *ShiftAPI) CreateShift(c *gin.Context) {
	var shift model.Shift
	if err := c.ShouldBindJSON(&shift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrValidation]})
		return
	}

	if err := shift.CreateShiftValidation(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[err]})
		return
	}

	shift.ID = rand.Int63()
	shift.CreatedAt = time.Now().UTC()

	h.db.Create(&shift)

	c.JSON(http.StatusCreated, gin.H{"data": shift})
}

// Add this new handler for status updates
func (h *ShiftAPI) UpdateShiftStatus(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrValidation]})
		return
	}

	var shift model.Shift
	if h.db.First(&shift, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shift not found"})
		return
	}

	switch request.Status {
	case
		model.StatusApproved,
		model.StatusRejected:
		// Valid status
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrInvalidState]})
		return
	}

	if request.Status == model.StatusApproved {
		// Check if the shift is already assigned
		if shift.Assigned {
			c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrInvalidState]})
			return
		}

		// Check if the worker is available
		var existing []model.Shift
		h.db.Where("worker_id =? AND status =?", shift.WorkerID, model.StatusApproved).Find(&existing)

		if err := shift.VerifyOverlap(existing); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": model.ErrMessages[err]})
			return
		}

		if err := shift.VerifyShiftLimits(existing); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": model.ErrMessages[err]})
			return
		}

		shift.Assigned = true
	}

	shift.Status = request.Status

	h.db.Save(&shift).Where("id =?", id)

	c.JSON(http.StatusOK, shift)
}

// Update GetAllShifts to support filtering
func (h *ShiftAPI) GetAllShifts(c *gin.Context) {
	var shifts []model.Shift
	query := h.db.Model(&model.Shift{})

	if workerID := c.Query("worker"); workerID != "" {
		query = query.Where("worker_id =?", workerID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status =?", status)
	}

	query = query.Order("created_at DESC")
	query.Find(&shifts)

	c.JSON(http.StatusOK, gin.H{"data": shifts})
}

func (h *ShiftAPI) RequestShift(c *gin.Context) {
	var request struct {
		ShiftID  int64 `json:"shift_id"`
		WorkerID int64 `json:"worker_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrValidation]})
		return
	}

	var shift model.Shift
	if h.db.First(&shift, request.ShiftID).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shift not found"})
		return
	}

	var existingShifts []model.Shift
	h.db.Where("worker_id = ?", request.WorkerID).Find(&existingShifts)

	if err := shift.VerifyOverlap(existingShifts); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": model.ErrMessages[err]})
		return
	}

	if err := shift.VerifyShiftLimits(existingShifts); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": model.ErrMessages[err]})
		return
	}

	shift.WorkerID = request.WorkerID
	shift.Assigned = true
	shift.RequestedAt = time.Now().UTC()
	shift.Status = model.StatusPending

	h.db.Updates(&shift).Where("id = ?", shift.ID)

	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *ShiftAPI) GetShift(c *gin.Context) {
	id := c.Param("id")

	var shift model.Shift
	if err := h.db.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": model.ErrMessages[model.ErrValidation]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *ShiftAPI) UpdateShift(c *gin.Context) {
	id := c.Param("id")

	var shift model.Shift
	if err := c.ShouldBindJSON(&shift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": model.ErrMessages[model.ErrValidation]})
		return
	}

	h.db.Save(&shift).Where("id =?", id)

	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *ShiftAPI) DeleteShift(c *gin.Context) {
	id := c.Param("id")

	h.db.Delete(&model.Shift{}, id)

	c.JSON(http.StatusOK, gin.H{"message": "Shift deleted"})
}
