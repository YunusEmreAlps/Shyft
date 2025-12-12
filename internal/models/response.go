package models

type ShiftScheduleListResponse struct {
	Data       []ShiftSchedule `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
	SortBy     string          `json:"sort_by"`
	SortOrder  string          `json:"sort_order"`
}
