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
	// 验证用户是否存在
	user := &models.User{}
	database.Connector.Where("name = ?", r.Name).First(&user)
	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	// 验证密码是否匹配
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
	// 验证用户名长度
	if len(r.Name) < config.App.MinUserLength {
		return errors.New("User name must longer than " + strconv.Itoa(config.App.MinUserLength))
	}

	// 验证用户名存在性
	user := models.User{}
	database.Connector.Where("name = ?", r.Name).First(&user)
	if user.ID > 0 {
		return errors.New("User name exists")
	}

	// 验证密码长度
	if len(r.Password) < config.App.MinPwdLength {
		return errors.New("Password must longer than " + strconv.Itoa(config.App.MinPwdLength))
	}

	// 验证部门是否存在
	department := models.Department{}
	database.Connector.First(&department, r.Department)
	if department.ID == 0 {
		return errors.New("Department not exists")
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

	// 处理部门主管冲突
	if r.Role == "manager" {
		user := models.User{}
		database.Connector.Where("id <> ? AND department_id = ? AND role = 'manager'", userID, r.Department).First(&user)
		if user.ID > 0 {
			return 0, errors.New("Department can only have one manager")
		}
	}

	return userID, nil
}

// UserPasswordRequest 修改密码请求
type UserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Validate 验证修改密码请求的合法性
func (r *UserPasswordRequest) Validate(role string, authID, userID int) (string, error) {
	// 用户必须输入原密码
	if role == "user" && r.OldPassword == "" {
		return "", errors.New("Old password missing")
	}

	// 用户只能修改自己的密码
	if role == "user" && authID != userID {
		return "", errors.New("You can only modify your own password")
	}

	// 验证用户是否存在
	user := models.User{}
	database.Connector.First(&user, userID)
	if user.ID == 0 {
		return "", errors.New("User not found")
	}

	// 验证用户的原密码
	if role == "user" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPassword))
		if err != nil {
			return "", errors.New("Password invalid")
		}
	}

	// 验证新密码的长度
	if len(r.NewPassword) < config.App.MinPwdLength {
		return "", errors.New("Password must longer than " + strconv.Itoa(config.App.MinPwdLength))
	}

	// 生成新密码的 bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), 10)
	if err != nil {
		return "", errors.New("Password hash generate error")
	}

	return string(hash), nil
}
