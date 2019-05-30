package requests

import (
	"errors"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// CreateDepartmentRequest 新增部门请求
type DepartmentCreateRequest struct {
	Name      string `binding:"required"`
	ManagerId uint   `binding:"required"`
}

// Validate 验证 CreateDepartment 请求有效性
func (r *DepartmentCreateRequest) Validate() error {
	department := models.Department{}
	//名字冲突检测
	database.Connector.Where("name = ?", r.Name).First(&department)
	if department.ID > 0 {
		return errors.New("Department name exists")
	}
	//名字长度检测
	if len(r.Name) < 2 {
		return errors.New("Department name not valid")
	}
	//部门主管ID存在性检测
	manager := models.User{}
	database.Connector.Where("id = ?", r.ManagerId).First(&manager)
	if manager.Role != "manager" {
		return errors.New("Manager not exist")
	}
	return nil
}
