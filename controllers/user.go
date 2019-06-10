package controllers

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
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

	// 验证提交数据的合法性
	user, err := req.Validate()
	if err != nil {
		response.Unauthorized(c, err.Error())
		c.Abort()
		return
	}

	// 生成新的 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"role":       user.Role,
		"expired_at": time.Now().Add(time.Hour * time.Duration(config.App.TokenValid)).Format(time.RFC3339),
	})
	tokenString, err := token.SignedString([]byte(config.App.EncryptKey))
	if err != nil {
		response.InternalServerError(c, "Token sign error")
		c.Abort()
		return
	}

	// 发送响应
	response.UserAuth(c, user.ID, tokenString)
}

// UserShow 获取单个用户信息
func UserShow(c *gin.Context) {
	// 获取 URL 中的用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 查询用户信息
	user := models.User{}
	database.Connector.Preload("Department").First(&user, userID)
	if user.ID < 1 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}

	// 用户只能获得自己的信息
	if role == "user" && authID != userID {
		response.Unauthorized(c, "You cannot get others information")
		c.Abort()
		return
	}

	// 部门主管只能获得本部门用户信息
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != user.DepartmentID {
			response.Unauthorized(c, "You cannot get other department information")
			c.Abort()
			return
		}
	}

	// 发送响应
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

	// 验证提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 生成密码的 bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		response.InternalServerError(c, "Password hash generate error")
		c.Abort()
		return
	}

	// 创建用户
	user := models.User{
		Name:         req.Name,
		Password:     string(hash),
		DepartmentID: uint(req.Department),
		Role:         "user",
	}
	database.Connector.Create(&user)
	if user.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	// 发送响应
	response.UserCreate(c, user.ID)
}

// UserList 获取用户列表
func UserList(c *gin.Context) {
	// 处理分页
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

	// 获取认证用户信息
	role, _ := c.Get("user_role")

	// 主管只能获得本部门用户列表
	if role == "manager" {
		managerID, _ := c.Get("user_id")
		manager := models.User{}
		database.Connector.First(&manager, managerID)
		departmentID := manager.DepartmentID

		db := database.Connector.Where("department_id = ?", departmentID).Preload("Department")
		db.Limit(perPage).Offset((page - 1) * perPage).Find(&users)
		db.Model(&models.User{}).Count(&total)
	} else {
		db := database.Connector.Preload("Department")
		db.Limit(perPage).Offset((page - 1) * perPage).Find(&users)
		db.Model(&models.User{}).Count(&total)
	}

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.UserList(c, total, page, users)
}

// UserDelete 删除用户
func UserDelete(c *gin.Context) {
	// 获取 URL 中用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	// 查询用户
	user := models.User{}
	database.Connector.First(&user, userID)
	if user.ID == 0 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}

	// 执行删除操作
	database.Connector.Delete(&user)

	// 发送响应
	response.UserDelete(c, userID)
}

// UserUpdate 修改用户
func UserUpdate(c *gin.Context) {
	var req requests.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	userID, err := req.Validate(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 查找被修改用户
	user := models.User{}
	database.Connector.First(&user, userID)
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

	// 发送响应
	response.UserUpdate(c)
}

// UserPassword 修改密码
func UserPassword(c *gin.Context) {
	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 获取 URL 中用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	var req requests.UserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	pwdHash, err := req.Validate(role.(string), authID.(int), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 执行修改密码操作
	user := models.User{}
	database.Connector.First(&user, userID)
	user.Password = pwdHash
	database.Connector.Save(&user)

	// 发送响应
	response.UserPassword(c)
}
