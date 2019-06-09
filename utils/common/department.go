package common

import (
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// DepartmentManager 获取指定部门的主管信息
func DepartmentManager(dep models.Department) *models.User {
	manager := models.User{}
	database.Connector.Where("department_id = ? AND role = 'manager'", dep.ID).First(&manager)

	if manager.ID == 0 {
		return nil
	}

	return &manager
}
