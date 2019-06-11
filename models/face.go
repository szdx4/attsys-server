package models

// Face 人脸模型
type Face struct {
	CommonFields
	UserID uint   `json:"-"`                                                       // 人脸所属用户
	Info   string `json:"info" gorm:"type:longtext"`                               // 人脸图片编码
	Status string `json:"status" gorm:"type:enum('wait','available','discarded')"` // 人脸审核状态
	User   User   `json:"user"`
}
