package requests

import (
	"errors"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// CreateDepartmentRequest 新增部门请求
type DepartmentCreateRequest struct {
	Name    string `binding:"required"`
	Manager uint   `binding:"required"`
}

// Validate 验证 CreateDepartment 请求有效性
func (r *DepartmentCreateRequest) Validate() error {
	department := models.Department{}
	database.Connector.Where("name = ?", r.Name).First(&department)
	if department.ID > 0 {
		return errors.New("Department name exists")
	}

	if len(r.)
	return nil
}

// Validate 验证创建用户请求的合法性
//func (r *UserCreateRequest) Validate() error {
//	user := models.User{}
//	database.Connector.Where("name = ?", r.Name).First(&user)
//	if user.ID > 0 {
//		return errors.New("User name exists")
//	}
//
//	if len(r.Password) < config.App.MinPwdLength {
//		return errors.New("Password not valid")
//	}
//
//	department := models.Department{}
//	database.Connector.First(&department, r.Department)
//	if department.ID < 1 {
//		return errors.New("Department not found")
//	}
//
//	return nil
//}
