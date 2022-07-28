package models

type Task struct {
	Guid        string `json:"id"`
	Description string `json:"description" form:"description"  binding:"required"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	CompletedAt string `json:"completedAt"`
	Userid      string `json:"user_id"`
}
