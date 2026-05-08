# DSL 语法概述

ComfyKit-Go 支持简单的 DSL（领域特定语言）语法，用于在工作流节点标题中定义参数。

## 基本语法

### 参数标记

| 语法 | 描述 | 示例 |
|------|------|------|
| `$param` | 基本参数 | `$prompt` |
| `$param!` | 必填参数 | `$prompt!` |
| `$~param` | 需要媒体上传的参数 | `$~image` |
| `$param.field` | 映射到特定字段的参数 | `$model.name` |
| `$output.varname` | 输出变量名 | `$output.result` |

### 使用示例

在 ComfyUI 工作流中，你可以在节点标题中使用这些标记：

```
Prompt, $prompt!, $seed
```

这个节点标题定义了两个参数：
- `$prompt!` - 必填的提示词参数
- `$seed` - 可选的随机种子参数

## 参数类型推断

ComfyKit-Go 会自动推断参数类型：

- 字符串（如 `"hello"`）
- 整数（如 `42`）
- 浮点数（如 `3.14`）
- 布尔值（`true` 或 `false`）
- 数组（如 `[1, 2, 3]`）
- 对象（如 `{"key": "value"}`）

## 工作流示例

```json
{
  "3": {
    "inputs": {
      "text": "$prompt!",
      "seed": "$seed"
    },
    "class_type": "CLIPTextEncode",
    "_title": "Prompt, $prompt!, $seed"
  }
}
```

## 下一步

- 查看 [参数](parameters.md) 了解更多参数语法
- 查看 [输出](outputs.md) 了解输出变量
- 查看 [最佳实践](best-practices.md) 了解使用建议
