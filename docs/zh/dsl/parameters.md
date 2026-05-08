# 参数

学习如何在工作流中定义和使用参数�?
## 基本参数

使用 `$` 前缀定义参数�?
```json
{
  "inputs": {
    "text": "$prompt",
    "seed": "$seed",
    "steps": "$steps"
  }
}
```

执行时传递值：

```go
params := map[string]interface{}{
    "prompt": "美丽的日�?,
    "seed":   12345,
    "steps":  25,
}

result, err := kit.Execute("workflow.json", params)
```

## 必填参数

使用 `!` 标记必填参数�?
```json
{
  "inputs": {
    "text": "$prompt!"
  }
}
```

如果未提供必填参数，ComfyKit �?panic�?
```go
// 这会 panic，因�?"prompt" 是必填的
result, err := kit.Execute("workflow.json", nil)
```

## 媒体上传参数

使用 `~` 标记需要媒体上传的参数�?
```json
{
  "inputs": {
    "image": "$~input_image"
  }
}
```

ComfyKit 将自动上传媒体文件：

```go
params := map[string]interface{}{
    "input_image": "/path/to/image.png",      // 本地文件
    // 或�?    "input_image": "https://example.com/img.jpg",  // URL
}
```

## 必填媒体上传

组合两个标记�?
```json
{
  "inputs": {
    "image": "$~input_image!"
  }
}
```

## 字段访问

使用点号访问嵌套字段�?
```json
{
  "inputs": {
    "text": "$prompt.text"
  }
}
```

使用嵌套参数�?
```go
params := map[string]interface{}{
    "prompt": map[string]interface{}{
        "text": "美丽的日�?,
        "style": "电影风格",
    },
}
```

## 支持的参数类�?
| 类型 | 描述 | 示例 |
|------|------|------|
| `str` | 字符�?| `"hello"` |
| `int` | 整数 | `42` |
| `float` | 浮点�?| `3.14` |
| `bool` | 布尔�?| `true` |
| `image` | 图片路径/URL | `/path/to/img.png` |
| `audio` | 音频路径/URL | `/path/to/audio.mp3` |
| `video` | 视频路径/URL | `/path/to/video.mp4` |

## 默认�?
您可以在工作流元数据中定义默认值：

```json
{
  "__metadata__": {
    "params": {
      "prompt": {
        "type": "str",
        "required": true,
        "default": "美丽的场�?
      },
      "seed": {
        "type": "int",
        "required": false,
        "default": 0
      }
    }
  }
}
```

## 参数验证

ComfyKit 验证参数�?
1. **必填检�?*: 确保提供必填参数
2. **类型检�?*: 尝试将值转换为预期类型
3. **媒体检�?*: 验证媒体文件存在�?URL 有效

## 最佳实�?
1. **使用描述性名�?*: 使用清晰的参数名称如 `prompt`, `seed`, `steps`
2. **标记必填参数**: 对必填参数使�?`!`
3. **使用媒体上传**: 对图�?音频/视频输入使用 `~`
4. **提供默认�?*: 尽可能在元数据中定义默认�?5. **文档参数**: 添加注释或元数据描述参数
