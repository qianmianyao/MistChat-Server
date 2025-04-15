# 更新日志 (CHANGELOG)

本文件记录项目的所有重要更改。

## [未发布]

### 2025-4-15

#### 代码优化
- 优化了 `internal/websocket/message_type` 目录下所有文件的注释
  - 移除了冗余和过度详细的注释描述
  - 统一了注释风格，使其符合项目编码规则
  - 确保注释简洁明了，专注于"做什么"而非"怎么做"
  - 影响文件:
    - `message_interface.go`: 优化了 Message 接口和 BaseMessage 结构的注释
    - `system.go`: 简化了 SystemMessage 相关方法的注释
    - `text.go`: 简化了 TextMessage 相关方法的注释
    - `message_parser.go`: 优化了消息解析函数的注释 