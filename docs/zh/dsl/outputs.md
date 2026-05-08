# 输出变量

## 输出变量语法

ComfyKit-Go 支持在工作流中定义输出变量：

```
$output.varname
```

## 使用示例

### 定义输出变量

在工作流节点标题中定义输出变量：

```
Save Image, $output.result
```

### 获取输出

执行工作流后，可以通过变量名获取输出：

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    // 处理错误
}

// 获取所有图片
images := result.Images

// 按变量名获取图片
resultImages := result.ImagesByVar["result"]
```

## 输出类型

ComfyKit-Go 支持多种输出类型：

### 图片输出

```go
// 获取所有图片
for _, imageURL := range result.Images {
    fmt.Println(imageURL)
}

// 按变量名获取图片
if images, ok := result.ImagesByVar["result"]; ok {
    for _, imageURL := range images {
        fmt.Println(imageURL)
    }
}
```

### 视频输出

```go
for _, videoURL := range result.Videos {
    fmt.Println(videoURL)
}
```

### 音频输出

```go
for _, audioURL := range result.Audios {
    fmt.Println(audioURL)
}
```

### 文本输出

```go
for _, text := range result.Texts {
    fmt.Println(text)
}
```

## 输出变量命名规范

### 推荐命名

- 使用小写字母
- 使用下划线分隔单词
- 避免使用特殊字符

**好的示例**：
```
$output.result
$output.main_image
$output.depth_map
```

**不好的示例**：
```
$output.Result          // 不推荐使用大写
$output.main-image      // 不推荐使用连字符
$output[result]         // 不推荐使用特殊字符
```

## 示例

### 工作流配置

```json
{
  "6": {
    "inputs": {
      "images": ["5", 0]
    },
    "class_type": "SaveImage",
    "_title": "Save Image, $output.result"
  }
}
```

### 代码使用

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫",
}

result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

// 获取输出
fmt.Printf("状态: %s\n", result.Status)
fmt.Printf("图片数量: %d\n", len(result.Images))

// 按变量名获取
if images, ok := result.ImagesByVar["result"]; ok {
    fmt.Printf("结果图片: %v\n", images)
}
```
