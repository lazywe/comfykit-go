# 参数语法

## 参数类型

ComfyKit-Go 支持多种参数类型：

### 1. 基本参数

```
$param
```

基本参数可以是任意类型的值。

**示例**：
```
Seed, $seed
```

### 2. 必填参数

```
$param!
```

必填参数如果没有提供会导致错误。

**示例**：
```
Prompt, $prompt!
```

### 3. 媒体参数

```
$~param
```

媒体参数用于上传文件（图片、音频、视频等）。

**示例**：
```
Image Input, $~image
```

### 4. 字段映射

```
$param.field
```

字段映射允许你将参数映射到工作流中的特定字段。

**示例**：
```
Model, $model.name
```

## 参数值类型

### 字符串

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫",
}
```

### 整数

```go
params := map[string]interface{}{
    "seed": 42,
}
```

### 浮点数

```go
params := map[string]interface{}{
    "cfg": 7.5,
}
```

### 布尔值

```go
params := map[string]interface{}{
    "enable": true,
}
```

### 数组

```go
params := map[string]interface{}{
    "steps": []int{10, 20, 30},
}
```

## 参数默认值

你可以为参数设置默认值：

```go
params := map[string]interface{}{
    "seed": 42, // 默认值
}
```

如果用户没有提供参数，将使用默认值。

## 参数验证

### 必填参数验证

如果标记为必填的参数没有提供，会返回错误：

```go
result, err := kit.Execute("workflow.json", map[string]interface{}{})
// 如果 workflow.json 包含 $prompt!，但没有提供 prompt 参数，会返回错误
```

### 类型验证

ComfyKit-Go 会尝试将参数转换为正确的类型：

```go
params := map[string]interface{}{
    "seed": "42", // 字符串会自动转换为整数
}
```

## 示例

### 工作流节点标题

```
Prompt: $prompt!, Seed: $seed, Steps: $steps
```

### 执行代码

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
    "steps":  30,
}
result, err := kit.Execute("workflow.json", params)
```
