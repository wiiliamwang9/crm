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
├── backend/                    # Go后端服务（精简化单文件架构）
│   ├── main.go                # 应用启动入口
│   ├── config.go              # 配置管理（数据库连接、服务器配置）
│   ├── config.yml             # YAML配置文件
│   ├── models.go              # 数据模型定义（所有实体模型）
│   ├── dto.go                 # 数据传输对象（请求/响应结构体）
│   ├── crm.go                 # 核心业务逻辑（所有业务函数）
│   ├── routes.go              # 路由配置（所有API路由定义）
│   ├── util.go                # 工具函数（字符串处理、类型转换等）
│   ├── go.mod                 # Go模块配置
│   ├── go.sum                 # 依赖版本锁定
│   └── init.sql               # 数据库初始化脚本
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

数据库连接信息在config.yml中配置：
- Host: db.lamdar.cn
- Port: 9524
- Database: crm
- Username: postgres
- Password: tpg1688

### 3. 启动后端服务
```bash
cd backend
go run main.go
```

服务将在 `http://localhost:8081` 启动（可在config.yml中配置端口）

### 4. 访问前端

打开浏览器访问 `http://localhost:8081` 即可使用系统

## API接口文档

### 客户管理 API

- `GET /api/v1/customers` - 获取客户列表（支持分页、搜索）
- `GET /api/v1/customers/:id` - 获取单个客户详细信息
- `POST /api/v1/customers` - 创建新客户
- `PUT /api/v1/customers/:id` - 更新客户信息
- `DELETE /api/v1/customers/:id` - 删除客户
- `GET /api/v1/customers/search` - 客户搜索（支持关键词和系统标签）

### 待办事项 API

- `GET /api/v1/todos` - 获取待办事项列表（支持客户筛选和分页）
- `POST /api/v1/todos` - 创建待办事项
- `PUT /api/v1/todos/:id` - 更新待办事项

### 跟进记录 API

- `GET /api/v1/activities` - 获取跟进记录列表（支持客户筛选和分页）
- `POST /api/v1/activities` - 创建跟进记录

### 用户管理 API

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情（智能判断员工/客户身份）

### 仪表板 API

- `POST /api/v1/dashboard/search` - 仪表板搜索（支持多维度筛选）

### 标签管理 API

- `GET /api/v1/tags` - 获取标签列表
- `POST /api/v1/tags` - 创建标签

### 提醒系统 API
- `GET /api/v1/reminders` - 获取提醒列表
- `POST /api/v1/reminders` - 创建提醒

### 系统 API

- `GET /health` - 健康检查

## 后端架构详解

### 精简化单文件架构

项目采用精简的单文件架构设计，将原来的多层分布式代码重构为功能明确的单文件模块：

#### 文件职责划分

1. **main.go** - 应用启动流程
    - 配置文件加载
    - 数据库连接初始化
    - GORM自动迁移
    - Gin路由设置和服务启动

2. **config.go** - 配置管理
    - YAML配置文件解析
    - 数据库连接配置
    - 服务器运行模式设置
    - 全局配置变量管理

3. **models.go** - 数据模型层
    - 所有实体模型定义（Customer、Todo、User、Activity等）
    - 自定义JSONB类型实现
    - 枚举类型定义
    - 模型方法（IsOverdue、GetTimeAgo等）

4. **dto.go** - 数据传输层
    - 请求/响应结构体定义
    - 模型转换函数（ToModel/FromModel）
    - 数据验证规则
    - 类型转换辅助函数

5. **crm.go** - 业务逻辑层
    - 所有业务函数实现
    - 客户管理：增删改查、搜索、数据聚合
    - 待办管理：创建、更新、状态追踪
    - 跟进记录：活动记录、时间计算
    - 用户管理：角色识别、详情查询
    - 仪表板：多维度数据筛选和聚合

6. **routes.go** - 路由配置
    - 统一的路由注册
    - CORS中间件配置
    - API端点定义和参数绑定
    - 响应格式标准化

7. **util.go** - 工具函数
    - 字符串处理（分割、去空格、连接）
    - 类型转换（字符串转int64、数组解析）
    - JSON序列化/反序列化封装

### 架构优势

- **简化维护**：代码集中在少数文件中，便于理解和修改
- **快速部署**：编译后仅需单个可执行文件和配置文件
- **清晰职责**：每个文件功能明确，避免代码重复
- **配置外化**：通过YAML文件管理环境配置
- **自动迁移**：启动时自动同步数据库结构

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
id          int4      跟进记录ID
customer_id int4      客户ID
user_id     int4      用户ID
kind        int4      类型
data        jsonb     沟通数据
created_at  timestamp 创建时间
remark      text      备注
```

## 技术特点

### 后端架构优势

- **精简架构**：采用单文件架构模式，降低维护复杂度
- **配置驱动**：YAML配置文件统一管理数据库和服务器配置
- **GORM自动迁移**：启动时自动创建和更新数据库表结构
- **统一响应格式**：标准化的JSON响应结构
- **类型安全**：完整的枚举定义和数据验证

### 业务逻辑优势

- **智能搜索**：支持多维度客户搜索和仪表板数据筛选
- **关联查询**：GORM预加载优化复杂关联查询性能
- **数据聚合**：智能的客户待办聚合和时间维度分析
- **用户角色识别**：动态判断用户身份（员工/客户）

### 数据模型特点

- **JSONB支持**：灵活的非结构化数据存储
- **PostgreSQL数组**：高效的多值字段支持
- **软删除机制**：BaseModel提供统一的删除标记
- **时间维度**：完整的创建、更新、删除时间追踪

## 开发规范

- **单一职责**：每个文件承担明确的功能职责
- **RESTful设计**：标准化的API端点设计
- **配置化管理**：外部化配置，便于环境切换
- **类型驱动**：强类型定义，减少运行时错误
- **响应统一**：标准化的JSON响应格式

## 贡献指南

1. Fork 项目到你的GitHub仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证，详情请查看 [LICENSE](LICENSE) 文件。