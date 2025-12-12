package models

type ListParams struct {
	Page     int    `form:"page" default:"1"`
	PageSize int    `form:"page_size" default:"10"`
	Search   string `form:"search"`

	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`

	OnlyActive        *bool  `form:"active" default:"true"`
	Status            *int   `form:"status"`
	Year              *int   `form:"year"`
	StartDate         string `form:"start_date"`
	EndDate           string `form:"end_date"`
	OrganizationID    *int   `form:"organization_id"`
	ManagerID         *int   `form:"manager_id"`
	OrganizationName  string `form:"organization_name"`
	OrganizationMail  string `form:"organization_mail"`
	OrganizationPhone string `form:"organization_phone"`
	ManagerName       string `form:"manager_name"`
	ManagerMail       string `form:"manager_mail"`
	ManagerPhone      string `form:"manager_phone"`
	UserID            *int   `form:"user_id"`
	UserName          string `form:"user_name"`
	UserMail          string `form:"user_mail"`
	UserPhone         string `form:"user_phone"`
	ShiftID           *int   `form:"shift_id"`
	ShiftStart        string `form:"shift_start"`
	ShiftEnd          string `form:"shift_end"`
	ShiftUser         string `form:"shift_user"`
	RangeYear         *int   `form:"range_year"`
}

func (p *ListParams) GetSortString() string {
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortOrder == "" {
		p.SortOrder = "DESC"
	}

	if p.SortBy == "organization.name" {
		return "(organization->0->>'name') " + p.SortOrder
	}
	if p.SortBy == "organization.mail" {
		return "(organization->0->>'mail') " + p.SortOrder
	}
	if p.SortBy == "organization.phone" {
		return "(organization->0->>'phone') " + p.SortOrder
	}
	if p.SortBy == "manager.name" {
		return "(manager->0->>'name') " + p.SortOrder
	}
	if p.SortBy == "manager.mail" {
		return "(manager->0->>'mail') " + p.SortOrder
	}
	if p.SortBy == "manager.phone" {
		return "(manager->0->>'phone') " + p.SortOrder
	}
	if p.SortBy == "user.name" {
		return "(users->0->>'name') " + p.SortOrder
	}
	if p.SortBy == "user.mail" {
		return "(users->0->>'mail') " + p.SortOrder
	}
	if p.SortBy == "user.phone" {
		return "(users->0->>'phone') " + p.SortOrder
	}
	if p.SortBy == "shift.start" {
		return "(shifts->0->>'start') " + p.SortOrder
	}
	if p.SortBy == "shift.end" {
		return "(shifts->0->>'end') " + p.SortOrder
	}
	if p.SortBy == "shift.user" {
		return "(shifts->0->>'user') " + p.SortOrder
	}

	jsonbColumns := map[string]string{
		"organization_name":  "organization->>'name'",
		"organization_mail":  "organization->>'mail'",
		"organization_phone": "organization->>'phone'",
		"manager_name":       "manager->>'name'",
		"manager_mail":       "manager->>'mail'",
		"manager_phone":      "manager->>'phone'",
		"user_id":            "users->>'id'",
		"user_name":          "users->>'name'",
		"user_mail":          "users->>'mail'",
		"user_phone":         "users->>'phone'",
		"shift_id":           "shifts->>'id'",
		"shift_start":        "shifts->>'start'",
		"shift_end":          "shifts->>'end'",
		"shift_user":         "shifts->>'user'",
	}

	allowedColumns := map[string]string{
		"ID":          "id",
		"alias":       "alias",
		"year":        "year",
		"start_date":  "start_date",
		"end_date":    "end_date",
		"status":      "status",
		"frequency":   "frequency",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"description": "description",
	}

	allowedOrders := map[string]bool{
		"ASC":  true,
		"DESC": true,
		"asc":  true,
		"desc": true,
	}

	var column string
	if jsonbCol, exists := jsonbColumns[p.SortBy]; exists {
		column = jsonbCol
	} else if normalCol, exists := allowedColumns[p.SortBy]; exists {
		column = normalCol
	} else {
		column = "created_at"
	}

	order := "DESC"
	if allowedOrders[p.SortOrder] {
		order = p.SortOrder
	}

	return column + " " + order
}
