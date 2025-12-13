package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shyft/internal/models"
	"shyft/pkg/httpErrors"
)

type createShiftScheduleDTO struct {
	Alias        string       `json:"alias" binding:"required"`
	Description  string       `json:"description" binding:"required"`
	Frequency    int          `json:"frequency" binding:"required"` // 1, 2, 3, 4, 5, 6, 7 (days of the week)
	Start_Date   time.Time    `json:"start_date" binding:"required"`
	End_Date     time.Time    `json:"end_date" binding:"required"`
	Year         int          `json:"year" binding:"required"`
	Status       int          `json:"status"` // 0: pending, 1: approved, 2: rejected
	Organization models.JSONB `json:"organization" binding:"required"`
	Manager      models.JSONB `json:"manager" binding:"required"`
	Users        models.JSONB `json:"users" binding:"required"`
	Shifts       models.JSONB `json:"shifts" binding:"required"`
}

// HandleCreateShiftSchedule godoc
// HandleCreateShiftSchedule handles the request to create a new shift schedule
// @Summary create a new shift schedule
// @Schemes
// @Description create a new shift schedule
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param alias body string true "Alias"
// @Param description body string true "Description"
// @Param frequency body int true "Frequency"
// @Param start_date body string true "Start Date"
// @Param end_date body string true "End Date"
// @Param year body int true "Year"
// @Param status body int true "Status"
// @Param organization body object true "Organization"
// @Param manager body object true "Manager"
// @Param users body object true "Users"
// @Param shifts body object true "Shifts"
// @Success 200 {object} RespondJson "successfully created shift schedule"
// @Failure 400 {object} RespondJson "cannot create shift schedule due to invalid request body"
// @Failure 422 {object} RespondJson "cannot create shift schedule due to invalid request body"
// @Failure 500 {object} RespondJson "cannot create shift schedule due to internal server error"
// @Router /shift-schedules [post]
func (ss *ShiftService) HandleCreateShiftSchedule(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get shift schedule from request body
	var params createShiftScheduleDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		return http.StatusBadRequest, nil, err
	}
	// if err := utils.ValidateStruct(c, &params); err != nil {
	// 	return http.StatusBadRequest, nil, err
	// }

	// Step 2: Validate shift schedule
	var shiftSchedule models.ShiftSchedule
	createParamsToShiftSchedule(&params, &shiftSchedule)

	// Step 3: Create shift schedule in database
	if err := ss.db.Create(&shiftSchedule).Error; err != nil {
		r, i := httpErrors.ErrorResponse(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r, i, errors.New("cannot create shift schedule due to not found")
		}
		return r, i, errors.New("cannot create shift schedule due to internal server error")
	}

	// Step 4: Return shift schedule
	return http.StatusOK, "Shift Schedule Successfully Created", nil
}

func createParamsToShiftSchedule(params *createShiftScheduleDTO, shift *models.ShiftSchedule) {
	shift.Alias = params.Alias
	shift.Frequency = params.Frequency
	shift.Description = params.Description
	shift.Start_Date = params.Start_Date
	shift.End_Date = params.End_Date
	shift.Year = params.Year
	shift.Status = params.Status
	shift.Organization = params.Organization
	shift.Manager = params.Manager
	shift.Users = params.Users
	shift.Shifts = params.Shifts
}
