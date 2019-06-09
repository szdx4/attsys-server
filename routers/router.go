package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/controllers"
	"github.com/szdx4/attsys-server/middleware"
)

// Router 设置路由和公共中间件，返回 Gin Engine 对象
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	gin.SetMode(config.App.RunMode)

	r.GET("/", controllers.Home)

	// User
	// 用户认证
	r.POST("/user/auth", controllers.UserAuth)
	// 获取指定用户信息
	r.GET("/user/:id", middleware.Token, controllers.UserShow)
	// 添加用户
	r.POST("/user", middleware.Token, middleware.Master, controllers.UserCreate)
	// 获取用户列表
	r.GET("/user", middleware.Token, middleware.Manager, controllers.UserList)
	// 删除用户
	r.DELETE("/user/:id", middleware.Token, middleware.Master, controllers.UserDelete)
	// 修改用户
	r.PUT("/user/:id", middleware.Token, middleware.Master, controllers.UserUpdate)
	// 修改密码
	r.PUT("/user/:id/password", middleware.Token, controllers.UserPassword)

	// Department
	// 添加部门
	r.POST("/department", middleware.Token, middleware.Master, controllers.DepartmentCreate)
	// 获取指定部门信息
	r.GET("/department/:id", middleware.Token, controllers.DepartmentShow)
	// 获取部门列表
	r.GET("/department", middleware.Token, middleware.Master, controllers.DepartmentList)
	// 编辑部门
	r.PUT("/department/:id", middleware.Token, middleware.Master, controllers.DepartmentUpdate)
	// 删除部门
	r.DELETE("/department/:id", middleware.Token, middleware.Master, controllers.DepartmentDelete)

	// Face
	r.GET("/face/user/:id", middleware.Token, controllers.FaceUserShow)
	r.POST("/face/user/:id", middleware.Token, controllers.FaceCreate)
	r.PUT("/face/:id", middleware.Token, middleware.Master, controllers.FaceUpdate)
	r.GET("/face", middleware.Token, middleware.Master, controllers.FaceList)

	// Hours
	// 获取工时记录
	r.GET("/hours", middleware.Token, controllers.HoursShow)

	// Shift
	// 添加排班
	r.POST("/shift/user/:id", middleware.Token, middleware.Manager, controllers.ShiftCreate)
	// 排班列表
	r.GET("/shift", middleware.Token, controllers.ShiftList)
	// 部门排班
	r.POST("/shift/department/:department_id", middleware.Token, middleware.Manager, controllers.ShiftDepartment)
	// 删除排班
	r.DELETE("/shift/:id", middleware.Token, middleware.Manager, controllers.ShiftDelete)

	// Leave
	// 申请请假
	r.POST("/leave/user/:id", middleware.Token, controllers.LeaveCreate)
	// 获取指定用户请假
	r.GET("/leave/user/:id", middleware.Token, controllers.LeaveShow)
	// 请假列表
	r.GET("/leave", middleware.Token, middleware.Manager, controllers.LeaveList)
	// 审批请假
	r.PUT("/leave/:id", middleware.Token, middleware.Manager, controllers.LeaveUpdate)
	r.DELETE("/leave/:id", middleware.Token, controllers.LeaveDelete)

	// Overtime
	// 申请加班
	r.POST("/overtime/user/:id", middleware.Token, controllers.OvertimeCreate)
	// 获取指定用户加班
	r.GET("/overtime/user/:id", middleware.Token, controllers.OvertimeShow)
	// 加班申请列表
	r.GET("/overtime", middleware.Token, middleware.Manager, controllers.OvertimeList)
	// 审批加班
	r.PUT("/overtime/:id", middleware.Token, middleware.Manager, controllers.OvertimeUpdate)

	// Sign
	r.GET("/sign/qrcode", middleware.Token, middleware.Manager, controllers.SignGetQrcode)
	r.POST("/sign/qrcode/:id", middleware.Token, controllers.SignWithQrcode)
	r.POST("/sign/face/:id", middleware.Token, middleware.Manager, controllers.SignWithFace)
	r.POST("/sign/off/:id", middleware.Token, controllers.SignOff)
	r.GET("/sign/user/:id", middleware.Token, controllers.SignStatus)

	// Message
	// 获取指定信息
	r.GET("/message/:id", middleware.Token, controllers.MessageShow)
	// 获取信息列表
	r.GET("/message", middleware.Token, controllers.MessageList)

	return r
}
