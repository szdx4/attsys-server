package requests

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// DepartmentCreateRequest 新增部门请求
type DepartmentCreateRequest struct {
	Name string `binding:"required"`
}

// Validate 验证 DepartmentCreateRequest 创建部门请求有效性
func (r *DepartmentCreateRequest) Validate() error {
	department := models.Department{}
	// 验证名字冲突
	database.Connector.Where("name = ?", r.Name).First(&department)
	if department.ID > 0 {
		return errors.New("Department name exists")
	}
	// 验证名字长度
	if len(r.Name) < 2 {
		return errors.New("Department name must longer than 2")
	}

	// 无误则返回空
	return nil
}

// DepartmentUpdateRequest 编辑部门请求
type DepartmentUpdateRequest struct {
	Name string `binding:"required"`
}

// Validate 验证 DepartmentUpdateRequest 编辑部门请求的有效性
func (r *DepartmentUpdateRequest) Validate(c *gin.Context) error {
	// 验证名字长度
	if len(r.Name) < 2 {
		return errors.New("Department name must longer than 2")
	}

	// 验证名字存在性
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.New("Department ID invalid")
	}
	department := models.Department{}
	database.Connector.Where("name = ? AND id <> ?", r.Name, departmentID).First(&department)
	if department.ID > 0 {
		return errors.New("Department name exists")
	}

	// 无误则返回空
	return nil
}
