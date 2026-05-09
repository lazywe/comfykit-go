package comfykit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lazywe/comfykit-go/comfyui"
)

const (
	defaultComfyUIBaseURL    = "http://127.0.0.1:8188"
	defaultRunningHubBaseURL = "https://www.runninghub.ai"
	defaultRunningHubTimeout = 0
	defaultRunningHubRetry   = 3
)

type ExecutorType string

const (
	ExecutorTypeHTTP      ExecutorType = "http"
	ExecutorTypeWebSocket ExecutorType = "websocket"
)

type ComfyKit struct {
	// Local ComfyUI configuration
	comfyUIBaseURL    string
	executorType      ExecutorType
	apiKey            string
	cookies           string
	httpExecutor      *comfyui.HTTPExecutor
	websocketExecutor *comfyui.WebSocketExecutor

	// RunningHub configuration
	runningHubBaseURL  string
	runningHubAPIKey   string
	runningHubTimeout  int
	runningHubRetry    int
	runningHubInstance string
	runningHubExecutor *comfyui.RunningHubExecutor
}

type ComfyKitOption func(*ComfyKit)

func WithComfyUIBaseURL(url string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.comfyUIBaseURL = url
	}
}

func WithExecutorType(t ExecutorType) ComfyKitOption {
	return func(k *ComfyKit) {
		k.executorType = t
	}
}

func WithAPIKey(key string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.apiKey = key
	}
}

func WithCookies(cookies string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.cookies = cookies
	}
}

func WithRunningHubBaseURL(url string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.runningHubBaseURL = url
	}
}

func WithRunningHubAPIKey(key string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.runningHubAPIKey = key
	}
}

func WithRunningHubTimeout(timeout int) ComfyKitOption {
	return func(k *ComfyKit) {
		k.runningHubTimeout = timeout
	}
}

func WithRunningHubRetry(retry int) ComfyKitOption {
	return func(k *ComfyKit) {
		k.runningHubRetry = retry
	}
}

func WithRunningHubInstance(instance string) ComfyKitOption {
	return func(k *ComfyKit) {
		k.runningHubInstance = instance
	}
}

func NewComfyKit(opts ...ComfyKitOption) *ComfyKit {
	k := &ComfyKit{
		comfyUIBaseURL:    getEnvOrDefault("COMFYUI_BASE_URL", defaultComfyUIBaseURL),
		executorType:      ExecutorType(getEnvOrDefault("COMFYUI_EXECUTOR_TYPE", "http")),
		apiKey:            os.Getenv("COMFYUI_API_KEY"),
		cookies:           os.Getenv("COMFYUI_COOKIES"),
		runningHubBaseURL: getEnvOrDefault("RUNNINGHUB_BASE_URL", defaultRunningHubBaseURL),
		runningHubTimeout: defaultRunningHubTimeout,
		runningHubRetry:   defaultRunningHubRetry,
	}

	for _, opt := range opts {
		opt(k)
	}

	if k.runningHubAPIKey == "" {
		k.runningHubAPIKey = k.apiKey
	}

	if envTimeout := os.Getenv("RUNNINGHUB_TIMEOUT"); envTimeout != "" {
		k.runningHubTimeout = parseIntOrDefault(envTimeout, defaultRunningHubTimeout)
	}
	if envRetry := os.Getenv("RUNNINGHUB_RETRY_COUNT"); envRetry != "" {
		k.runningHubRetry = parseIntOrDefault(envRetry, defaultRunningHubRetry)
	}
	if envInstance := os.Getenv("RUNNINGHUB_INSTANCE_TYPE"); envInstance != "" {
		k.runningHubInstance = envInstance
	}

	k.executorType = ExecutorType(strings.ToLower(string(k.executorType)))

	return k
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseIntOrDefault(s string, def int) int {
	var v int
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		return def
	}
	return v
}

func (k *ComfyKit) Close() error {
	var err error

	if k.runningHubExecutor != nil {
		if e := k.runningHubExecutor.Close(); e != nil {
			err = e
		}
	}

	if k.httpExecutor != nil {
		k.httpExecutor.Close()
	}

	if k.websocketExecutor != nil {
		k.websocketExecutor.Close()
	}

	return err
}

func (k *ComfyKit) GetComfyUIBaseURL() string {
	return k.comfyUIBaseURL
}

func (k *ComfyKit) GetExecutorType() ExecutorType {
	return k.executorType
}

func (k *ComfyKit) GetAPIKey() string {
	return k.apiKey
}

func (k *ComfyKit) GetRunningHubBaseURL() string {
	return k.runningHubBaseURL
}

func (k *ComfyKit) GetRunningHubAPIKey() string {
	return k.runningHubAPIKey
}

func (k *ComfyKit) GetRunningHubTimeout() int {
	return k.runningHubTimeout
}

func (k *ComfyKit) getHttpExecutor() *comfyui.HTTPExecutor {
	if k.httpExecutor == nil {
		k.httpExecutor = comfyui.NewHTTPExecutor(k.comfyUIBaseURL, k.apiKey, k.cookies)
	}
	return k.httpExecutor
}

func (k *ComfyKit) getWebsocketExecutor() *comfyui.WebSocketExecutor {
	if k.websocketExecutor == nil {
		k.websocketExecutor = comfyui.NewWebSocketExecutor(k.comfyUIBaseURL, k.apiKey, k.cookies)
	}
	return k.websocketExecutor
}

func (k *ComfyKit) getRunninghubExecutor() *comfyui.RunningHubExecutor {
	if k.runningHubExecutor == nil {
		k.runningHubExecutor = comfyui.NewRunningHubExecutor(
			k.runningHubBaseURL,
			k.runningHubAPIKey,
			k.runningHubTimeout,
			k.runningHubRetry,
			k.runningHubInstance,
		)
	}
	return k.runningHubExecutor
}

func (k *ComfyKit) getLocalExecutor() comfyui.Executor {
	if k.executorType == ExecutorTypeWebSocket {
		return k.getWebsocketExecutor()
	}
	return k.getHttpExecutor()
}

func (k *ComfyKit) isRunninghubWorkflowId(workflow string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(workflow)
}

func (k *ComfyKit) isUrl(workflow string) bool {
	return strings.HasPrefix(workflow, "http://") || strings.HasPrefix(workflow, "https://")
}

func (k *ComfyKit) isFilePath(workflow string) bool {
	if _, err := os.Stat(workflow); err == nil {
		return true
	}
	return strings.Contains(workflow, "/") || strings.Contains(workflow, "\\")
}

func (k *ComfyKit) Execute(workflow string, params map[string]interface{}) (*comfyui.ExecuteResult, error) {
	if k.isRunninghubWorkflowId(workflow) {
		return k.getRunninghubExecutor().ExecuteByID(workflow, params)
	}

	if k.isUrl(workflow) {
		tempPath, err := downloadFile(workflow)
		if err != nil {
			return nil, fmt.Errorf("failed to download workflow: %w", err)
		}
		defer os.Remove(tempPath)
		return k.getLocalExecutor().ExecuteWorkflow(tempPath, params)
	}

	if k.isFilePath(workflow) {
		if comfyui.IsRunningHubWorkflow(workflow) {
			return k.getRunninghubExecutor().ExecuteWorkflow(workflow, params)
		}
		return k.getLocalExecutor().ExecuteWorkflow(workflow, params)
	}

	return k.getLocalExecutor().ExecuteWorkflow(workflow, params)
}

func (k *ComfyKit) ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*comfyui.ExecuteResult, error) {
	tempFile, err := os.CreateTemp("", "workflow_*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if err := json.NewEncoder(tempFile).Encode(workflowJSON); err != nil {
		return nil, fmt.Errorf("failed to write workflow JSON: %w", err)
	}
	tempFile.Close()

	return k.Execute(tempPath, params)
}

// ExecuteAsyncByID 异步创建 RunningHub 任务
// 参数：
//   workflowID - RunningHub 工作流 ID
//   params - 用户提供的参数映射
// 返回：任务 ID、输出节点映射、错误信息
func (k *ComfyKit) ExecuteAsyncByID(workflowID string, params map[string]interface{}) (string, map[string]string, error) {
	return k.getRunninghubExecutor().ExecuteAsyncByID(workflowID, params)
}

// GetTaskCompletion 检查任务状态并在完成时获取结果
// 参数：
//   taskID - 任务 ID
//   outputID2Var - 输出节点 ID 到变量名的映射
// 返回：执行结果、是否完成、错误信息
func (k *ComfyKit) GetTaskCompletion(taskID string, outputID2Var map[string]string) (*comfyui.ExecuteResult, bool, error) {
	return k.getRunninghubExecutor().GetTaskCompletion(taskID, outputID2Var)
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "workflow_*.json")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

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

func init() {
	time.Local = time.FixedZone("UTC", 0)
}
