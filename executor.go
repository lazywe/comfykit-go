package comfykit

import (
	"github.com/lazywe/comfykit-go/comfyui"
)

type Executor interface {
	ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*comfyui.ExecuteResult, error)
}
