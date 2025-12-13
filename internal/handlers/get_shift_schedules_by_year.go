package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shyft/internal/models"
)

// HandleGetShiftScheduleByYear godoc
// HandleGetShiftScheduleByYear handles the request to get a shift schedules by year
// @Summary get a shift schedules by year
// @Schemes
// @Description get a shift schedules by year
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year path string true "Shift Schedule Year"
// @Success 200 {object} RespondJson "get shifts by year successfully"
// @Failure 400 {object} RespondJson "cannot get shifts schedule by year due to invalid request body"
// @Failure 422 {object} RespondJson "cannot get shifts schedule by year due to invalid request body"
// @Failure 500 {object} RespondJson "cannot get shifts schedule by year due to internal server error"
// @Router /shift-schedules/{year} [get]
func (ss *ShiftService) HandleGetShiftScheduleByYear(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get shift year from path and validate
	year := c.Param("year")
	if year == "" {
		return http.StatusBadRequest, nil, nil
	}

	// Step 2: Get shift schedules by year from database
	var shiftSchedules []models.ShiftSchedule
	if err := ss.db.Where("deleted_at IS NULL AND year = ?", year).Find(&shiftSchedules).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("cannot get shift schedule by year due to not found")
		}
		return http.StatusInternalServerError, nil, errors.New("cannot get shift schedule by year due to internal server error")
	}

	// Step 3: Return shift schedules by year
	return http.StatusOK, shiftSchedules, nil
}
