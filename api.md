# API 定义

## 错误响应格式

    {
        "status": 4xx/5xx,
        "message": "error message"
    }

## 用户相关

### 用户认证

POST `/user/auth`

#### JSON 参数

 - name: 用户姓名
 - password: 用户密码

#### 响应

    {
        "status": 200,
        "token": "{Token}"
    }

### 用户列表

GET `/user`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "name": "test",
                "role": "user/manager/master",
                "department": 1,
                "hours": 10
            },
            ...
        ]
    }

### 添加用户

POST `/user`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 姓名
 - password: 密码
 - department: 部门

#### 响应

    {
        "status": 201,
        "user_id": 2
    }

### 修改用户

PUT `/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 姓名
 - department: 部门
 - role: 身份
 - hours: 工时

#### 响应

    {
        "status": 200
    }

### 获取指定用户信息

GET `/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "name": "test",
            ...
        }
    }

### 删除用户

DELETE `/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 部门相关

### 部门列表

GET `/department`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "name": "X 部门",
                "manager": 1
            },
            ...
        ]
    }

### 添加部门

POST `/department`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 部门名称
 - manager: 主管 ID

#### 响应

    {
        "status": 201,
        "department_id": 2
    }

### 编辑部门

PUT `/department/{department_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 部门名称
 - manager: 主管 ID

#### 响应

    {
        "status": 200
    }

### 获取指定部门信息

GET `/department/{department_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "id": 1,
            ...
        }
    }

### 删除部门

DELETE `/department/{department_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 人脸相关

### 获取指定用户可用的人脸信息

GET `/user/{user_id}/face`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "info": "xxxxxxx",
                "status": "wait/available/discarded"
            },
            ...
        ]
    }

### 更新指定用户人脸信息

POST `/user/{user_id}/face`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - info: 人脸信息

#### 响应

    {
        "status": 201,
        "face_id": 2
    }

### 审批人脸信息

PUT `/face/{face_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - status: 人脸状态

#### 响应

    {
        "status": 200
    }

## 工时相关

### 获取工时记录

GET `/hours`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - user_id: 用户 ID，可选
 - start_at: 开始日期，可选
 - end_at: 结束日期，可选
 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "date": "2019-02-02 11:11:11",
                "hours": 12
            },
            ...
        ]
    }

## 排班相关

### 排班列表

GET `/shift`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - user_id: 用户 ID，可选
 - department_id: 部门 ID，可选
 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "type": "normal/overtime/allovertime",
                "status": "no/on/off/leave"
            },
            ...
        ]
    }

### 添加排班

POST `/user/{user_id}/shift`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - start_at: 开始时间
 - end_at: 结束时间
 - type: 排班类型

#### 响应

    {
        "status": 201,
        "shift_id": 2
    }

### 部门排班

POST `/department/{department_id}/shift`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - start_at: 开始时间
 - end_at: 结束时间
 - type: 排班类型

#### 响应

    {
        "status": 201,
        "shift_ids": [2, 3]
    }

### 更新排班状态

PUT `/shift/{shift_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - status: 排班状态

#### 响应

    {
        "status": 200
    }

### 删除排班

DELETE `/shift/{shift_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 请假相关

### 请假列表

GET `/leave`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "身体原因",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 获取指定用户请假

GET `/user/{user_id}/leave`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "身体原因",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 申请请假

POST `/user/{user_id}/leave`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - start_at: 开始时间
 - end_at: 结束时间
 - remark: 请假理由

#### 响应

    {
        "status": 201,
        "leave_id": 2
    }

### 审批请假

PUT `/leave/{leave_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - status: 状态

#### 响应

    {
        "status": 200
    }

## 加班相关

### 加班申请列表

GET `/overtime`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "任务未完成",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 获取指定用户加班

GET `/user/{user_id}/overtime`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "page": 1,
        "data": [
            {
                "id": 1,
                "user_id": 1,
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "任务未完成",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 申请加班

POST `/user/{user_id}/overtime`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - start_at: 开始时间
 - end_at: 结束时间
 - remark: 理由

#### 响应

    {
        "status": 201,
        "overtime_id": 2
    }

### 审批加班

PUT `/overtime/{overtime_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - status: 状态

#### 响应

    {
        "status": 200
    }

## 签到相关

### 获取用户当前签到状态

GET `/user/{user_id}/sign`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    no/leave = {
        "status": 204
    }

    on/off = {
        "status": 200,
        "data": {
            "id": 1,
            "shift_id": 1,
            "start_at": "2019-02-02 11:11:11",
            "end_at": "2019-02-02 11:11:11"
        }
    }

### 获取二维码

GET `/sign/qrcode`

#### 说明

二维码解析后为签到使用的 Token

#### 响应

    {
        "status": 200,
        "data": {
            "qrcode": "data:image/png;...",
            "count": 10,
            "expired_at": "2019-02-02 11:11:11"
        }
    }

### 二维码签到

POST `/user/{user_id}/sign/qrcode`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - token: 签到 Token

#### 响应

    {
        "status": 200,
        "sign_id: 3
    }

### 人脸签到

POST `/user/{user_id}/sign/face`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - face: 人脸图片的 base64 编码

#### 响应

    {
        "status": 200,
        "sign_id: 3
    }

### 签退

POST `/sign/{sign_id}/off`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }
