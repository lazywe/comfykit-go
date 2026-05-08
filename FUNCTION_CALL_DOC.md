# ComfyKit Python 函数调用文档

## 概述

ComfyKit 是一个统一�?API，用于执�?ComfyUI 工作流。支持本�?ComfyUI 执行�?RunningHub 云端执行�?

---

## 一、核心类：ComfyKit

### 1.1 构造函�?

```python
def __init__(
    self,
    # ComfyUI 配置（本地执行）
    comfyui_url: Optional[str] = None,      # 默认: "http://127.0.0.1:8188"
    executor_type: Literal["http", "websocket"] = "http",  # 默认: "http"
    api_key: Optional[str] = None,           # ComfyUI API key
    cookies: Optional[str] = None,           # ComfyUI cookies
    
    # RunningHub 配置（云端执行）
    runninghub_url: Optional[str] = None,    # 默认: "https://www.runninghub.ai"
    runninghub_api_key: Optional[str] = None, # RunningHub API key
    runninghub_timeout: Optional[int] = None, # 任务超时（秒�?
    runninghub_retry_count: int = 3,         # API重试次数
    runninghub_instance_type: Optional[str] = None,  # 实例类型（如 "plus"�?
)
```

**配置优先�?*: 构造函数参�?> 环境变量 > 默认�?

### 1.2 核心方法

#### execute
```python
async def execute(
    self,
    workflow: Union[str, Path],      # 工作流源（文件路�?URL/RunningHub ID�?
    params: Optional[Dict[str, Any]] = None  # 参数（统一简单格式）
) -> ExecuteResult
```

**自动检测工作流类型**:
- **RunningHub ID**: 纯数字字符串（如 "12345"�?
- **URL**: �?`http://` �?`https://` 开�?
- **文件路径**: 存在的文件或包含 `/`、`\` 的路�?
- **RunningHub工作流文�?*: 包含 `_source: "runninghub"` 的本地文�?

#### execute_json
```python
async def execute_json(
    self,
    workflow_json: Dict[str, Any],   # 工作流JSON字典
    params: Optional[Dict[str, Any]] = None
) -> ExecuteResult
```

#### close
```python
async def close()  # 关闭所有执行器和资�?
```

---

## 二、数据模�?

### 2.1 ExecuteResult

```python
class ExecuteResult(BaseModel):
    status: str                      # 执行状�?
    prompt_id: Optional[str]         # 任务ID
    duration: Optional[float]        # 执行时长（秒�?
    
    images: List[str]                # 图片URL列表
    images_by_var: Dict[str, List[str]]  # 按变量名分组的图�?
    
    audios: List[str]                # 音频URL列表
    audios_by_var: Dict[str, List[str]]  # 按变量名分组的音�?
    
    videos: List[str]                # 视频URL列表
    videos_by_var: Dict[str, List[str]]  # 按变量名分组的视�?
    
    texts: List[str]                 # 文本内容列表
    texts_by_var: Dict[str, List[str]]   # 按变量名分组的文�?
    
    outputs: Optional[Dict[str, Any]]    # 原始输出数据
    msg: Optional[str]               # 消息（错误信息等�?
```

### 2.2 WorkflowParam

```python
class WorkflowParam(BaseModel):
    name: str              # 参数�?
    type: str = "str"      # 参数类型（str/int/float/bool�?
    required: bool = True  # 是否必填
    default: Optional[Any] # 默认�?
    need_upload: bool = False  # 是否需要上�?
```

### 2.3 WorkflowParamMapping

```python
class WorkflowParamMapping(BaseModel):
    param_name: str         # 参数�?
    node_id: str            # 节点ID
    input_field: str        # 输入字段�?
    node_class_type: str    # 节点类型
    need_upload: bool = False  # 是否需要上�?
```

### 2.4 WorkflowOutputMapping

```python
class WorkflowOutputMapping(BaseModel):
    node_id: str     # 节点ID
    output_var: str  # 输出变量�?
```

### 2.5 WorkflowMetadata

```python
class WorkflowMetadata(BaseModel):
    title: str                    # 工作流标�?
    params: Dict[str, WorkflowParam]  # 参数字典
    mapping_info: WorkflowMappingInfo # 映射信息
    workflow_id: Optional[str]    # RunningHub工作流ID
    is_runninghub: bool = False   # 是否为RunningHub工作�?
```

---

## 三、执行器基类：ComfyUIExecutor

### 3.1 构造函�?

```python
def __init__(self, base_url: str = None, api_key: str = None, cookies: str = None)
```

### 3.2 核心方法

#### execute_workflow（抽象方法）
```python
async def execute_workflow(self, workflow_file: str, params: Dict[str, Any] = None) -> ExecuteResult
```

#### _randomize_seed_in_workflow
```python
def _randomize_seed_in_workflow(self, workflow_data: Dict[str, Any]) -> tuple[Dict[str, Any], Dict[str, int]]
```
将工作流中所�?`seed == 0` 的节点替换为随机63位种子�?

#### _apply_params_to_workflow
```python
async def _apply_params_to_workflow(self, workflow_data: Dict[str, Any], metadata: WorkflowMetadata, params: Dict[str, Any]) -> Dict[str, Any]
```
应用参数到工作流，支持默认值和必填校验�?

#### _extract_output_nodes
```python
def _extract_output_nodes(self, metadata: WorkflowMetadata) -> Dict[str, str]
```
从元数据中提取输出节点映射（node_id -> 变量名）�?

#### _split_media_by_suffix
```python
def _split_media_by_suffix(self, node_output: Dict[str, Any], base_url: str) -> Tuple[List[str], List[str], List[str]]
```
按文件扩展名分类媒体（图�?视频/音频）�?

#### _map_outputs_by_var
```python
def _map_outputs_by_var(self, output_id_2_var: Dict[str, str], output_id_2_media: Dict[str, List[str]]) -> Dict[str, List[str]]
```
按变量名映射输出�?

#### _extend_flat_list_from_dict
```python
def _extend_flat_list_from_dict(self, media_dict: Dict[str, List[str]]) -> List[str]
```
将字典中的列表展平为单一列表�?

### 3.3 媒体上传相关方法

#### _handle_media_upload
```python
async def _handle_media_upload(self, node_data: Dict[str, Any], input_field: str, param_value: Any)
```
处理媒体上传（支持URL和本地文件）�?

#### _upload_media_from_source
```python
async def _upload_media_from_source(self, media_url: str) -> str
```
从URL下载媒体并上传到ComfyUI�?

#### _upload_media
```python
async def _upload_media(self, media_path: str) -> str
```
上传本地媒体文件到ComfyUI�?

---

## 四、HTTP执行器：HttpExecutor

### 4.1 构造函�?

```python
def __init__(self, base_url: str = None, api_key: str = None, cookies: str = None)
```

### 4.2 核心方法

#### execute_workflow
```python
async def execute_workflow(self, workflow_file: str, params: Dict[str, Any] = None) -> ExecuteResult
```

#### _queue_prompt
```python
async def _queue_prompt(self, workflow: Dict[str, Any], client_id: str, prompt_ext_params: Optional[Dict[str, Any]] = None) -> str
```
提交工作流到ComfyUI队列，返回prompt_id�?

#### _wait_for_results
```python
async def _wait_for_results(self, prompt_id: str, client_id: str, timeout: Optional[int] = None, output_id_2_var: Optional[Dict[str, str]] = None) -> ExecuteResult
```
轮询等待执行结果�?

---

## 五、WebSocket执行器：WebSocketExecutor

### 5.1 构造函�?

```python
def __init__(self, base_url: str = None, api_key: str = None, cookies: str = None)
```

### 5.2 核心方法

#### execute_workflow
```python
async def execute_workflow(self, workflow_file: str, params: Dict[str, Any] = None) -> ExecuteResult
```

#### _parse_ws_url
```python
def _parse_ws_url()  # 解析WebSocket URL
```

#### _parse_ws_message
```python
def _parse_ws_message(self, message: dict, prompt_id: str) -> tuple[bool, dict]
```
解析WebSocket消息，判断执行是否完成�?

#### _build_result_from_collected_outputs
```python
def _build_result_from_collected_outputs(self, collected_outputs: Dict[str, Any], prompt_id: str, output_id_2_var: Optional[Dict[str, str]] = None) -> ExecuteResult
```
从收集的输出构建执行结果�?

---

## 六、RunningHub执行器：RunningHubExecutor

### 6.1 构造函�?

```python
def __init__(self, base_url: str = None, api_key: str = None, timeout: int = None, retry_count: int = 3, instance_type: str = None)
```

### 6.2 核心方法

#### execute_by_id
```python
async def execute_by_id(self, workflow_id: str, params: Dict[str, Any] = None) -> ExecuteResult
```
通过RunningHub工作流ID直接执行�?

#### execute_workflow
```python
async def execute_workflow(self, workflow_file: str, params: Dict[str, Any] = None) -> ExecuteResult
```
执行本地工作流文件（包含RunningHub配置）�?

#### _convert_params_to_node_info_list
```python
async def _convert_params_to_node_info_list(self, metadata, params: dict, seed_changes: Dict[str, int] = None) -> List[dict]
```
将参数转换为RunningHub nodeInfoList格式�?

#### _wait_for_task_completion
```python
async def _wait_for_task_completion(self, task_id: str, output_id_2_var: Optional[Dict[str, str]] = None, max_wait_time: int = None) -> ExecuteResult
```
等待RunningHub任务完成�?

#### _process_task_result
```python
async def _process_task_result(self, task_id: str, result_data: List[Dict[str, Any]], output_id_2_var: Optional[Dict[str, str]] = None) -> ExecuteResult
```
处理RunningHub任务结果，转换为ExecuteResult格式�?

---

## 七、RunningHub客户端：RunningHubClient

### 7.1 构造函�?

```python
def __init__(self, api_key: str = None, base_url: str = None, timeout: int = None, retry_count: int = None, instance_type: str = None)
```

### 7.2 API方法

#### get_workflow_json
```python
async def get_workflow_json(self, workflow_id: str) -> Dict[str, Any]
```
获取工作流JSON定义�?

#### upload_file
```python
async def upload_file(self, file_path: str) -> str
```
上传文件到RunningHub，返回fileName�?

#### create_task
```python
async def create_task(self, workflow_id: str, node_info_list: List[Dict] = None) -> Dict[str, Any]
```
创建工作流执行任务�?

#### query_task_status
```python
async def query_task_status(self, task_id: str) -> Dict[str, Any]
```
查询任务状态�?

#### query_task_result
```python
async def query_task_result(self, task_id: str) -> List[Dict[str, Any]]
```
查询任务执行结果�?

#### close
```python
async def close()  # 关闭会话
```

---

## 八、工作流解析器：WorkflowParser

### 8.1 DSL语法

| 语法 | 说明 | 示例 |
|------|------|------|
| `$param` | 基础参数 | `$prompt` |
| `$param!` | 必填参数 | `$prompt!` |
| `$~param` | 需要上�?| `$~image` |
| `$~param!` | 必填且需上传 | `$~image!` |
| `$param.field` | 指定字段 | `$prompt.text` |
| `$output.name` | 输出标记 | `$output.result` |

### 8.2 核心方法

#### parse_workflow_file
```python
def parse_workflow_file(self, file_path: str, tool_name: Optional[str] = None) -> Optional[WorkflowMetadata]
```
解析工作流文件，提取参数和输出映射�?

#### parse_workflow
```python
def parse_workflow(self, workflow_data: Dict[str, Any], title: str) -> Optional[WorkflowMetadata]
```
解析工作流JSON数据�?

#### parse_node
```python
def parse_node(self, node_id: str, node_data: Dict[str, Any]) -> tuple[List[WorkflowParam], List[WorkflowParamMapping], Optional[WorkflowOutputMapping]]
```
解析单个节点的参数和输出�?

#### parse_param_marker
```python
def parse_param_marker(self, marker: str) -> Optional[Dict[str, Any]]
```
解析参数标记（如 `$prompt.text!`）�?

---

## 九、媒体上传节点类�?

```python
MEDIA_UPLOAD_NODE_TYPES = {
    'LoadImage',
    'LoadAudio',
    'LoadVideo',
    'VHS_LoadAudioUpload',
    'VHS_LoadVideo',
}
```

---

## 十、输出节点类�?

```python
known_output_types = {
    'SaveImage',
    'SaveVideo',
    'SaveAudio',
    'VHS_SaveVideo',
    'VHS_SaveAudio'
}
```

---

## 十一、文件扩展名分类

| 类型 | 扩展�?|
|------|--------|
| 图片 | `.png`, `.jpg`, `.jpeg`, `.webp`, `.bmp`, `.tiff` |
| 视频 | `.mp4`, `.mov`, `.avi`, `.webm`, `.gif` |
| 音频 | `.mp3`, `.wav`, `.flac`, `.ogg`, `.aac`, `.m4a`, `.wma`, `.opus` |
