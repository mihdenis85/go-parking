package api

import (
	"net/http"
	"strconv"

	"github.com/amend-parking-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(svc *service.Service) *Handlers {
	return &Handlers{service: svc}
}

func (h *Handlers) GetCountOfFreeSpaces(c *gin.Context) {
	count, err := h.service.GetCountOfFreeSpaces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, count)
}

func (h *Handlers) GetOccupiedSpaces(c *gin.Context) {
	spaces, err := h.service.GetOccupiedSpaces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spaces)
}

func (h *Handlers) ParkCar(c *gin.Context) {
	var body AddParkingSpaceLogSchema
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.service.AddParkingSpaceLog(
		c.Request.Context(),
		body.FirstName,
		body.LastName,
		body.CarMake,
		body.LicensePlate,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "no free parking spaces available" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

func (h *Handlers) FreeUpParkingSpace(c *gin.Context) {
	placeNumberStr := c.Query("place_number")
	placeNumber, err := strconv.Atoi(placeNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place_number"})
		return
	}

	log, err := h.service.FreeUpParkingSpace(c.Request.Context(), placeNumber)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "parking space is already free" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

func (h *Handlers) GetParkingSpaceLogs(c *gin.Context) {
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	if firstName == "" || lastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "first_name and last_name are required"})
		return
	}

	logs, err := h.service.GetParkingSpaceLogsByFirstNameAndLastName(
		c.Request.Context(),
		firstName,
		lastName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
