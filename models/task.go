package models



type Task struct {
	Guid        string `json:"id"`
	Description string `json:"description"  binding:"required"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

