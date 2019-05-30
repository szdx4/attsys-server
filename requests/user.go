package requests

import (
	"errors"

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
	Department int    `binding:"required"`
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
