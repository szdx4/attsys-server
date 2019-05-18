package requests

import (
	"errors"

	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
	"golang.org/x/crypto/bcrypt"
)

// UserAuthRequest 用户认证请求
type UserAuthRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Validate 验证 UserAuthRequest 请求中用户信息的有效性
func (r *UserAuthRequest) Validate() (*models.User, error) {
	user := &models.User{
		Name: r.Username,
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
