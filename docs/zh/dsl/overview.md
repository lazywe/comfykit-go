# DSL 概述

ComfyKit 使用简单的 DSL（领域特定语言）来定义工作流中的参数和输出�?
## 语法

DSL 使用简单的 `$` 前缀来标记参数：

```json
{
  "inputs": {
    "text": "$prompt",
    "seed": "$seed",
    "image": "$~input_image"
  }
}
```

## 参数类型

### 基本参数

```json
"text": "$prompt"
```

接受任何值的基本参数�?
### 必填参数

```json
"text": "$prompt!"
```

标记参数为必填�?
### 媒体上传参数

```json
"image": "$~input_image"
```

标记参数需要媒体上传�?
### 必填媒体上传

```json
"image": "$~input_image!"
```

必填且需要媒体上传的参数�?
### 字段访问

```json
"text": "$prompt.text"
```

访问参数中的嵌套字段�?
## 输出标记

```json
{
  "inputs": {
    "filename_prefix": "$output.result"
  }
}
```

用变量名标记输出节点�?
## 完整示例

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "CLIPTextEncode",
      "inputs": {
        "text": "$prompt!"  // 必填参数
      }
    },
    {
      "id": "2",
      "type": "LoadImage",
      "inputs": {
        "image": "$~input_image"  // 媒体上传参数
      }
    },
    {
      "id": "3",
      "type": "SaveImage",
      "inputs": {
        "images": ["4", 0],
        "filename_prefix": "$output.result"  // 输出标记
      }
    }
  ]
}
```

## 优势

- **简�?*: 易于理解的语�?- **灵活**: 支持各种参数类型
- **自动**: ComfyKit 自动处理参数注入
- **验证**: 内置必填参数验证
- **媒体处理**: 自动媒体上传支持
