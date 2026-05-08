# DSL 最佳实践

## 参数命名

### 使用描述性名称

为参数使用描述性名称，便于理解和维护：

**好的示例**：
```
$prompt, $seed, $steps
```

**不好的示例**：
```
$p, $s, $st
```

### 保持一致性

在整个工作流中保持参数命名一致：

```
# 始终使用相同的命名风格
$prompt, $negative_prompt, $seed, $steps
```

## 必填参数

### 标记关键参数为必填

将影响结果的关键参数标记为必填：

```
$prompt!, $model!
```

### 提供默认值

为可选参数提供默认值：

```go
params := map[string]interface{}{
    "seed":  42,
    "steps": 20,
    "cfg":   7.5,
}
```

## 媒体参数

### 使用媒体参数标记

对于需要上传文件的参数，使用媒体参数标记：

```
$~image, $~video
```

### 文件路径处理

确保文件路径正确：

```go
params := map[string]interface{}{
    "image": "path/to/image.png", // 使用绝对路径或相对于工作目录的路径
}
```

## 输出变量

### 使用有意义的变量名

为输出变量使用有意义的名称：

```
$output.main_image
$output.depth_map
$output.mask
```

### 组织输出

按功能组织输出变量：

```
$output.result           // 主要结果
$output.preview          // 预览图
$output.intermediate     // 中间结果
```

## 工作流设计

### 保持简洁

避免在单个节点中定义过多参数：

**好的示例**：
```
Prompt: $prompt!, Seed: $seed
```

**不好的示例**：
```
Prompt: $prompt!, Seed: $seed, Steps: $steps, CFG: $cfg, Model: $model!, Lora: $lora, Clip Skip: $clip_skip
```

### 使用注释

在工作流中添加注释说明：

```json
{
  "_comment": "Text Encoder Node",
  "_title": "Prompt: $prompt!, Seed: $seed"
}
```

## 错误处理

### 检查必填参数

确保所有必填参数都已提供：

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫", // 必填
}

result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("参数错误: %v\n", err)
    return
}
```

### 处理执行错误

检查执行结果状态：

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("执行错误: %v\n", err)
    return
}

if result.Status == "error" {
    fmt.Printf("工作流错误: %s\n", result.Message)
    return
}
```

## 性能考虑

### 减少参数数量

过多的参数会增加解析时间，尽量精简：

```
# 只定义必要的参数
$prompt!, $seed
```

### 使用默认值

为常用参数设置合理的默认值，减少用户输入：

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫",
    // seed 使用默认值 42
    // steps 使用默认值 20
}
```

## 文档

### 文档化参数

为工作流添加参数文档：

```json
{
  "_parameters": {
    "prompt": "生成图像的提示词",
    "seed": "随机种子，用于复现结果",
    "steps": "采样步数"
  }
}
```

### 提供示例

为工作流提供使用示例：

```json
{
  "_examples": [
    {
      "name": "生成猫图片",
      "parameters": {
        "prompt": "一只可爱的猫",
        "seed": 42
      }
    }
  ]
}
```
