package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shyft/internal/models"
	"shyft/pkg/httpErrors"
)

// HandleDeleteShiftSchedule godoc
// HandleDeleteShiftSchedule handles the request to delete a shift schedule
// @Summary delete a shift schedule
// @Schemes
// @Description delete a shift schedule
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shift Schedule ID"
// @Success 200 {object} RespondJson "successfully deleted shift schedule"
// @Failure 400 {object} RespondJson "cannot delete shift schedule due to invalid request body"
// @Failure 422 {object} RespondJson "cannot delete shift schedule due to invalid request body"
// @Failure 500 {object} RespondJson "cannot delete shift schedule due to internal server error"
// @Router /shift-schedules/{id} [delete]
func (ss *ShiftService) HandleDeleteShiftSchedule(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get shift schedule id from path	and validate
	id := c.Param("id")
	if id == "" {
		return http.StatusBadRequest, nil, nil
	}

	// Step 2: Delete shift schedule by id from database (soft delete)
	var shiftSchedule models.ShiftSchedule
	if err := ss.db.Where("id = ?", id).Delete(&shiftSchedule).Error; err != nil {
		r, i := httpErrors.ErrorResponse(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r, i, errors.New("cannot delete shift schedule due to not found")
		}
		return r, i, errors.New("cannot delete shift schedule due to internal server error")
	}

	// Step 3: Return shift schedule by id
	return http.StatusOK, "Shift Schedule Successfully Deleted", nil
}
