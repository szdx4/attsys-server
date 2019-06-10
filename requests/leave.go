package requests

import (
	"errors"

	"github.com/szdx4/attsys-server/utils/common"
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
	startAt, err := common.ParseTime(r.StartAt)
	if err != nil {
		return errors.New("start_at not valid")
	}
	endAt, err := common.ParseTime(r.EndAt)
	if err != nil {
		return errors.New("end_at not valid")
	}

	// 验证请假时间的先后有效性
	if startAt.After(endAt) {
		return errors.New("Time not valid")
	}

	// 无误则返回空
	return nil
}

// LeaveUpdateRequest 审批请假请求
type LeaveUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

// Validate 验证 LeaveUpdateRequest 请求中的有效性
func (r *LeaveUpdateRequest) Validate() error {
	// 验证状态的有效性
	if r.Status != "pass" && r.Status != "reject" {
		return errors.New("Status not valid")
	}

	// 无误则返回空
	return nil
}
