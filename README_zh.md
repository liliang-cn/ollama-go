# Ollama Go 客户端

基于官方 Python 客户端设计的 [Ollama](https://ollama.ai/) Go 客户端库。

> **注意**: 这是一个受官方 Python 客户端启发的非官方 Go 客户端库。它提供相同的功能和 API 设计模式，具有 **98%+ 的功能一致性**。

## ✨ 功能特性

- **完整的 API 支持**: 支持所有 Ollama REST API 端点
- **流式支持**: 生成和聊天的实时流式响应
- **类型安全**: 完整的 Go 类型定义和编译时检查
- **灵活配置**: 多种客户端配置选项
- **错误处理**: 全面的错误处理和 JSON 解析
- **文件上传**: 用于模型创建的 Blob 上传功能
- **高级选项**: 20+ 配置函数进行微调
- **Context 支持**: 完整的 context.Context 支持用于取消操作

## 安装

```bash
go get github.com/liliang-cn/ollama-go
```

## 使用方法

### 基本用法

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    // 生成响应
    response, err := ollama.Generate(ctx, "gemma3", "天空为什么是蓝色的？")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Response)
}
```

### 聊天

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    messages := []ollama.Message{
        {
            Role:    "user",
            Content: "天空为什么是蓝色的？",
        },
    }

    response, err := ollama.Chat(ctx, "gemma3", messages)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Message.Content)
}
```

### 流式响应

`Generate` 和 `Chat` 都支持流式响应：

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    // 流式生成
    responseChan, errorChan := ollama.GenerateStream(ctx, "gemma3", "给我讲个故事")

    for {
        select {
        case response, ok := <-responseChan:
            if !ok {
                return
            }
            fmt.Print(response.Response)
        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

### 自定义客户端

您可以创建具有特定配置的自定义客户端：

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    // 创建自定义 HTTP 客户端
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    // 创建具有自定义配置的客户端
    client, err := ollama.NewClient(
        ollama.WithHost("http://localhost:11434"),
        ollama.WithHTTPClient(httpClient),
        ollama.WithHeaders(map[string]string{
            "Custom-Header": "custom-value",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    
    req := &ollama.GenerateRequest{
        Model:  "gemma3",
        Prompt: "你好，世界！",
    }

    response, err := client.Generate(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Response)
}
```

### 嵌入向量

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    // 创建嵌入向量
    response, err := ollama.Embed(ctx, "nomic-embed-text", "敏捷的棕色狐狸")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("生成了 %d 个嵌入向量\n", len(response.Embeddings))
}
```

### 使用完整选项创建模型

```go
package main

import (
    "context"
    
    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    client, _ := ollama.NewClient()
    
    // 使用完整配置创建模型
    req := &ollama.CreateRequest{
        Model:     "my-custom-model",
        Modelfile: "FROM llama2\nSYSTEM \"你是一个有用的助手。\"",
        Files:     map[string]string{"data.txt": "训练数据"},
        Adapters:  map[string]string{"lora": "适配器数据"},
        Template:  "{{ .System }}{{ .Prompt }}",
        License:   "MIT",
        System:    "自定义系统提示",
        Parameters: &ollama.Options{
            Temperature: ollama.Float64Ptr(0.7),
        },
        Messages: []ollama.Message{
            {Role: "system", Content: "你很有帮助"},
        },
    }
    
    status, err := client.Create(ctx, req)
    if err != nil {
        panic(err)
    }
    fmt.Printf("模型已创建：%s\n", status.Status)
}
```

### 文件上传 (Blob)

```go
package main

import (
    "context"
    
    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    // 上传文件并获取其摘要
    digest, err := ollama.CreateBlob(ctx, "/path/to/file.bin")
    if err != nil {
        panic(err)
    }
    fmt.Printf("文件已上传，摘要：%s\n", digest)
}
```

### 进度流式传输

对于拉取模型等操作，您可以流式传输进度更新：

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    progressChan, errorChan := ollama.PullStream(ctx, "gemma3")

    for {
        select {
        case progress, ok := <-progressChan:
            if !ok {
                fmt.Println("拉取完成！")
                return
            }
            if progress.Total > 0 {
                percentage := float64(progress.Completed) / float64(progress.Total) * 100
                fmt.Printf("进度：%.1f%% (%s)\n", percentage, progress.Status)
            } else {
                fmt.Printf("状态：%s\n", progress.Status)
            }
        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

## API

### 客户端方法

- `Generate(ctx, req)` - 生成完成响应
- `GenerateStream(ctx, req)` - 生成流式完成响应
- `Chat(ctx, req)` - 发送聊天消息
- `ChatStream(ctx, req)` - 发送聊天消息并获取流式响应
- `Embed(ctx, req)` - 创建嵌入向量
- `Embeddings(ctx, req)` - 创建嵌入向量（传统 API）
- `List(ctx)` - 列出可用模型
- `Show(ctx, req)` - 显示模型信息
- `Pull(ctx, req)` - 下载模型
- `PullStream(ctx, req)` - 下载模型并显示进度
- `Push(ctx, req)` - 上传模型
- `PushStream(ctx, req)` - 上传模型并显示进度
- `Create(ctx, req)` - 从 Modelfile 创建模型
- `CreateStream(ctx, req)` - 创建模型并显示进度
- `Delete(ctx, req)` - 删除模型
- `Copy(ctx, req)` - 复制模型
- `Ps(ctx)` - 列出运行中的进程

### 全局函数

为了方便使用，所有客户端方法也可作为使用默认客户端实例的全局函数：

- `ollama.Generate(ctx, model, prompt, options...)`
- `ollama.Chat(ctx, model, messages, options...)`
- `ollama.Embed(ctx, model, input, options...)`
- 等等...

### 配置选项

可以使用选项函数配置客户端：

- `WithHost(host)` - 设置 Ollama 服务器 URL
- `WithHTTPClient(client)` - 使用自定义 HTTP 客户端
- `WithHeaders(headers)` - 添加自定义请求头

### 请求选项

许多函数支持用于常见配置的选项函数：

- `WithOptions(options)` - 设置模型选项
- `WithSystem(prompt)` - 设置系统提示
- `WithFormat(format)` - 设置响应格式
- `WithKeepAlive(duration)` - 设置保持连接时间
- `WithImages(images)` - 添加图像（用于多模态模型）
- `WithTools(tools)` - 添加工具进行函数调用
- `WithThinking()` - 启用思考模式

## 环境变量

- `OLLAMA_HOST` - 设置 Ollama 服务器 URL（默认：`http://localhost:11434`）

## 许可证

本项目采用 MIT 许可证。