package requests

import (
	"errors"
)

// OvertimeCreateRequest 申请加班
type OvertimeCreateRequest struct {
	Remark string `json:"remark" binding:"required"`
}

// Validate 验证 OvertimeCreateRequest 请求中的有效性
func (r *OvertimeCreateRequest) Validate() error {
	if len(r.Remark) == 0 {
		return errors.New("Remark cannot be empty")
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
