package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shyft/internal/models"
	"shyft/pkg/httpErrors"
)

// HandleGetOnlyDeletedShiftSchedules godoc
// HandleGetOnlyDeletedShiftSchedules handles the request to get all deleted shift schedules
// @Summary get all deleted shift schedules
// @Schemes
// @Description get all deleted shift schedules
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shift Schedule ID"
// @Success 200 {object} RespondJson "get all shift schedules successfully"
// @Failure 400 {object} RespondJson "cannot get only deleted shift schedules due to invalid request body"
// @Failure 422 {object} RespondJson "cannot get only deleted shift schedules due to invalid request body"
// @Failure 500 {object} RespondJson "cannot get only deleted shift schedules due to internal server error"
// @Router /shift-schedules/deleted [get]
func (ss *ShiftService) HandleGetOnlyDeletedShiftSchedules(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get all shift schedules from database
	var shiftSchedules []models.ShiftSchedule
	if err := ss.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&shiftSchedules).Error; err != nil {
		r, i := httpErrors.ErrorResponse(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r, i, errors.New("cannot get all shift schedules due to not found")
		}
		return r, i, errors.New("cannot get all shift schedules due to internal server error")
	}

	// Step 2: Return all shift schedules
	return http.StatusOK, shiftSchedules, nil
}
