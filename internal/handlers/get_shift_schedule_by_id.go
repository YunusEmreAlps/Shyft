package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shyft/internal/models"
	"shyft/pkg/httpErrors"
)

// HandleGetShiftScheduleByID godoc
// HandleGetShiftScheduleByID handles the request to get a shift schedule by id
// @Summary get a shift schedule by id
// @Schemes
// @Description get a shift schedule by id
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shift Schedule ID"
// @Success 200 {object} RespondJson "get shift by id successfully"
// @Failure 400 {object} RespondJson "cannot get shift schedule by id due to invalid request body"
// @Failure 422 {object} RespondJson "cannot get shift schedule by id due to invalid request body"
// @Failure 500 {object} RespondJson "cannot get shift schedule by id due to internal server error"
// @Router /shift-schedules/id [get]
func (ss *ShiftService) HandleGetShiftScheduleByID(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get shift id from path and validate
	id := c.Param("id")
	if id == "" {
		return http.StatusBadRequest, nil, nil
	}

	// Step 2: Get shift schedule by id from database
	var shiftSchedule models.ShiftSchedule
	if err := ss.db.Where("deleted_at IS NULL AND id = ?", id).Find(&shiftSchedule).Error; err != nil {
		r, i := httpErrors.ErrorResponse(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r, i, errors.New("cannot get shift schedule by id due to not found")
		}
		return r, i, errors.New("cannot get shift schedule by id due to internal server error")
	}

	// Step 3: Return shift schedule by id
	return http.StatusOK, shiftSchedule, nil
}
