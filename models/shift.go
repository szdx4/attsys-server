package models

// shift 排班模型
type Shift struct {
	CommonFields
	UserID  uint   `json:"user_id"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}
