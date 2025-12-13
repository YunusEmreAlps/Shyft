package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shyft/internal/models"
)

// HandleGetShiftScheduleByWeek godoc
// HandleGetShiftScheduleByWeek handles the request to get shift schedules by the current week
// @Summary get shift schedules by current week
// @Schemes
// @Description get shift schedules by current week
// @Tags Shift
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RespondJson "get shifts by current week successfully"
// @Failure 400 {object} RespondJson "cannot get shifts schedule by current week due to invalid request body"
// @Failure 422 {object} RespondJson "cannot get shifts schedule by current week due to invalid request body"
// @Failure 500 {object} RespondJson "cannot get shifts schedule by current week due to internal server error"
// @Router /shift-schedules/week [get]
func (ss *ShiftService) HandleGetShiftScheduleByWeek(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get all shift schedules
	var shiftSchedules []models.ShiftSchedule
	if err := ss.db.Where("deleted_at IS NULL").Find(&shiftSchedules).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Step 2: Get current date and time this week format
	today := time.Now()
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	if today.Weekday() == time.Sunday {
		weekStart = weekStart.AddDate(0, 0, -6)
	}
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekEnd = time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 0, weekEnd.Location())

	// Step 3: Filter shift schedules by current week
	var data []map[string]interface{}
	for _, shiftSchedule := range shiftSchedules {
		//
		temp := map[string]interface{}{
			"id":           shiftSchedule.ID,
			"alias":        shiftSchedule.Alias,
			"description":  shiftSchedule.Description,
			"frequency":    shiftSchedule.Frequency,
			"status":       shiftSchedule.Status,
			"week_start":   weekStart.Format(time.RFC3339),
			"week_end":     weekEnd.Format(time.RFC3339),
			"manager":      shiftSchedule.Manager,
			"organization": shiftSchedule.Organization,
			"users":        shiftSchedule.Users,
			"shifts":       []interface{}{},
		}

		for _, shift := range shiftSchedule.Shifts {
			// Unmarshal shift to map
			shiftMap, ok := shift.(map[string]interface{})
			if !ok {
				continue // Skip if shift is not a map
			}
			startTime, err := time.Parse(time.RFC3339, shiftMap["start"].(string))
			if err != nil {
				continue // Skip invalid date formats
			}
			endTime, err := time.Parse(time.RFC3339, shiftMap["end"].(string))
			if err != nil {
				continue // Skip invalid date formats
			}

			// Check if shift is in current week range
			if (startTime.After(weekStart) || startTime.Equal(weekStart)) && (endTime.Before(weekEnd) || endTime.Equal(weekEnd)) {
				temp["shifts"] = append(temp["shifts"].([]interface{}), shift)
			} else if (startTime.Before(weekStart) || startTime.Equal(weekStart)) && (endTime.After(weekEnd) || endTime.Equal(weekEnd)) {
				temp["shifts"] = append(temp["shifts"].([]interface{}), shift)
			} else if (startTime.After(weekStart) || startTime.Equal(weekStart)) && (startTime.Before(weekEnd) || startTime.Equal(weekEnd)) {
				temp["shifts"] = append(temp["shifts"].([]interface{}), shift)
			} else if (endTime.After(weekStart) || endTime.Equal(weekStart)) && (endTime.Before(weekEnd) || endTime.Equal(weekEnd)) {
				temp["shifts"] = append(temp["shifts"].([]interface{}), shift)
			}
		}

		// Only add organization data if there are shifts in the current week
		if len(temp["shifts"].([]interface{})) > 0 {
			data = append(data, temp)
		} else {
			// If empty or not all departments want to be shown
			data = append(data, map[string]interface{}{
				"id":           shiftSchedule.ID,
				"alias":        shiftSchedule.Alias,
				"description":  shiftSchedule.Description,
				"frequency":    shiftSchedule.Frequency,
				"status":       shiftSchedule.Status,
				"week_start":   weekStart.Format(time.RFC3339),
				"week_end":     weekEnd.Format(time.RFC3339),
				"organization": shiftSchedule.Organization,
				"manager":      shiftSchedule.Manager,
				"users":        shiftSchedule.Users,
				"shifts":       []interface{}{},
			})
		}
	}

	// Step 4: Return shift schedules by week
	return http.StatusOK, data, nil
}
