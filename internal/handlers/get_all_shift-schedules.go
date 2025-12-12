package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"shift-scheduler-service/internal/models"
	"shift-scheduler-service/internal/repository"
	"shift-scheduler-service/pkg/httpErrors"
)

// HandleGetAllShiftSchedules godoc
// HandleGetAllShiftSchedules handles the request to get all shift schedules with advanced filtering, sorting and pagination
// @Summary get all shift schedules with filters
// @Description get all shift schedules with pagination, search, sort and filter options
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1) default(1)
// @Param page_size query int false "Page size" minimum(1) maximum(100) default(10)
// @Param search query string false "Search in alias, description, organization, manager"
// @Param sort_by query string false "Sort by field (alias, year, start_date, end_date, status, organization_name, manager_name)" default(created_at)
// @Param sort_order query string false "Sort order (ASC or DESC)" default(DESC)
// @Param active query bool false "Filter only active records (not deleted)"
// @Param status query int false "Filter by status (0:pending, 1:approved, 2:rejected)"
// @Param year query int false "Filter by year"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param organization_name query string false "Filter by organization name"
// @Param organization_mail query string false "Filter by organization mail"
// @Param organization_phone query string false "Filter by organization phone"
// @Param manager_name query string false "Filter by manager name"
// @Param manager_mail query string false "Filter by manager mail"
// @Param manager_phone query string false "Filter by manager phone"
// @Param user_id query int false "Filter by user ID"
// @Param user_name query string false "Filter by user name"
// @Param user_mail query string false "Filter by user mail"
// @Param user_phone query string false "Filter by user phone"
// @Param shift_id query int false "Filter by shift ID"
// @Param shift_start query string false "Filter by shift start (YYYY-MM-DD HH:MM:SS)"
// @Param shift_end query string false "Filter by shift end (YYYY-MM-DD HH:MM:SS)"
// @Param shift_user query string false "Filter by shift user"
// @Success 200 {object} models.ShiftScheduleListResponse "get all shift schedules successfully"
// @Failure 400 {object} RespondJson "cannot get all shift schedules due to invalid request parameters"
// @Failure 500 {object} RespondJson "cannot get all shift schedules due to internal server error"
// @Router /shift-schedules [get]
func (ss *ShiftService) HandleGetAllShiftSchedules(c *gin.Context) (int, interface{}, error) {
	// Step 1: Parse query parameters
	var params models.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return http.StatusBadRequest, nil, errors.New("invalid query parameters: " + err.Error())
	}

	// Step 2: Set default values
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	// Step 3: Use repository for listing
	repo := repository.NewShiftScheduleRepository(ss.db)
	result, err := repo.List(params)
	if err != nil {
		r, i := httpErrors.ErrorResponse(err)
		return r, i, errors.New("cannot get shift schedules: " + err.Error())
	}
	return http.StatusOK, result, nil
}
