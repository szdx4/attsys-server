package routers

import (
	"github.com/szdx4/attsys-server/controllers"
	"github.com/szdx4/attsys-server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
)

// Router 设置路由和公共中间件，返回 Gin Engine 对象
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.App.RunMode)

	r.GET("/", controllers.Home)

	// User
	// 用户认证
	r.POST("/user/auth", controllers.UserAuth)
	// 获取指定用户信息
	r.GET("/user/:id", middleware.Token, controllers.UserShow)
	// 添加用户
	r.POST("/user", middleware.Token, middleware.Master, controllers.UserCreate)
	//r.POST("/user", controllers.UserCreate) //测试用
	// 获取用户列表
	r.GET("/user", middleware.Token, middleware.Manager, controllers.UserList)
	// 删除用户
	r.DELETE("/user/:id", middleware.Token, middleware.Manager, controllers.UserDelete)
	// 修改用户
	r.PUT("/user/:id", middleware.Token, middleware.Master, controllers.UserUpdate)

	// Department
	// 添加部门
	r.POST("/department", middleware.Token, middleware.Master, controllers.DepartmentCreate)
	// 获取指定部门信息
	r.GET("/department/:id", middleware.Token, controllers.DepartmentShow)
	// 获取部门列表
	r.GET("/department", middleware.Token, controllers.DepartmentList)
	// 编辑部门
	r.PUT("/department/:id", middleware.Token, middleware.Master, controllers.DepartmentUpdate)
	// 删除部门
	r.DELETE("/department/:id", middleware.Token, middleware.Master, controllers.DepartmentDelete)

	// Face
	r.GET("/face/user/:id", middleware.Token, controllers.FaceUserShow)

	// Hours
	// 获取工时记录
	r.GET("/hours", middleware.Token, controllers.HoursShow)

	// Shift
	// 添加排班
	r.POST("/shift/user/:id", middleware.Token, controllers.ShiftCreate)
	// 排班列表
	r.GET("/shift", middleware.Token, controllers.ShiftList)
	// 部门排班
	r.POST("/shift/department/:department_id", middleware.Token, middleware.Manager, controllers.ShiftDepartment)
	// 更新排班状态
	r.PUT("/shift/:shift_id", middleware.Token, middleware.Manager, controllers.ShiftUpdate)
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

	// Overtime
	// 申请加班
	r.POST("/overtime/user/:id", middleware.Token, controllers.OvertimeCreate)
	// 获取指定用户加班
	r.GET("/overtime/user/:id", middleware.Token, middleware.Manager, controllers.OvertimeShow)
	//

	// Sign

	return r
}
