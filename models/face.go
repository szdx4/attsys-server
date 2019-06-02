package models

// Face 人脸模型
type Face struct {
	CommonFields
	UserID uint   `json:"user_id"`
	Info   string `json:"info" gorm:"type:text"`
	Status string `json:"status" gorm:"type:enum('wait','available','discarded')"`
	User   User   `json:"-"`
}
