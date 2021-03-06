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
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
	r.Use(cors.New(corsConfig))

	gin.SetMode(config.App.RunMode)

	r.GET("/", controllers.Home)

	// User
	// 用户认证
	r.POST("/user/auth", controllers.UserAuth)
	// 批量添加用户
	r.POST("/user/batch", middleware.Token, middleware.Master, controllers.UserBatch)
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
	// 获取指定人脸信息
	r.GET("/face/user/:id", middleware.Token, controllers.FaceUserShow)
	// 更新人脸信息
	r.POST("/face/user/:id", middleware.Token, controllers.FaceCreate)
	// 审批人脸信息
	r.PUT("/face/:id", middleware.Token, middleware.Master, controllers.FaceUpdate)
	// 获取人脸列表
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
	// 全单位排班
	r.POST("/shift/all", middleware.Token, middleware.Master, controllers.ShiftAll)
	// 删除排班
	r.DELETE("/shift/:id", middleware.Token, middleware.Manager, controllers.ShiftDelete)
	// 修改排班
	r.PUT("/shift/:id", middleware.Token, middleware.Manager, controllers.ShiftUpdate)

	// Leave
	// 申请请假
	r.POST("/leave/user/:id", middleware.Token, controllers.LeaveCreate)
	// 获取指定用户请假
	r.GET("/leave/user/:id", middleware.Token, controllers.LeaveShow)
	// 请假列表
	r.GET("/leave", middleware.Token, middleware.Manager, controllers.LeaveList)
	// 审批请假
	r.PUT("/leave/:id", middleware.Token, middleware.Manager, controllers.LeaveUpdate)
	// 销假
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
	// 获取二维码
	r.GET("/sign/qrcode", middleware.Token, middleware.Manager, controllers.SignGetQrcode)
	// 二维码签到
	r.POST("/sign/qrcode/:id", middleware.Token, controllers.SignWithQrcode)
	// 人脸签到
	r.POST("/sign/face", middleware.Token, middleware.Manager, controllers.SignWithFace)
	// 签退
	r.POST("/sign/off/:id", middleware.Token, controllers.SignOff)
	// 获取用户当前签到状态
	r.GET("/sign/user/:id", middleware.Token, controllers.SignStatus)

	// Message
	// 获取指定信息
	r.GET("/message/:id", middleware.Token, controllers.MessageShow)
	// 获取信息列表
	r.GET("/message", middleware.Token, controllers.MessageList)

	// Status
	// 获取用户相关数据
	r.GET("/status/user", middleware.Token, middleware.Manager, controllers.StatusUser)
	// 获取签到相关数据
	r.GET("/status/sign", middleware.Token, middleware.Manager, controllers.StatusSign)
	// 获取用户工作时间和加班时间
	r.GET("/status/hours/:user_id", middleware.Token, controllers.StatusHour)

	return r
}
