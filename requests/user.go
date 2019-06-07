package requests

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	//"math/bits"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
	"golang.org/x/crypto/bcrypt"
)

// UserAuthRequest 用户认证请求
type UserAuthRequest struct {
	Name     string `binding:"required"`
	Password string `binding:"required"`
}

// Validate 验证 UserAuthRequest 请求中用户信息的有效性
func (r *UserAuthRequest) Validate() (*models.User, error) {
	user := &models.User{}
	database.Connector.Where("name = ?", r.Name).First(&user)
	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
	if err != nil {
		return nil, errors.New("Password invalid")
	}

	return user, nil
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Name       string `binding:"required"`
	Password   string `binding:"required"`
	Department uint   `binding:"required"`
}

// Validate 验证创建用户请求的合法性
func (r *UserCreateRequest) Validate() error {
	if len(r.Name) < config.App.MinUserLength {
		return errors.New("User name must longer than " + strconv.Itoa(config.App.MinUserLength))
	}

	user := models.User{}
	database.Connector.Where("name = ?", r.Name).First(&user)
	if user.ID > 0 {
		return errors.New("User name exists")
	}

	if len(r.Password) < config.App.MinPwdLength {
		return errors.New("Password must longer than " + strconv.Itoa(config.App.MinPwdLength))
	}

	return nil
}

// UserUpdateRequest 更改用户信息请求
type UserUpdateRequest struct {
	Name       string `binding:"required"`
	Department uint   `binding:"required"`
	Role       string `binding:"required"`
	Hours      uint   `binding:"required"`
}

// Validate 验证 UserUpdateRequest 请求中信息的有效性
func (r *UserUpdateRequest) Validate(c *gin.Context) (int, error) {
	// 验证名字的有效性
	if len(r.Name) < config.App.MinUserLength {
		return 0, errors.New("User name must longer than " + strconv.Itoa(config.App.MinUserLength))
	}

	// 验证名字的冲突性
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("User ID invalid")
	}
	user := models.User{}
	database.Connector.Where("name = ? AND id <> ?", r.Name, userID).First(&user)
	if user.ID > 0 {
		return 0, errors.New("User name exists")
	}

	// 验证部门的存在与否
	department := models.Department{}
	database.Connector.First(&department, r.Department)
	if department.ID == 0 {
		return 0, errors.New("Department not exists")
	}

	// 验证角色的有效性
	if r.Role != "user" && r.Role != "manager" && r.Role != "master" {
		return 0, errors.New("User role not valid")
	}
	if r.Role == "manager" {
		user := models.User{}
		database.Connector.Where("user_id <> ? AND department_id = ? AND role = 'manager'", userID, r.Department).First(&user)
		if user.ID > 0 {
			return 0, errors.New("Department can only have one manager")
		}
	}

	// 验证 Hours 的有效性
	if r.Hours < 0 {
		return 0, errors.New("User hours not valid")
	}

	return userID, nil
}
