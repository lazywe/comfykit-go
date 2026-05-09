package comfyui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// RunningHubExecutor 是用于在 RunningHub 云平台上执行 ComfyUI 工作流的执行器
// 它封装了与 RunningHub API 的交互逻辑，包括任务创建、状态查询和结果获取
type RunningHubExecutor struct {
	*ComfyUIExecutor                   // 嵌入基础执行器，继承参数映射能力
	timeout          int               // 任务超时时间（秒）
	retryCount       int               // API 请求重试次数
	instanceType     string            // RunningHub 实例类型（如 "plus", "pro"）
	client           *RunningHubClient // RunningHub API 客户端
}

// NewRunningHubExecutor 创建一个新的 RunningHub 执行器实例
// 参数：
//
//	baseURL - RunningHub API 基础地址
//	apiKey - RunningHub API 密钥
//	timeout - 任务超时时间（秒）
//	retryCount - API 请求重试次数
//	instanceType - RunningHub 实例类型
func NewRunningHubExecutor(baseURL, apiKey string, timeout, retryCount int, instanceType string) *RunningHubExecutor {
	return &RunningHubExecutor{
		ComfyUIExecutor: NewComfyUIExecutor(baseURL, apiKey, ""),
		timeout:         timeout,
		retryCount:      retryCount,
		instanceType:    instanceType,
		client:          NewRunningHubClient(apiKey, baseURL, timeout, retryCount, instanceType),
	}
}

// Close 关闭执行器，释放资源
func (e *RunningHubExecutor) Close() error {
	return e.client.Close()
}

// prepareAndCreateTask 准备工作流并创建任务（公共代码）
// 参数：
//
//	workflowID - RunningHub 工作流 ID
//	params - 用户提供的参数映射
//
// 返回：任务 ID、输出节点映射、工作流元数据、错误信息
func (e *RunningHubExecutor) prepareAndCreateTask(workflowID string, params map[string]interface{}) (string, map[string]string, *WorkflowMetadata, error) {
	// 从 RunningHub 获取工作流 JSON
	workflowJSON, err := e.client.GetWorkflowJSON(workflowID)
	if err != nil {
		return "", nil, nil, err
	}

	// 随机化工作流中的种子值
	workflowJSON, seedChanges := randomizeSeedInWorkflow(workflowJSON)

	// 解析工作流元数据
	parser := NewWorkflowParser()
	metadata := parser.ParseWorkflow(workflowJSON, "workflow_"+workflowID)
	if metadata == nil {
		return "", nil, nil, fmt.Errorf("failed to parse workflow metadata")
	}

	metadata.WorkflowID = workflowID
	metadata.IsRunningHub = true

	// 将用户参数转换为节点信息列表
	nodeInfoList := e.ConvertParamsToNodeInfoList(metadata, params, seedChanges)
	// 提取输出节点映射
	outputID2Var := extractOutputNodes(metadata)

	// 在 RunningHub 上创建任务
	taskData, err := e.client.CreateTask(workflowID, nodeInfoList)
	if err != nil {
		return "", nil, nil, err
	}

	// 提取任务 ID
	taskID, ok := taskData["taskId"].(string)
	if !ok || taskID == "" {
		return "", nil, nil, fmt.Errorf("failed to get task ID")
	}

	return taskID, outputID2Var, metadata, nil
}

// ExecuteByID 根据工作流 ID 在 RunningHub 上执行工作流
// 参数：
//
//	workflowID - RunningHub 工作流 ID
//	params - 用户提供的参数映射
//
// 返回：执行结果和错误信息
func (e *RunningHubExecutor) ExecuteByID(workflowID string, params map[string]interface{}) (*ExecuteResult, error) {
	startTime := time.Now()

	// 准备并创建任务
	taskID, outputID2Var, _, err := e.prepareAndCreateTask(workflowID, params)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to prepare task: " + err.Error()), nil
	}

	// 等待任务完成
	result, err := e.WaitForTaskCompletion(taskID, outputID2Var)
	if err != nil {
		return result, err
	}

	// 设置执行时长
	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

// ExecuteAsyncByID 异步创建 RunningHub 任务
// 参数：
//
//	workflowID - RunningHub 工作流 ID
//	params - 用户提供的参数映射
//
// 返回：任务 ID、输出节点映射、错误信息
func (e *RunningHubExecutor) ExecuteAsyncByID(workflowID string, params map[string]interface{}) (string, map[string]string, error) {
	// 准备并创建任务
	taskID, outputID2Var, _, err := e.prepareAndCreateTask(workflowID, params)
	return taskID, outputID2Var, err
}

// GetTaskCompletion 检查任务状态并在完成时获取结果
// 参数：
//
//	taskID - 任务 ID
//	outputID2Var - 输出节点 ID 到变量名的映射
//
// 返回：执行结果、是否完成、错误信息
func (e *RunningHubExecutor) GetTaskCompletion(taskID string, outputID2Var map[string]string) (*ExecuteResult, bool, error) {
	// 查询任务状态
	statusInfo, err := e.client.QueryTaskStatus(taskID)
	if err != nil {
		return nil, false, err
	}

	taskStatus, _ := statusInfo["status"].(string)
	statusMsg, _ := statusInfo["msg"].(string)

	switch taskStatus {
	case "SUCCESS":
		// 任务成功，获取结果
		resultData, err := e.client.QueryTaskResult(taskID)
		if err != nil {
			return NewExecuteResult("error").ErrorResult("Failed to get task result: " + err.Error()), true, nil
		}
		result := e.ProcessTaskResult(taskID, resultData, outputID2Var)
		return result, true, nil

	case "FAILED":
		// 任务失败
		errorMsg := "RunningHub task " + taskID + " failed"
		if statusMsg != "" {
			errorMsg += ": " + statusMsg
		}
		return NewExecuteResult("error").ErrorResult(errorMsg), true, nil

	case "QUEUED", "RUNNING":
		// 任务还在进行中
		return nil, false, nil

	default:
		// 未知状态
		return NewExecuteResult("error").ErrorResult("Unknown task status: " + taskStatus), true, nil
	}
}

// ExecuteWorkflow 从本地文件加载工作流并在 RunningHub 上执行
// 参数：
//
//	workflowFile - 本地工作流文件路径
//	params - 用户提供的参数映射
//
// 返回：执行结果和错误信息
func (e *RunningHubExecutor) ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error) {
	startTime := time.Now()

	// 检查工作流文件是否存在
	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		return NewExecuteResult("error").ErrorResult("Workflow file does not exist: " + workflowFile), nil
	}

	// 加载工作流 JSON
	workflowJSON, err := loadWorkflowJSON(workflowFile)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to load workflow: " + err.Error()), nil
	}

	// 随机化工作流中的种子值
	workflowJSON, seedChanges := randomizeSeedInWorkflow(workflowJSON)

	// 解析工作流元数据
	parser := NewWorkflowParser()
	metadata := parser.ParseWorkflow(workflowJSON, "")
	if metadata == nil {
		return NewExecuteResult("error").ErrorResult("Cannot parse workflow metadata"), nil
	}

	// 获取工作流 ID（必须在元数据中定义）
	workflowID := metadata.WorkflowID
	if workflowID == "" {
		return NewExecuteResult("error").ErrorResult("RunningHub workflow_id not found in metadata"), nil
	}

	// 将用户参数转换为节点信息列表
	nodeInfoList := e.ConvertParamsToNodeInfoList(metadata, params, seedChanges)

	// 在 RunningHub 上创建任务
	taskData, err := e.client.CreateTask(workflowID, nodeInfoList)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to create task: " + err.Error()), nil
	}

	// 提取任务 ID
	taskID, ok := taskData["taskId"].(string)
	if !ok || taskID == "" {
		return NewExecuteResult("error").ErrorResult("Failed to get task ID"), nil
	}

	// 提取输出节点映射并等待任务完成
	outputID2Var := extractOutputNodes(metadata)
	result, err := e.WaitForTaskCompletion(taskID, outputID2Var)
	if err != nil {
		return result, err
	}

	// 设置执行时长
	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

// ConvertParamsToNodeInfoList 将用户参数转换为 RunningHub API 所需的节点信息列表
// 参数：
//
//	metadata - 工作流元数据
//	params - 用户提供的参数映射
//	seedChanges - 随机化后的种子值映射
//
// 返回：节点信息列表
func (e *RunningHubExecutor) ConvertParamsToNodeInfoList(metadata *WorkflowMetadata, params map[string]interface{}, seedChanges map[string]int64) []map[string]interface{} {
	nodeInfoList := []map[string]interface{}{}

	// 处理用户参数映射
	for _, mapping := range metadata.MappingInfo.ParamMappings {
		if val, ok := params[mapping.ParamName]; ok {
			nodeClassType := mapping.NodeClassType
			needUpload := mapping.NeedUpload

			// 如果需要上传媒体文件，处理文件上传
			if needUpload || IsMediaUploadNode(nodeClassType) {
				val = e.HandleRunningHubMediaUpload(val)
			}

			nodeInfo := map[string]interface{}{
				"nodeId":     mapping.NodeID,
				"fieldName":  mapping.InputField,
				"fieldValue": val,
			}
			nodeInfoList = append(nodeInfoList, nodeInfo)
		}
	}

	// 处理随机化的种子值
	for nodeID, seedValue := range seedChanges {
		alreadySet := false
		// 检查种子是否已经被用户设置
		for _, ni := range nodeInfoList {
			if ni["nodeId"] == nodeID && ni["fieldName"] == "seed" {
				alreadySet = true
				break
			}
		}
		if !alreadySet {
			nodeInfoList = append(nodeInfoList, map[string]interface{}{
				"nodeId":     nodeID,
				"fieldName":  "seed",
				"fieldValue": seedValue,
			})
		}
	}

	return nodeInfoList
}

// HandleRunningHubMediaUpload 处理媒体文件上传
// 支持从 URL 下载并上传，或直接上传本地文件
// 参数：
//
//	paramValue - 媒体文件路径或 URL
//
// 返回：上传后的文件名或原始值（如果上传失败）
func (e *RunningHubExecutor) HandleRunningHubMediaUpload(paramValue interface{}) interface{} {
	if strVal, ok := paramValue.(string); ok {
		// 如果是 URL，先下载再上传
		if strings.HasPrefix(strVal, "http://") || strings.HasPrefix(strVal, "https://") {
			mediaValue, err := e.UploadMediaFromURL(strVal)
			if err == nil {
				return mediaValue
			}
		} else if _, err := os.Stat(strVal); err == nil {
			// 如果是本地文件，直接上传
			uploadedName, err := e.client.UploadFile(strVal)
			if err == nil {
				return uploadedName
			}
		}
	}
	return paramValue
}

// UploadMediaFromURL 从 URL 下载媒体文件并上传到 RunningHub
// 参数：
//
//	mediaURL - 媒体文件的 URL
//
// 返回：上传后的文件名和错误信息
func (e *RunningHubExecutor) UploadMediaFromURL(mediaURL string) (string, error) {
	// 下载文件到临时路径
	tempPath, err := downloadFile(mediaURL)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempPath) // 确保函数结束后删除临时文件
	return e.client.UploadFile(tempPath)
}

// downloadFile 从 URL 下载文件到临时目录
// 参数：
//
//	url - 要下载的文件 URL
//
// 返回：临时文件路径和错误信息
func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "workflow_*.json")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 读取响应内容并写入临时文件
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 || err != nil {
			break
		}
		tempFile.Write(buf[:n])
	}

	return tempFile.Name(), nil
}

// WaitForTaskCompletion 等待 RunningHub 任务完成
// 参数：
//
//	taskID - 任务 ID
//	outputID2Var - 输出节点 ID 到变量名的映射
//
// 返回：执行结果和错误信息
func (e *RunningHubExecutor) WaitForTaskCompletion(taskID string, outputID2Var map[string]string) (*ExecuteResult, error) {
	maxWaitTime := e.timeout // 最大等待时间（秒）
	checkInterval := 2       // 状态检查间隔（秒）
	startTime := time.Now()

	for {
		// 检查是否超时
		elapsedTime := time.Since(startTime).Seconds()
		if maxWaitTime > 0 && elapsedTime >= float64(maxWaitTime) {
			return NewExecuteResult("error").ErrorResult("RunningHub task timeout after " + fmt.Sprintf("%d", maxWaitTime) + " seconds"), nil
		}

		// 使用 GetTaskCompletion 检查任务状态
		result, completed, err := e.GetTaskCompletion(taskID, outputID2Var)
		if err != nil {
			// 遇到错误，等待一下再试
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		if completed {
			return result, nil
		}

		// 任务还在运行中，继续等待
		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}

// ProcessTaskResult 处理 RunningHub 任务结果，将原始结果转换为 ExecuteResult
// 参数：
//
//	taskID - 任务 ID
//	resultData - 原始结果数据
//	outputID2Var - 输出节点 ID 到变量名的映射
//
// 返回：格式化后的执行结果
func (e *RunningHubExecutor) ProcessTaskResult(taskID string, resultData interface{}, outputID2Var map[string]string) *ExecuteResult {
	result := NewExecuteResult("completed")
	result.PromptID = taskID

	// 按类型分组输出
	outputID2Images := map[string][]string{}
	outputID2Videos := map[string][]string{}
	outputID2Audios := map[string][]string{}
	outputID2Texts := map[string][]string{}

	// 解析结果列表
	if list, ok := resultData.([]interface{}); ok {
		for idx, item := range list {
			if itemMap, ok := item.(map[string]interface{}); ok {
				fileURL, _ := itemMap["fileUrl"].(string)
				fileType, _ := itemMap["fileType"].(string)
				nodeID, _ := itemMap["nodeId"].(string)
				if nodeID == "" {
					nodeID = fmt.Sprintf("%d", idx) // 如果没有 nodeId，使用索引作为标识
				}

				if fileURL != "" {
					fileType = strings.ToLower(fileType)
					// 根据文件类型分类
					switch {
					case strings.Contains(fileType, "image") ||
						strings.Contains(fileType, "png") ||
						strings.Contains(fileType, "jpg") ||
						strings.Contains(fileType, "jpeg") ||
						strings.Contains(fileType, "gif") ||
						strings.Contains(fileType, "webp"):
						outputID2Images[nodeID] = append(outputID2Images[nodeID], fileURL)

					case strings.Contains(fileType, "video") ||
						strings.Contains(fileType, "mp4") ||
						strings.Contains(fileType, "avi") ||
						strings.Contains(fileType, "mov") ||
						strings.Contains(fileType, "mkv"):
						outputID2Videos[nodeID] = append(outputID2Videos[nodeID], fileURL)

					case strings.Contains(fileType, "audio") ||
						strings.Contains(fileType, "mp3") ||
						strings.Contains(fileType, "wav") ||
						strings.Contains(fileType, "flac"):
						outputID2Audios[nodeID] = append(outputID2Audios[nodeID], fileURL)

					case strings.Contains(fileType, "text") ||
						strings.Contains(fileType, "txt") ||
						strings.Contains(fileType, "json") ||
						strings.Contains(fileType, "xml"):
						textContent := e.DownloadTextFromURL(fileURL)
						if textContent != "" {
							outputID2Texts[nodeID] = append(outputID2Texts[nodeID], textContent)
						}
					}
				}
			}
		}
	}

	// 将分类结果设置到执行结果中
	if len(outputID2Images) > 0 {
		result.ImagesByVar = mapOutputsByVar(outputID2Var, outputID2Images)
		result.Images = extendFlatListFromDict(result.ImagesByVar)
	}
	if len(outputID2Videos) > 0 {
		result.VideosByVar = mapOutputsByVar(outputID2Var, outputID2Videos)
		result.Videos = extendFlatListFromDict(result.VideosByVar)
	}
	if len(outputID2Audios) > 0 {
		result.AudiosByVar = mapOutputsByVar(outputID2Var, outputID2Audios)
		result.Audios = extendFlatListFromDict(result.AudiosByVar)
	}
	if len(outputID2Texts) > 0 {
		result.TextsByVar = mapOutputsByVar(outputID2Var, outputID2Texts)
		result.Texts = extendFlatListFromDict(result.TextsByVar)
	}

	// 保存原始数据
	result.Outputs = map[string]interface{}{"raw_data": resultData}
	return result
}

// DownloadTextFromURL 从 URL 下载文本内容
// 参数：
//
//	textURL - 文本文件的 URL
//
// 返回：文本内容（如果下载失败返回空字符串）
func (e *RunningHubExecutor) DownloadTextFromURL(textURL string) string {
	resp, err := http.Get(textURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
