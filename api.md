# API 定义

## 错误响应格式

    {
        "status": 4xx/5xx,
        "message": "error message"
    }

## 用户相关

### 用户认证 （ok）

POST `/user/auth`

#### JSON 参数

 - name: 用户姓名
 - password: 用户密码

#### 响应

    {
        "status": 200,
        "token": "{Token}"
    }

### 用户列表 （ok）

GET `/user`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "name": "test",
                "role": "user/manager/master",
                "department": {
                    "id": 1,
                    "name": "test"
                },
                "hours": 10
            },
            ...
        ]
    }

### 添加用户 （ok）

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

### 批量添加用户 （ok）

POST `/user/batch`

#### 说明

批量添加用户的csv文档格式样例见 user_batch.csv 文件

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - batch: 批量用户信息

#### 响应

    [
        {
            "status": 201,
            "user_id": 2
        },
        {
            "status": 400,
            "message": "why failed"
        },
        ...
    ]
    
### 修改用户 （ok）

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

### 获取指定用户信息 （ok）

GET `/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
             "id": 1,
             "name": "test",
             "role": "user/manager/master",
             "department": {
                 "id": 1,
                 "name": "test"
             },
             "hours": 10
        }
    }

### 删除用户 （ok）

DELETE `/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

### 修改密码

PUT `/user/{user_id}/password`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - old_password: 原密码，manager 和 master 权限可选
 - new_password: 新密码

#### 响应

    {
        "status": 200
    }

## 部门相关

### 部门列表 （ok）

GET `/department`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "name": "X 部门"
            },
            ...
        ]
    }

### 添加部门 （ok）

POST `/department`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 部门名称

#### 响应

    {
        "status": 201,
        "department_id": 2
    }

### 编辑部门 （ok）

PUT `/department/{department_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - name: 部门名称

#### 响应

    {
        "status": 200
    }

### 获取指定部门信息 （ok）

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

### 删除部门 （ok）

DELETE `/department/{department_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 人脸相关

### 获取指定用户可用的人脸信息 ok

GET `/face/user/{user_id}`

#### 说明

人脸格式样例见 `face.base64.txt` 文件。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "id": 1,
            "user": {
                "id": 1,
                "name": "test name",
                ...
            },
            "info": "xxxxxxx",
            "status": "wait/available/discarded"
        }
    }

### 更新指定用户人脸信息 ok

POST `/face/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - info: 人脸信息

#### 响应

    {
        "status": 201,
        "face_id": 2
    }

### 获取人脸列表 ok

GET `/face`

### Header

 - Authorization: Bearer {Token}

### URL Query 参数

 - user_id: 获取指定用户人脸信息，可选
 - status: 获取制定状态人脸信息，可选 
 - page: 分页

### 响应

    {
        "status": 200,
        "total": 100
        "current_page": 1
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "name",
                    ...
                },
                "info": "xxxxxxx",
                "status": "wait/available/discarded"
            },
            ...
        ]
    }

### 审批人脸信息 ok

PUT `/face/{face_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 工时相关

### 获取工时记录 （ok）

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
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "test",
                    ...
                },
                "date": "2019-02-02 11:11:11",
                "hours": 12
            },
            ...
        ]
    }

## 排班相关

### 排班列表 （ok）

GET `/shift`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - user_id: 用户 ID，可选
 - department_id: 部门 ID，可选
 - start_at: 开始日期，可选
 - end_at: 结束日期，可选
 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "test",
                    ...
                },
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "type": "normal/allovertime",
                "status": "no/on/off/leave"
            },
            ...
        ]
    }

### 添加排班 （ok）

POST `/shift/user/{user_id}`

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

### 部门排班 （ok）

POST `/shift/department/{department_id}`

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

### 全单位排班 （ok）

POST '/shift/all'

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
    
### 删除排班 （ok）

DELETE `/shift/{shift_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

### 修改排班

PUT `/shift/{shift_id}`

#### 说明

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - start_at: 开始时间
 - end_at: 结束时间
 - effect: temp/all  temp 代表只修改这一次排班，即临时修改;all 代表修改这一次和之后所有的排班，即永久修改

#### 响应

    {
        "status": 200
    }

## 请假相关

### 请假列表 （ok）

GET `/leave`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "name",
                    ...
                },
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "身体原因",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 获取指定用户请假 （ok）

GET `/leave/user/{user_id}`

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
                "user": {
                    "id": 1,
                    "name": "name",
                    ...
                },
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "身体原因",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 申请请假 （ok）

POST `/leave/user/{user_id}`

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

### 审批请假 （ok）

PUT `/leave/{leave_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - status: 状态

#### 响应

    {
        "status": 200
    }

### 销假

DELETE `/leave/{leave_id}`

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200
    }

## 加班相关

### 加班申请列表 （ok）

GET `/overtime`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "test",
                    ...
                },
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "任务未完成",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 获取指定用户加班 （ok）

GET `/overtime/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### URL Query 参数

 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "test",
                    ...
                },
                "start_at": "2019-02-02 11:11:11",
                "end_at": "2019-02-02 11:11:11",
                "remark": "任务未完成",
                "status": "wait/pass/reject"
            },
            ...
        ]
    }

### 申请加班 （ok）

POST `/overtime/user/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - remark: 理由

#### 响应

    {
        "status": 201,
        "overtime_id": 2
    }

### 审批加班 （ok）

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

GET `/sign/user/{user_id}`

#### 说明

如果用户已经签到，则返回对应签到条目的 ID，即 `sign_id`;如果没有签到，则返回 `204 No Content`。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "sign_id": 2,
        "shift": {
            "id": 1,
            "start_at": "2019-xxxx",
            "end_at": "2019-xxxx",
            ...
        }
    }

### 获取二维码 ok

GET `/sign/qrcode`

#### 说明

二维码解析后为签到使用的 Token

#### 响应

    {
        "status": 200,
        "data": {
            "qrcode": "data:image/png;...",
            "expired_at": "2019-02-02 11:11:11"
        }
    }

### 二维码签到 ok

POST `/sign/qrcode/{user_id}`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - token: 签到 Token

#### 响应

    {
        "status": 200,
        "sign_id: 3,
        "user": {
            "id": 1,
            "name": "xxx",
            ...
        }
    }

### 人脸签到 ok

POST `/sign/face`

#### Header

 - Authorization: Bearer {Token}

#### JSON 参数

 - face: 人脸图片的 base64 编码

#### 响应

    {
        "status": 200,
        "sign_id: 3,
        "user": {
            "id": 1,
            "name": "xxx",
            ...
        }
    }

### 签退 ok

POST `/sign/off/{sign_id}`

#### 说明

响应中的 `overtime` 为是否可以申请加班。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "overtime": true
    }

## 信息相关

### 获取信息列表

GET `/message`

#### Header

 - Authorization: Bearer {Token}

#### URL Query

 - from_user_id: 获取指定发送用户的信息，可选
 - to_user_id: 获取指定接收用户的信息，可选
 - status: 获取指定状态的信息，可选
 - page: 分页

#### 响应

    {
        "status": 200,
        "total": 100,
        "current_page": 1,
        "per_page": 20,
        "data": [
            {
                "id": 1,
                "from": {
                    "id": 1,
                    "name": "test",
                    ...
                },
                "to": {
                    "id": 2,
                    "name": "test2",
                    ...
                },
                "title": "test title",
                "status": "unread/read"
            },
            ...
        ]
    }

### 获取指定信息

GET `/message/{message_id}`

#### 说明

调用此接口后 message 会自动变为 read 状态。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "id": 1,
            "from": {
                "id": 1,
                "name": "test",
                ...
            },
            "to": {
                "id": 2,
                "name": "test2",
                ...
            },
            "title": "test title",
            "content": "test content",
            "status": "unread/read"
        }
    }

## 系统状态相关

### 获取用户相关数据

GET `/status/user`

#### 说明

部门主管和经理可以获取到用户总数。

注：部门主管返回的数据仅包含本部门用户。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "users": 100
        }
    }

### 获取签到相关数据

GET `/status/sign`

#### 说明

部门主管和经理可以获取到当天已签到人数、迟到人数、请假人数。

注：部门主管返回的数据仅包含本部门用户。

#### Header

 - Authorization: Bearer {Token}

#### 响应

    {
        "status": 200,
        "data": {
            "signed": 10,
            "lated": 3,
            "leaved": 2
        }
    }

### 获取用户的工作时间和加班时间

GET '/status/hours/{user_id}'

#### Header

- Authorization: Bearer {Token}
 
#### URL Query
 
- start_at: 获取开始时间，必选
- end_at: 获取结束时间，必选
 
 #### 响应
 
     {
         "status": 200,
         "shift_hour": 300,
         "overtime_hour": 100
     }