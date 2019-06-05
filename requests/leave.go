package requests

import (
	"errors"
	"github.com/szdx4/attsys-server/config"
)

// LeaveCreateRequest 申请请假请求
type LeaveCreateRequest struct {
	StartAt string `json:"start_at" binding:"required"`
	EndAt   string `json:"end_at" binding:"required"`
	Remark  string `json:"remark" binding:"required"`
}

// Validate 验证 LeaveCreateRequest 请求中的有效性
func (r *LeaveCreateRequest) Validate() error {
	// 将接收的 string 格式转换成 Time
	startAt, err := config.StrToTime(r.StartAt)
	if err != nil {
		return errors.New("start_at not valid")
	}
	endAt, err := config.StrToTime(r.EndAt)
	if err != nil {
		return errors.New("end_at not valid")
	}

	// 验证给出请假时间的有效性
	if startAt.After(endAt) {
		return errors.New("Time not valid")
	}
	return nil
}
