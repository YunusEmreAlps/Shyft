package repository

import (
	"shift-scheduler-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type ShiftScheduleRepository struct {
	db *gorm.DB
}

func NewShiftScheduleRepository(db *gorm.DB) *ShiftScheduleRepository {
	return &ShiftScheduleRepository{db: db}
}

func (r *ShiftScheduleRepository) List(params models.ListParams) (*models.ShiftScheduleListResponse, error) {
	var schedules []models.ShiftSchedule
	var total int64

	query := r.db.Model(&models.ShiftSchedule{})

	// Apply all filters
	query = r.applyFilters(query, params)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	query = query.Order(params.GetSortString())

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&schedules).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.ShiftScheduleListResponse{
		Data:       schedules,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		SortBy:     params.SortBy,
		SortOrder:  params.SortOrder,
	}, nil
}

func (r *ShiftScheduleRepository) applyFilters(query *gorm.DB, params models.ListParams) *gorm.DB {
	// OnlyActive filter
	if params.OnlyActive != nil {
		if *params.OnlyActive {
			query = query.Where("deleted_at IS NULL")
		} else {
			query = query.Unscoped().Where("deleted_at IS NOT NULL")
		}
	}

	// Search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where(
			`alias ILIKE ? OR description ILIKE ? OR 
        EXISTS (SELECT 1 FROM jsonb_array_elements(organization) AS o WHERE o->>'name' ILIKE ?) OR
        EXISTS (SELECT 1 FROM jsonb_array_elements(manager) AS m WHERE m->>'name' ILIKE ?)`,
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Status filter
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// Year filter
	if params.Year != nil {
		query = query.Where("year = ?", *params.Year)
	}

	// Organization filters
	if params.OrganizationName != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM jsonb_array_elements(organization) AS o WHERE o->>'name' ILIKE ?)",
			"%"+params.OrganizationName+"%",
		)
	}
	if params.OrganizationMail != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(organization) AS o WHERE o->>'mail' ILIKE ?)", "%"+params.OrganizationMail+"%")
	}
	if params.OrganizationPhone != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(organization) AS o WHERE o->>'phone' ILIKE ?)", "%"+params.OrganizationPhone+"%")
	}

	// Manager filters
	if params.ManagerName != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM jsonb_array_elements(manager) AS m WHERE m->>'name' ILIKE ?)",
			"%"+params.ManagerName+"%",
		)
	}
	if params.ManagerMail != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM jsonb_array_elements(manager) AS m WHERE m->>'mail' ILIKE ?)",
			"%"+params.ManagerMail+"%",
		)
	}
	if params.ManagerPhone != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM jsonb_array_elements(manager) AS m WHERE m->>'phone' ILIKE ?)",
			"%"+params.ManagerPhone+"%",
		)
	}

	// User filters
	if params.UserID != nil {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(users) AS u WHERE (u->>'id')::int = ?)", *params.UserID)
	}
	if params.UserName != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM jsonb_array_elements(users) AS m WHERE m->>'name' ILIKE ?)",
			"%"+params.UserName+"%",
		)
	}
	if params.UserMail != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(users) AS m WHERE m->>'mail' ILIKE ?)", "%"+params.UserMail+"%")
	}
	if params.UserPhone != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(users) AS m WHERE m->>'phone' ILIKE ?)", "%"+params.UserPhone+"%")
	}

	// Shifts filters
	if params.ShiftID != nil {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(shifts) AS u WHERE (u->>'id')::int = ?)", *params.ShiftID)
	}

	if params.ShiftStart != "" {
		if shiftStart, err := time.Parse("2006-01-02", params.ShiftStart); err == nil {
			query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(shifts) AS s WHERE (s->>'start')::date >= ?)", shiftStart)
		}
	}

	if params.ShiftEnd != "" {
		if shiftEnd, err := time.Parse("2006-01-02", params.ShiftEnd); err == nil {
			query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(shifts) AS s WHERE (s->>'end')::date <= ?)", shiftEnd)
		}
	}

	if params.ShiftUser != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements(shifts) AS s WHERE s->>'user' ILIKE ?)", "%"+params.ShiftUser+"%")
	}

	// Date filters
	if params.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			query = query.Where("start_date >= ?", startDate)
		}
	}
	if params.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			query = query.Where("end_date <= ?", endDate)
		}
	}

	if params.RangeYear != nil {
		loc, _ := time.LoadLocation("Europe/Istanbul")
		yearStart := time.Date(*params.RangeYear, 1, 1, 0, 0, 0, 0, loc)
		yearEnd := time.Date(*params.RangeYear, 12, 31, 23, 59, 59, 0, loc)
		query = query.Where("start_date <= ? AND end_date >= ?", yearEnd, yearStart)
	}

	return query
}
