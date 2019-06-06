package requests

import (
	"errors"
	"github.com/szdx4/attsys-server/config"
)

// OvertimeCreateRequest 申请加班
type OvertimeCreateRequest struct {
	StartAt string `json:"start_at" binding:"required"`
	EndAt   string `json:"end_at" binding:"required"`
	Remark  string `json:"remark" binding:"required"`
}

// Validate 验证 OvertimeCreateRequest 请求中的有效性
func (r *OvertimeCreateRequest) Validate() error {
	// 将接收的 string 格式转换成 Time
	startAt, err := config.StrToTime(r.StartAt)
	if err != nil {
		return errors.New("start_at not valid")
	}
	endAt, err := config.StrToTime(r.EndAt)
	if err != nil {
		return errors.New("end_at not valid")
	}

	// 判断加班时间阈

	// 验证给出请假时间的有效性
	if startAt.After(endAt) {
		return errors.New("Time not valid")
	}
	return nil
}

// OvertimeUpdateRequest 审批加班
type OvertimeUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

// Validate 验证 OvertimeUpdateRequest 请求中的有效性
func (r *OvertimeUpdateRequest) Validate() error {
	// 验证状态的有效性
	if r.Status != "wait" && r.Status != "pass" && r.Status != "reject" {
		return errors.New("Status not valid")
	}
	return nil
}
