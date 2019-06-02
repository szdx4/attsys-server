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

// Validate 验证 CreateDepartment 创建部门请求有效性
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
	//部门主管 ID 存在性检测
	manager := models.User{}
	database.Connector.Where("id = ?", r.Manager).First(&manager)
	if manager.ID == 0 {
		return errors.New("Manager not exist")
	}
	if !(manager.Role == "manager" || manager.Role == "master") {
		return errors.New("Manager not exist")
	}

	return nil
}

// DepartmentUpdateRequest 编辑部门请求
type DepartmentUpdateRequest struct {
	Name    string `binding:"required"`
	Manager uint   `binding:"required"`
}

// Validate 验证 DepartmentUpdateRequest 编辑部门请求的有效性
func (r *DepartmentUpdateRequest) Validate() error {
	//检测名字长度
	if len(r.Name) < 2 {
		return errors.New("Department name not valid")
	}
	//部门主管 ID 存在性检测
	manager := models.User{}
	database.Connector.Where("id = ?", r.Manager).First(&manager)
	if manager.ID == 0 {
		return errors.New("Manager not exist")
	}
	if !(manager.Role == "manager" || manager.Role == "master") {
		return errors.New("Manager not exist")
	}

	return nil
}
