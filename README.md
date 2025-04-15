# Parchment API

## Create a Swagger API Documentation
```bash
swag init -g ./cmd/server/main.go --parseDepth 2 --parseDependency -o ./docs
```

# 生成私钥
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

# 提取公钥（给前端）
openssl rsa -pubout -in private.pem -out public.pem

# 一次性签名
用户输入密码 → bcrypt → 签名 + salt + timestamp → 发到后端 → 后端验证签名 + bcrypt 校验

# Parchment-Server 项目结构详细说明

```
parchment-server/                   # 项目根目录
├── .gitattributes                  # Git文件属性配置
├── .github/                        # GitHub相关配置目录
├── .gitignore                      # Git忽略文件配置
├── README.md                       # 项目说明文档
│
├── api/                            # API定义和路由
│   └── v1/                         # API v1版本
│       └── router.go               # API路由配置文件
│
├── cmd/                            # 应用入口点目录
│   └── server/                     # 服务器应用
│       └── main.go                 # 服务器主程序入口
│
├── config/                         # 配置文件目录
│   └── dev.yaml                    # 开发环境配置文件
│
├── docs/                           # API文档目录
│   ├── docs.go                     # Swagger自动生成的文档代码
│   ├── swagger.json                # Swagger API描述(JSON格式)
│   └── swagger.yaml                # Swagger API描述(YAML格式)
│
├── internal/                       # 内部应用代码
│   ├── handler/                    # HTTP请求处理器
│   │   ├── chat/                   # 聊天相关处理器
│   │   │   └── chat.go             # 聊天核心处理逻辑
│   │   ├── hello/                  # Hello示例处理器
│   │   │   └── hello.go            # Hello处理逻辑
│   │   └── common.go               # 通用处理器代码
│   │
│   ├── models/                     # 数据模型
│   │   ├── config/                 # 配置相关模型
│   │   ├── dot/                    # 数据传输对象
│   │   └── entity/                 # 数据库实体模型
│   │       └── chat_models.go      # 聊天相关数据模型
│   │
│   ├── services/                   # 业务逻辑服务
│   │   └── chat/                   # 聊天相关服务
│   │       ├── create.go           # 创建聊天/房间服务
│   │       ├── find.go             # 查询聊天/房间服务
│   │       ├── find_test.go        # 查询服务测试
│   │       └── update.go           # 更新聊天/房间服务
│   │
│   └── websocket/                  # WebSocket实现
│       ├── client.go               # WebSocket客户端实现
│       ├── hub.go                  # WebSocket连接中心
│       └── message_type/           # 消息类型定义
│           ├── message_interface.go # 消息接口定义
│           ├── message_parser.go   # 消息解析器
│           ├── system.go           # 系统消息定义
│           └── text.go             # 文本消息定义
│
├── pkg/                            # 可导出的公共库代码
│   ├── config/                     # 配置相关工具
│   │   └── config.go               # 配置加载与解析
│   │
│   ├── database/                   # 数据库工具
│   │   └── database.go             # 数据库连接与操作
│   │
│   ├── encryption/                 # 加密工具
│   │   ├── generate_id.go          # ID生成工具
│   │   └── generate_id_test.go     # ID生成测试
│   │
│   ├── global/                     # 全局变量与常量
│   │   └── global.go               # 全局变量定义
│   │
│   ├── logger/                     # 日志工具
│   │   └── logger.go               # 日志配置与实现
│   │
│   └── utils/                      # 实用工具函数
│       └── result.go               # HTTP响应结果工具
│
└── go.mod                          # Go模块定义文件
```

## 目录详细说明

### 1. 顶级目录

- **`.gitattributes`**: 定义Git特定的属性配置，如行尾处理、二进制文件标记等
- **`.github/`**: 包含GitHub相关配置，如工作流、issue模板等
- **`.gitignore`**: 指定Git应忽略的文件和目录
- **`README.md`**: 项目的主要说明文档
- **`go.mod`**: Go模块定义，包含项目依赖管理

### 2. api/ 目录

包含API定义和路由配置，按版本分组：

- **`v1/`**: API第一版本
    - **`router.go`**: 定义API路由和中间件，处理HTTP请求到对应处理器的映射

### 3. cmd/ 目录

包含所有可执行命令入口点：

- **`server/`**: 主服务器应用
    - **`main.go`**: 程序主入口，负责初始化组件、连接数据库、启动HTTP服务器

### 4. config/ 目录

包含应用程序的配置文件：

- **`dev.yaml`**: 开发环境配置，包含数据库连接、服务端口等配置信息

### 5. docs/ 目录

API文档相关文件：

- **`docs.go`**: 由Swagger生成的Go代码，包含API文档结构
- **`swagger.json`/`swagger.yaml`**: Swagger格式的API规范文件，用于生成API文档

### 6. internal/ 目录

应用内部代码，不对外导出：

- **`handler/`**: HTTP请求处理器
    - **`chat/`**: 聊天相关API处理
        - **`chat.go`**: 聊天功能核心处理逻辑
    - **`hello/`**: 示例API处理器
        - **`hello.go`**: 简单的示例处理器
    - **`common.go`**: 处理器公共代码

- **`models/`**: 数据模型定义
    - **`config/`**: 配置相关模型
    - **`dot/`**: 数据传输对象(DTO)，用于API请求和响应
    - **`entity/`**: 数据库实体模型
        - **`chat_models.go`**: 聊天相关的数据库模型定义

- **`services/`**: 业务逻辑层
    - **`chat/`**: 聊天服务
        - **`create.go`**: 创建聊天会话、房间的业务逻辑
        - **`find.go`**: 查询聊天会话、房间的业务逻辑
        - **`update.go`**: 更新聊天会话、房间的业务逻辑

- **`websocket/`**: WebSocket通信实现
    - **`client.go`**: WebSocket客户端连接处理
    - **`hub.go`**: WebSocket连接中心，管理所有连接
    - **`message_type/`**: 消息类型定义
        - **`message_interface.go`**: 消息接口定义
        - **`message_parser.go`**: 消息解析器实现
        - **`system.go`**: 系统消息类型定义
        - **`text.go`**: 文本消息类型定义

### 7. pkg/ 目录

可被外部导入的公共库代码：

- **`config/`**: 配置工具
    - **`config.go`**: 配置文件加载和解析实现

- **`database/`**: 数据库工具
    - **`database.go`**: 数据库连接和操作封装

- **`encryption/`**: 加密和安全工具
    - **`generate_id.go`**: 安全ID生成工具
    - **`generate_id_test.go`**: ID生成单元测试

- **`global/`**: 全局变量和常量
    - **`global.go`**: 全局变量定义

- **`logger/`**: 日志工具
    - **`logger.go`**: 日志配置和实现

- **`utils/`**: 通用工具函数
    - **`result.go`**: HTTP响应结果工具，封装统一的响应格式

## 文件功能说明

1. **`main.go`**: 应用程序的入口点，负责初始化各组件，包括配置、日志、数据库连接、WebSocket Hub，并启动HTTP服务器

2. **`router.go`**: 定义API路由并设置中间件，将HTTP请求映射到对应的处理器

3. **`chat.go` (handler)**: 处理与聊天相关的HTTP请求，如创建房间、加入房间等

4. **`client.go`**: WebSocket客户端连接管理，处理消息的接收和发送

5. **`hub.go`**: WebSocket连接中心，管理所有WebSocket连接，分发消息

6. **`message_parser.go`**: 解析不同类型的WebSocket消息

7. **`chat_models.go`**: 定义聊天相关的数据库实体模型，如房间、消息等

8. **`database.go`**: 提供数据库连接和操作的工具函数

9. **`logger.go`**: 配置和初始化日志组件，提供不同级别的日志记录能力

10. **`result.go`**: 封装统一的HTTP响应格式，简化API响应处理

这种项目结构遵循了Go项目的标准布局和最佳实践，模块化程度高，便于维护和扩展。项目明确分离了API定义、业务逻辑和数据访问层，使代码结构清晰，责任边界明确。