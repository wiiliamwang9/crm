# CRM客户管理系统

这是一个功能完整的客户关系管理系统，基于Go后端和现代化前端架构，支持客户管理、待办事项、提醒功能等企业级功能。

## 技术栈

### 后端
- Go 1.23+ (工具链 1.24.6)
- Gin Web框架
- GORM ORM框架
- PostgreSQL数据库
- 支持复杂数据类型（JSON、数组等）

### 前端
- 原生HTML/CSS/JavaScript
- 模块化组件架构
- 响应式设计
- API配置化管理

## 核心功能

### 客户管理
- 客户信息的完整CRUD操作
- 高级查询和分页功能
- Excel文件批量导入客户数据
- 智能客户信息合并和更新
- 支持复杂的客户字段类型（数组、JSONB等）
- 客户分级管理（S/A/B/C/X级）
- 客户状态跟踪
- 多维度客户分类

### 待办事项管理
- 完整的待办事项生命周期管理
- 多种优先级设置（低/中/高/紧急）
- 灵活的提醒机制（微信/企业微信/短信）
- 待办状态跟踪（待处理/已完成/延期/取消）
- 操作日志记录
- 客户关联的任务管理

### 提醒系统
- 多渠道提醒支持（微信、企业微信、短信）
- 灵活的提醒频率设置（单次/每日/每周/每月）
- 提醒模板管理
- 用户个性化提醒配置
- 免打扰时间设置
- 失败重试机制

### 用户管理
- 用户权限管理
- 组织架构支持（主管-下属关系）
- 个性化配置

## 项目架构

```
crm/
├── backend/                    # Go后端服务
│   ├── main.go                # 应用入口文件
│   ├── go.mod                 # Go模块配置
│   ├── start.sh               # 启动脚本
│   ├── config/                # 配置模块
│   │   ├── database.go        # 数据库连接配置
│   │   └── server.go          # 服务器配置
│   ├── models/                # 数据模型层
│   │   ├── customer.go        # 客户数据模型
│   │   ├── todo.go            # 待办事项模型
│   │   ├── reminder.go        # 提醒系统模型
│   │   ├── user.go            # 用户模型
│   │   └── group.go           # 客户组模型
│   ├── handlers/              # HTTP处理器层
│   │   ├── customer_crud.go   # 客户CRUD操作
│   │   ├── customer_query.go  # 客户查询操作
│   │   ├── excel_handler.go   # Excel导入处理
│   │   ├── todo.go            # 待办事项处理
│   │   ├── reminder_handler.go # 提醒处理
│   │   └── user_handler.go    # 用户管理
│   ├── services/              # 业务逻辑层
│   │   ├── customer_service.go
│   │   ├── todo_service.go
│   │   ├── reminder_service.go
│   │   └── user_service.go
│   ├── routes/                # 路由配置
│   │   ├── routes.go          # 主路由
│   │   ├── customer_routes.go # 客户相关路由
│   │   ├── todo_routes.go     # 待办相关路由
│   │   ├── reminder_routes.go # 提醒相关路由
│   │   └── user_routes.go     # 用户相关路由
│   ├── middleware/            # 中间件
│   │   └── response.go        # 统一响应处理
│   ├── dto/                   # 数据传输对象
│   │   └── customer_dto.go
│   ├── migrations/            # 数据库迁移脚本
│   └── sql/                   # SQL脚本
│       ├── core_tables.sql
│       ├── todo.sql
│       └── reminder.sql
├── web/                       # 前端应用
│   ├── index.html            # 主页面
│   ├── index-new.html        # 新版本页面
│   ├── config/               # 配置文件
│   │   ├── api-config.js     # API配置
│   │   └── test-config.html  # 测试配置页面
│   ├── modules/              # 功能模块
│   │   ├── customer.js       # 客户管理模块
│   │   ├── customer-list.js  # 客户列表模块
│   │   └── excel-import.js   # Excel导入模块
│   ├── components/           # 可重用组件
│   │   └── Modal.js         # 模态框组件
│   ├── utils/               # 工具函数
│   │   ├── api.js           # API工具
│   │   └── utils.js         # 通用工具
│   ├── styles/              # 样式文件
│   │   └── main.css
│   ├── query/               # 查询相关页面
│   │   ├── query_target.html
│   │   └── style.css
│   └── target/              # 目标客户管理
│       ├── target.html
│       └── components/      # 目标管理组件
├── nginx_bak.conf            # Nginx配置备份
└── README.md                # 项目文档
```

## 安装和运行

### 1. 安装依赖
```bash
cd backend
go mod tidy
```

### 2. 配置数据库
数据库连接信息已配置为：
- Host: db.lamdar.cn
- Port: 9524
- Database: walkman
- Username: postgres
- Password: tpg1688

### 3. 启动后端服务
```bash
cd backend
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 4. 访问前端
打开浏览器访问 `http://localhost:8080` 即可使用系统

## API接口文档

### 客户管理 API
- `GET /api/v1/customers` - 获取客户列表（支持分页、搜索、筛选）
- `GET /api/v1/customers/:id` - 获取单个客户详细信息
- `POST /api/v1/customers` - 创建新客户
- `PUT /api/v1/customers/:id` - 更新客户信息
- `DELETE /api/v1/customers/:id` - 删除客户
- `GET /api/v1/customers/query` - 高级客户查询
- `POST /api/v1/upload-excel` - Excel批量导入客户数据

### 待办事项 API
- `GET /api/v1/todos` - 获取待办事项列表
- `GET /api/v1/todos/:id` - 获取待办详情
- `POST /api/v1/todos` - 创建待办事项
- `PUT /api/v1/todos/:id` - 更新待办事项
- `DELETE /api/v1/todos/:id` - 删除待办事项
- `POST /api/v1/todos/:id/complete` - 完成待办事项
- `GET /api/v1/todos/logs/:id` - 获取待办操作日志

### 提醒系统 API
- `GET /api/v1/reminders` - 获取提醒列表
- `GET /api/v1/reminders/:id` - 获取提醒详情
- `POST /api/v1/reminders` - 创建提醒
- `PUT /api/v1/reminders/:id` - 更新提醒
- `DELETE /api/v1/reminders/:id` - 删除提醒
- `GET /api/v1/reminder-templates` - 获取提醒模板
- `GET /api/v1/reminder-config` - 获取用户提醒配置

### 用户管理 API
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户信息

## Excel导入格式

系统支持从销售记录Excel中提取客户信息。Excel文件应包含以下列（按顺序）：

1. 公司名称
2. 仓库名称  
3. 销售单日期
4. 销售单号
5. 退货单号
6. 退货单日期
7. 销售员ID
8. 销售员
9. 公司手机ID
10. 公司手机号
11. **客户ID** (将作为原客户ID保存)
12. **客户** (客户名称，必填)
13. 商品编码
14. 批次号
15. 任务标记
16. 商品分类一
17. 商品分类二
18. 商品分类三
19. 商品分类四
20. 商品名称
21. 商品别名
22. 单位
23. 单位重量
24. 进价
25. 商品单价（含税）
26. 数量
27. 商品金额（含税）
28. 税率
29. 折扣金额
30. 销售金额（含税）
31. 包装费
32. 运费
33. 客户承担费用
34. 余额充值
35. 余额支付
36. 收款方式
37. 收款金额
38. 收款日期
39. 收入金额（含税）
40. 审核人
41. 审核时间
42. 审核状态
43. 审核备注
44. 单据状态
45. **客户电话** (将保存为客户联系方式)
46. 发货方式
47. 发货备注
48. 快递单号
49. **收货人** (将作为联系人姓名)
50. **收货号码** (将合并到客户联系方式)
51. **收货地址** (将解析为省市区地址信息)
52. 创建者
53. 备注

### 导入逻辑说明

- **唯一标识**：以客户电话号码为唯一标识符
- **智能更新**：如果电话号码已存在，则更新客户信息而非重复创建
- **数组字段合并**：电话号码、销售员等数组字段会智能合并
- **信息补全**：新记录会补充已有客户的空缺信息（如地址、联系人等）
- **地址解析**：收货地址自动解析为省、市、区信息
- **销售员关联**：自动提取销售员ID并关联到客户记录

## 部署说明

### 开发环境
```bash
# 启动开发服务器
cd backend
go run main.go
```

### 生产环境
```bash
# 使用启动脚本
cd backend
chmod +x start.sh
./start.sh
```

### Docker部署
项目支持容器化部署，可参考nginx_bak.conf配置反向代理。

## 主要特性说明

### 智能客户管理
- **唯一标识**：以客户电话号码为唯一标识符，避免重复数据
- **智能合并**：相同电话号码的客户信息自动合并更新
- **多维度分类**：支持客户分级（S/A/B/C/X）、状态管理、标签系统
- **地址智能解析**：自动解析省市区地址信息

### 高效任务管理
- **优先级管理**：四级优先级系统（低/中/高/紧急）
- **状态追踪**：完整的任务生命周期管理
- **智能提醒**：多渠道、多频率提醒机制
- **操作审计**：完整的操作日志记录

### 灵活提醒系统
- **多渠道支持**：微信、企业微信、短信提醒
- **个性化配置**：用户可设置免打扰时间
- **失败重试**：自动重试机制确保提醒送达
- **模板管理**：支持自定义提醒模板

## 数据库表结构

### 客户表(customers)

```
id	int4
name varchar(256) 客户名，比如“阿亮烟酒茶”
contact_name varchar(256) 联系人名字，比如“李雨亮 李老板”
gender int4 客户性别
avatar varchar(2048) 头像
photos varchar(2048)[] 门头照片
remark text 备注

source varchar(256) 客户来源
created_at timestamp 创建时间
created_by int4 创建人
updated_at timestamp 最后更新时间
updated_by int4 最后更新人

phones varchar(128)[] 手机号
wechats varchar(128)[] 微信号
douyins varchar(128)[] 客户抖音号
kwais varchar(128)[] 客户快手号
redbooks varchar(128)[] 客户小红书号
wework_openids varchar(128)[] 企业微信中的客户识别号

province varchar(256) 省
city varchar(256) 市
district varchar(256) 县
district_id int4 地址编号
street varchar(256) 街道
address varchar(2048) 完整地址
lat	float8 经纬度
lon	float8 经纬度

category varchar(256) 客户分类
flags int4 标记位
tags varchar(128)[] 标签
level int4 客户分级（0=未分级,1=S,2=A,3=B,4=C,10=X）
state int4 客户状态（0=未知,1=未开发,2=开发中,3=已开发,4=已拉黑,5=已倒闭,6=同事,7=叛徒,8=同行）
kind int4 种类（0=未知,1=个体夫妻店,2=加盟连锁店,3=工厂直营店,4=其他）
added_wechat bool 已添加微信

work_phone varchar(256)[] 工作手机号
work_wechat varchar(256)[] 工作微信
credit_sale decimal 允许赊账金额
sellers int4[] 所属销售员
last_visited timestamp 最后线下联系时间
last_called timestamp 最后线上联系时间

group_id int4[] 所属组（比如连锁店、亲戚关系）
birth_place varchar(256) 出生地（方言、口音）
birth_year int 出生年
birth_month int 出生月
birth_date int 出生日
favors jsonb[] 买货偏好（比如：[{"product": "毛尖", avgPrice: "125", avgQuantity: "10"}]）
products varchar(512)[] 主营产品
annual_turnover varchar(512) 年营业额
shipping_infos jsonb[] 收货信息（比如[{"phone": "18872675676", "address": "丹阳街道八里街214号莉姐", "city": "孝感市", "receiver": "孝南莉姐", "isDefault": true, "districtId": 420902000000, "district": "孝南区", "province": "湖北省"}]）
extra_info jsonb 附加信息


```

### 待办事项表(todos)
```
id             uint64      待办ID
customer_id    uint64      关联客户ID
creator_id     uint64      创建人ID
executor_id    uint64      执行人ID
title          string      待办标题
content        string      待办内容详情
status         enum        待办状态（pending/completed/overdue/cancelled）
planned_time   timestamp   计划执行时间
completed_time timestamp   完成时间
is_reminder    bool        是否提醒
reminder_type  enum        提醒方式（wechat/enterprise_wechat/both/sms）
priority       enum        优先级（low/medium/high/urgent）
tags           json        标签
created_at     timestamp   创建时间
updated_at     timestamp   更新时间
```

### 提醒记录表(reminders)
```
id             uint64      提醒ID
todo_id        uint64      关联待办ID
user_id        uint64      提醒用户ID
type           enum        提醒方式
title          string      提醒标题
content        string      提醒内容
status         enum        提醒状态（pending/sent/failed/cancelled）
frequency      enum        提醒频率（once/daily/weekly/monthly）
schedule_time  timestamp   计划提醒时间
sent_time      timestamp   实际发送时间
retry_count    int         重试次数
created_at     timestamp   创建时间
```

### 客户组表(groups)
```
id     int4         组ID
name   varchar(256) 组名
roles  jsonb        组成人员（比如： {"朋友": [2, 3]}）
```

### 用户表(users)
```
id         int4         用户ID
name       varchar(256) 用户名
manager_id int4         主管ID
```

### 跟进记录表(activities)
```
id          int4      活动ID
customer_id int4      客户ID
user_id     int4      用户ID
kind        int4      类型
data        jsonb     沟通数据
created_at  timestamp 创建时间
remark      text      备注
```

## 技术亮点

- **模块化架构**：清晰的分层架构，便于维护和扩展
- **数据一致性**：完善的事务管理和数据校验
- **高性能查询**：优化的数据库索引和查询语句
- **安全可靠**：参数验证、SQL注入防护
- **日志审计**：完整的操作日志记录
- **响应式设计**：适配多种设备尺寸

## 开发规范

- 遵循RESTful API设计原则
- 统一的错误处理和响应格式
- 代码注释完整，便于维护
- 数据库设计规范，字段含义明确
- 前后端分离，便于独立开发和部署

## 贡献指南

1. Fork 项目到你的GitHub仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证，详情请查看 [LICENSE](LICENSE) 文件。