package requests

// CreateDepartmentRequest 新增部门请求
type CreateDepartmentRequest struct {
	Name    string `binding:"required"`
	Manager uint   `binding:"required"`
}

// Validate 验证 CreateDepartment 请求有效性
func (*CreateDepartmentRequest) Validate() bool {
	return true
}
