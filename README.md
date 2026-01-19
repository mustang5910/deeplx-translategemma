# DeepLX-TranslateGemma

这是一个基于 Go (Go-Zero) 开发的翻译服务，它提供了与 **DeepLX** 兼容的 API 接口，底层使用 **Ollama** 运行 Google 的 **TranslateGemma** 模型进行高质量翻译。

## ✨ 特性

- **DeepLX 兼容**: API 请求参数与返回结构完全兼容 DeepLX，可直接替换原有服务使用。
- **本地大模型**: 利用 Ollama 本地运行 `translategemma` 模型，隐私安全，无需联网调用外部收费 API。
- **高性能**: 基于 Go-Zero 框架，并发性能优秀。

## 🛠️ 前置要求

在运行本项目之前，请确保你已经安装并配置好了以下环境：

1. **Ollama**: 需要安装并运行 Ollama 服务。
2. **TranslateGemma 模型**: 在 Ollama 中下载该模型。
   ```bash
   ollama pull translategemma
   ```
3. **Go 环境** (如果是源码编译运行): Go 1.20+

## 🚀 快速开始

### 方式一：源码运行

1. **克隆项目**
   ```bash
   git clone https://github.com/mustang5910/deeplx-translategemma
   cd deeplx-translategemma
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置**
   检查 `etc/deeplx-api.yaml` 文件，确保 Ollama 地址正确。
   ```yaml
   Name: deeplx-api
   Host: 0.0.0.0
   Port: 8888
   Model: translategemma
   OllamaUrl: http://127.0.0.1:11434
   ```

4. **运行**
   ```bash
   go run deeplx.go
   ```

### 方式二：Docker (可选)

*(如果有 Dockerfile，可以取消注释以下步骤)*
<!--
1. 构建镜像
   docker build -t deeplx-gemma .
2. 运行容器
   docker run -d -p 8888:8888 --network host deeplx-gemma
-->

## ⚙️ 配置说明

配置文件位于 `etc/deeplx-api.yaml`：

| 配置项      | 说明                   | 默认值                   |
| ----------- | ---------------------- | ------------------------ |
| `Port`      | 服务监听端口           | `8888`                   |
| `Model`     | 使用的 Ollama 模型名称 | `translategemma`         |
| `OllamaUrl` | Ollama 服务地址        | `http://127.0.0.1:11434` |

## 🔌 API 接口

### 翻译接口

- **URL**: `/translate` (假设路由通常是这个，或者根路径，需确认 handler)
- **Method**: `POST`
- **Content-Type**: `application/json`

**请求示例:**

```bash
curl -X POST http://localhost:8888/translate \
-d '{
    "text": "Hello world",
    "source_lang": "EN",
    "target_lang": "ZH"
}'
```

**响应示例:**

```json
{
    "code": 200,
    "data": "你好世界",
    "source_lang": "EN",
    "target_lang": "ZH"
}
```

## 📝 开发

项目使用 `goctl` 生成。

- 修改 `.api` 文件后重新生成代码（如果存在）：
  ```bash
  goctl api go -api deeplx.api -dir .
  ```
