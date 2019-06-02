package requests

import (
	"errors"
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
	user := &models.User{
		Name: r.Name,
	}
	database.Connector.First(user)
	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
	if err != nil {
		return nil, err
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
	if len(r.Name) < 2 {
		return errors.New("User name not valid")
	}

	user := models.User{}
	database.Connector.Where("name = ?", r.Name).First(&user)
	if user.ID > 0 {
		return errors.New("User name exists")
	}

	if len(r.Password) < config.App.MinPwdLength {
		return errors.New("Password not valid")
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
func (r *UserUpdateRequest) Validate() error {
	// 验证名字的有效性
	if len(r.Name) < 2 {
		return errors.New("User name not valid")
	}
	// 验证 department 的存在与否
	department := models.Department{}
	database.Connector.Where("id = ?", r.Department).First(&department)
	if department.ID == 0 {
		return errors.New("User not exists")
	}
	// 验证 role 的有效性
	if r.Role != "user" && r.Role != "manager" && r.Role != "master" {
		return errors.New("User role not valid")
	}
	// 验证 Hours 的有效性
	if r.Hours <= 0 {
		return errors.New("User hours wrong")
	}
	return nil
}
