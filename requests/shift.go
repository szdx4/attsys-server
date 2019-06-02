package requests

import "errors"

// ShiftCreateRequest 添加排班
type ShiftCreateRequest struct {
	StartAt string
	EndAt   string
	Type    string
}

// Validate 验证 ShiftCreateRequest 请求中的有效性
func (r *ShiftCreateRequest) Validate() error {
	// 验证开始时间的有效性
	//if len(r.StartAt) != 9 {
	//	return errors.New("Start time not valid")
	//}
	//// 验证结束时间的有效性
	//if len(r.EndAt) != 9 {
	//	return errors.New("End time not valid")
	//}
	// 验证类型的有效性
	if r.Type != "normal" && r.Type != "overtime" && r.Type != "allovertime" {
		return errors.New("Type not valid")
	}
	return nil
}
