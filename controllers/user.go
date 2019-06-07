package controllers

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth 用户认证
func UserAuth(c *gin.Context) {
	var req requests.UserAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Bad Request")
		c.Abort()
		return
	}

	user, err := req.Validate()
	if err != nil {
		response.Unauthorized(c, "Wrong username or password")
		c.Abort()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"role":       user.Role,
		"expired_at": time.Now().UTC().Add(time.Hour * 2).Format(time.UnixDate),
	})

	tokenString, err := token.SignedString([]byte(config.App.EncryptKey))
	if err != nil {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.UserAuth(c, user.ID, tokenString)
}

// UserShow 获取单个用户信息
func UserShow(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	user := models.User{}
	database.Connector.First(&user, userID)
	if user.ID < 1 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}

	response.UserShow(c, user)
}

// UserCreate 新建用户
func UserCreate(c *gin.Context) {
	var req requests.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	user := models.User{
		Name:         req.Name,
		Password:     string(hash),
		DepartmentID: uint(req.Department),
		Role:         "user",
	}
	database.Connector.Create(&user)

	if user.ID < 1 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.UserCreate(c, user.ID)
}

// UserList 获取用户列表
func UserList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage

	users := []models.User{}
	total := 0
	database.Connector.Limit(perPage).Offset((page - 1) * perPage).Find(&users)
	database.Connector.Model(&models.User{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.UserList(c, total, page, users)
}

// UserDelete 删除用户
func UserDelete(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	user := models.User{}
	database.Connector.Where("id = ?", userID).First(&user)

	if user.ID == 0 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}

	database.Connector.Delete(&user)

	response.UserDelete(c, userID)
}

// UserUpdate 修改用户
func UserUpdate(c *gin.Context) {
	// 请求合法性检验
	var req requests.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	if err := req.Validate(c); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	user := models.User{}
	database.Connector.Where("id = ?", userID).First(&user)
	if user.ID == 0 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}
	// 修改用户的相应信息
	user.Name = req.Name
	user.DepartmentID = uint(req.Department)
	user.Role = req.Role
	user.Hours = uint(req.Hours)
	database.Connector.Save(&user)

	response.UserUpdate(c)
}
