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

## 项目结构说明
```text
project/
│
├── api/                 # API定义、接口文档和API版本控制
│   └── v1/              # REST API接口定义
│
├── cmd/                 # 命令行应用程序入口点
│   └── server/          # 服务器主入口
│       └── main.go      # 应用程序主函数
│
├── config/              # 配置文件目录
│   ├── dev.yaml         # 开发环境配置
│   ├── test.yaml        # 测试环境配置
│   └── prod.yaml        # 生产环境配置
│
├── docs/                # Swagger规范文件
│
├── internal/            # 内部应用代码
│   ├── handler/         # HTTP请求处理器
│   │   ├── auth/        # 认证相关处理器
│   │   ├── user/        # 用户相关处理器
│   │   └── common.go    # 通用处理器函数
│   │
│   ├── models/          # 数据模型
│   │   ├── dto/         # 数据传输对象
│   │   └── entity/      # 数据库实体
│   │
│   └── services/        # 业务逻辑层
│       ├── auth/        # 认证服务
│       ├── user/        # 用户服务
│       └── common.go    # 通用服务函数
│
├── pkg/                 # 可被外部应用程序导入的库代码
│   ├── database/        # 数据库工具
│   ├── logger/          # 日志工具
│   ├── middleware/      # 中间件
│   ├── utils/           # 实用工具
│   └── validator/       # 验证器
│
└── test/                # 测试相关文件
    ├── integration/     # 集成测试
    ├── mock/            # 模拟数据
    └── unit/            # 单元测试
```