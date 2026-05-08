# 输出

学习如何定义和访问工作流输出�?
## 输出标记

使用 `$output.name` 标记输出节点�?
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "SaveImage",
      "inputs": {
        "images": ["2", 0],
        "filename_prefix": "$output.result"
      }
    }
  ]
}
```

## 访问输出

```go
result, err := kit.Execute("workflow.json", params)
if err != nil || result.Status != "completed" {
    return
}

// 访问所有图�?for _, img := range result.Images {
    fmt.Println(img)
}

// 按变量名访问
if images, ok := result.ImagesByVar["result"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}
```

## 多个输出

定义多个输出变量�?
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "SaveImage",
      "inputs": {
        "images": ["2", 0],
        "filename_prefix": "$output.main"
      }
    },
    {
      "id": "3",
      "type": "SaveImage",
      "inputs": {
        "images": ["4", 0],
        "filename_prefix": "$output.preview"
      }
    }
  ]
}
```

分别访问它们�?
```go
result, err := kit.Execute("workflow.json", params)

// 主输�?if images, ok := result.ImagesByVar["main"]; ok {
    fmt.Println("主图�?", len(images))
}

// 预览输出
if images, ok := result.ImagesByVar["preview"]; ok {
    fmt.Println("预览图片:", len(images))
}
```

## 输出类型

ComfyKit 支持多种输出类型�?
### 图片

```json
{
  "type": "SaveImage",
  "inputs": {
    "filename_prefix": "$output.images"
  }
}
```

### 音频

```json
{
  "type": "SaveAudio",
  "inputs": {
    "filename_prefix": "$output.audio"
  }
}
```

### 视频

```json
{
  "type": "SaveVideo",
  "inputs": {
    "filename_prefix": "$output.video"
  }
}
```

### 文本

```json
{
  "type": "SaveText",
  "inputs": {
    "filename_prefix": "$output.text"
  }
}
```

## 输出结构

```go
type ExecuteResult struct {
    // 扁平列表
    Images   []string
    Audios   []string
    Videos   []string
    Texts    []string
    
    // 按变量分�?    ImagesByVar  map[string][]string
    AudiosByVar  map[string][]string
    VideosByVar  map[string][]string
    TextsByVar   map[string][]string
    
    // 原始输出数据
    Outputs map[string]interface{}
}
```

## 最佳实�?
1. **使用描述性名�?*: 使用有意义的变量�?2. **分组相关输出**: 对相关输出使用相同的变量�?3. **处理所有输出类�?*: 检查图片、音频、视频和文本
4. **使用原始输出**: 访问 `Outputs` 获取自定义数�?5. **文档输出**: 添加注释描述每个输出包含的内�?