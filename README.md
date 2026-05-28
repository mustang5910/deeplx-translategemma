# DeepLX-TranslateGemma

兼容 [DeepLX](https://github.com/OwO-Network/DeepLX) API 的翻译服务，使用任意 OpenAI 兼容的大语言模型（如 Ollama、vLLM、OpenAI 等）作为后端翻译引擎。

## 特性

- **DeepLX 兼容** — 提供 `/translate` 端点，请求/响应格式与 DeepLX 一致，可无缝替换现有客户端（如沉浸式翻译插件）。
- **多模型支持** — 兼容任何 OpenAI API 格式的后端（Gemma、Qwen、Llama 等），只需更改配置中的 Base URL。
- **200+ 语言** — 内置完整的语言代码映射表，支持 ISO 639-1 及带区域的 BCP 47 标签。
- **并发控制** — 基于信号量的并发限流，防止后端过载。
- **智能重试** — 指数退避 + 随机抖动，自动重试 429/5xx 错误。
- **可定制提示词** — 通过模板变量 `{{.Text}}` 和 `{{.TargetLang}}` 自定义翻译 Prompt。
- **Docker 支持** — 两阶段构建，最小化运行时镜像，内置健康检查。

## 快速开始

### 前置条件

- 一个运行中的 OpenAI 兼容 API 服务（Ollama、vLLM 或任意兼容服务）
- Go 1.25+（仅源码构建需要）

### Docker 启动

```bash
# 1. 准备配置文件
cp etc/deeplx-api.yaml.example etc/deeplx-api.yaml

# 2. 编辑配置，修改 OpenAIBaseURL 和 Model
vim etc/deeplx-api.yaml

# 3. 启动
docker compose up -d
```

### 源码构建

```bash
cp etc/deeplx-api.yaml.example etc/deeplx-api.yaml
# 编辑配置文件...

go build -o deeplx-translategemma .
./deeplx-translategemma -f etc/deeplx-api.yaml
```

服务默认监听 `0.0.0.0:8888`。

## API

### 翻译

```http
POST /translate
Content-Type: application/json

{
    "text": "Hello, world!",
    "source_lang": "EN",
    "target_lang": "ZH"
}
```

响应：

```json
{
    "code": 200,
    "data": "你好，世界！",
    "source_lang": "English",
    "target_lang": "Chinese"
}
```

### 健康检查

```http
GET /health
```

返回 `"ok"`。

## 配置说明

配置文件位于 `etc/deeplx-api.yaml`。完整配置项：

```yaml
Name: deeplx-api
Host: 0.0.0.0
Port: 8888
Model: translategemma           # 后端模型名称
OpenAIBaseURL: http://127.0.0.1:11434/v1   # OpenAI 兼容的 API 端点
OpenAIKey: "sk-xxx"             # API Key（Ollama 等本地服务可填任意值）
MaxConcurrent: 10               # 最大并发请求数（≤0 默认 10）
MaxRetries: 0                   # API 调用失败重试次数（0 不重试）
# Prompt: "..."                 # 自定义翻译提示词模板，不填使用内置默认值
```

### 自定义 Prompt 模板

支持 Go template 语法，可用变量：

| 变量              | 含义                                |
| ----------------- | ----------------------------------- |
| `{{.Text}}`       | 待翻译文本                          |
| `{{.TargetLang}}` | 目标语言名称（如 English、Chinese） |
| `{{.SourceLang}}` | 源语言名称（如 English、Chinese）   |
| `{{.TargetCode}}` | 目标语言代码（如 en、zh）           |
| `{{.SourceCode}}` | 源语言代码（如 en、zh）             |

示例：

```yaml
Prompt: "请将以下{{.SourceLang}}文本翻译为{{.TargetLang}}，只输出译文：\n\n{{.Text}}"
```

## 语言列表

支持的语言代码从 `internal/logic/languages.go` 的 `LanguageMap` 中加载。语言代码大小写不敏感，支持简写（如 `EN`、`zh`）和完整 BCP 47 标签（如 `EN-US`、`zh-Hans`）。

部分常用语言代码：

| 代码 | 语言       |
| ---- | ---------- |
| EN   | English    |
| ZH   | Chinese    |
| JA   | Japanese   |
| KO   | Korean     |
| FR   | French     |
| DE   | German     |
| ES   | Spanish    |
| RU   | Russian    |
| AR   | Arabic     |
| PT   | Portuguese |

## 目录结构

```
.
├── deeplx.go                      # 程序入口
├── Dockerfile                     # Docker 构建文件
├── docker-compose.yml             # Docker Compose 编排
├── docker-entrypoint.sh           # 容器入口脚本（自动拷贝示例配置）
├── go.mod / go.sum                # Go 依赖管理
├── docs/
│   └── deeplx.api                 # go-zero API 定义文件
├── etc/
│   ├── deeplx-api.yaml            # 运行配置（gitignore）
│   └── deeplx-api.yaml.example    # 示例配置
├── internal/
│   ├── config/
│   │   └── config.go              # 配置结构体定义
│   ├── handler/
│   │   ├── routes.go              # 路由注册
│   │   ├── translatehandler.go    # 翻译接口处理器
│   │   └── healthhandler.go       # 健康检查处理器
│   ├── logic/
│   │   ├── translatelogic.go      # 核心翻译逻辑（LLM 调用 + 重试 + 并发控制）
│   │   ├── healthlogic.go         # 健康检查逻辑
│   │   └── languages.go           # 200+ 语言代码映射表
│   ├── svc/
│   │   └── servicecontext.go      # 服务上下文（依赖注入 + 并发控制初始化）
│   └── types/
│       └── types.go               # 请求/响应类型定义
└── .github/workflows/
    ├── ci.yml                     # CI（lint + build + test）
    ├── docker-publish.yml         # Docker 镜像发布到 ghcr.io
    └── release.yml                # 二进制发布（linux/amd64 + linux/arm64）
```

## 工作原理

1. 接收 DeepLX 格式的翻译请求
2. 将语言代码解析为语言名称（如 `ZH` → `Chinese`）
3. 使用配置的 Prompt 模板渲染翻译指令
4. 通过 `openai-go` 库调用 OpenAI 兼容的 Chat Completions API
5. 返回 LLM 的输出作为翻译结果

## 技术栈

- **语言**: Go 1.25
- **框架**: [go-zero](https://github.com/zeromicro/go-zero) v1
- **LLM SDK**: [openai-go](https://github.com/openai/openai-go) v3
- **构建**: Docker multi-stage build + GitHub Actions

## 许可证

MIT
