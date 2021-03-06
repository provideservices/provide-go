package baseline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/provideplatform/provide-go/api"
	"github.com/provideplatform/provide-go/common"
)

const defaultBaselineHost = "localhost:8080"
const defaultBaselinePath = "api/v1"
const defaultBaselineScheme = "http"

// Service for the baseline api
type Service struct {
	api.Client
}

// InitBaselineService convenience method to initialize a `baseline.Service` instance
func InitBaselineService(token string) *Service {
	host := defaultBaselineHost
	if os.Getenv("BASELINE_API_HOST") != "" {
		host = os.Getenv("BASELINE_API_HOST")
	}

	path := defaultBaselinePath
	if os.Getenv("BASELIEN_API_PATH") != "" {
		host = os.Getenv("BASELIEN_API_PATH")
	}

	scheme := defaultBaselineScheme
	if os.Getenv("BASELINE_API_SCHEME") != "" {
		scheme = os.Getenv("BASELINE_API_SCHEME")
	}

	return &Service{
		api.Client{
			Host:   host,
			Path:   path,
			Scheme: scheme,
			Token:  common.StringOrNil(token),
		},
	}
}

// ConfigureStack updates the global configuration on the local baseline stack
func ConfigureStack(token string, params map[string]interface{}) error {
	status, _, err := InitBaselineService(token).Put("config", params)
	if err != nil {
		return fmt.Errorf("failed to configure baseline stack; status: %v; %s", status, err.Error())
	}

	if status != 204 {
		return fmt.Errorf("failed to configure baseline stack; status: %v", status)
	}

	return nil
}

// ListWorkgroups retrieves a paginated list of baseline workgroups scoped to the given API token
func ListWorkgroups(token, applicationID string, params map[string]interface{}) ([]*Workgroup, error) {
	status, resp, err := InitBaselineService(token).Get("workgroups", params)
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to list baseline workgroups; status: %v", status)
	}

	workgroups := make([]*Workgroup, 0)
	for _, item := range resp.([]interface{}) {
		workgroup := &Workgroup{}
		workgroupraw, _ := json.Marshal(item)
		json.Unmarshal(workgroupraw, &workgroup)
		workgroups = append(workgroups, workgroup)
	}

	return workgroups, nil
}

// CreateWorkgroup initializes a new or previously-joined workgroup on the local baseline stack
func CreateWorkgroup(token string, params map[string]interface{}) (*Workgroup, error) {
	status, resp, err := InitBaselineService(token).Post("workgroups", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create workgroup; status: %v; %s", status, err.Error())
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to create workgroup; status: %v", status)
	}

	workgroup := &Workgroup{}
	workgroupraw, _ := json.Marshal(resp)
	err = json.Unmarshal(workgroupraw, &workgroup)

	return workgroup, nil
}

// UpdateWorkgroup updates a previously-initialized workgroup on the local baseline stack
func UpdateWorkgroup(id, token string, params map[string]interface{}) error {
	uri := fmt.Sprintf("workgroups/%s", id)
	status, _, err := InitBaselineService(token).Post(uri, params)
	if err != nil {
		return fmt.Errorf("failed to update workgroup; status: %v; %s", status, err.Error())
	}

	if status != 204 {
		return fmt.Errorf("failed to update workgroup; status: %v", status)
	}

	return nil
}

// ListWorkflows retrieves a paginated list of baseline workflows scoped to the given API token
func ListWorkflows(token, applicationID string, params map[string]interface{}) ([]*Workflow, error) {
	status, resp, err := InitBaselineService(token).Get("workflows", params)
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to list baseline workflows; status: %v", status)
	}

	workflows := make([]*Workflow, 0)
	for _, item := range resp.([]interface{}) {
		workflow := &Workflow{}
		workflowraw, _ := json.Marshal(item)
		json.Unmarshal(workflowraw, &workflow)
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// CreateWorkflow initializes a new workflow on the local baseline stack
func CreateWorkflow(token string, params map[string]interface{}) (*Workflow, error) {
	status, resp, err := InitBaselineService(token).Post("workflows", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow; status: %v; %s", status, err.Error())
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to create workflow; status: %v", status)
	}

	workflow := &Workflow{}
	workflowraw, _ := json.Marshal(resp)
	err = json.Unmarshal(workflowraw, &workflow)

	return workflow, nil
}

// ListWorksteps retrieves a paginated list of baseline worksteps scoped to the given API token
func ListWorksteps(token, applicationID string, params map[string]interface{}) ([]*Workstep, error) {
	status, resp, err := InitBaselineService(token).Get("worksteps", params)
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to list baseline worksteps; status: %v", status)
	}

	worksteps := make([]*Workstep, 0)
	for _, item := range resp.([]interface{}) {
		workstep := &Workstep{}
		workstepraw, _ := json.Marshal(item)
		json.Unmarshal(workstepraw, &workstep)
		worksteps = append(worksteps, workstep)
	}

	return worksteps, nil
}

// CreateWorkstep initializes a new workstep on the local baseline stack
func CreateWorkstep(token string, params map[string]interface{}) (*Workstep, error) {
	status, resp, err := InitBaselineService(token).Post("worksteps", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create workstep; status: %v; %s", status, err.Error())
	}

	if status != 200 {
		return nil, fmt.Errorf("failed to create workstep; status: %v", status)
	}

	workstep := &Workstep{}
	workstepraw, _ := json.Marshal(resp)
	err = json.Unmarshal(workstepraw, &workstep)

	return workstep, nil
}

// CreateObject is a generic way to baseline a business object
func CreateObject(token string, params map[string]interface{}) (interface{}, error) {
	status, resp, err := InitBaselineService(token).Post("objects", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create baseline object; status: %v; %s", status, err.Error())
	}

	if status != 202 {
		return nil, fmt.Errorf("failed to create baseline object; status: %v", status)
	}

	return resp, nil
}

// UpdateObject updates a business object
func UpdateObject(token, id string, params map[string]interface{}) error {
	uri := fmt.Sprintf("objects/%s", id)
	status, _, err := InitBaselineService(token).Put(uri, params)
	if err != nil {
		return fmt.Errorf("failed to update baseline state; status: %v; %s", status, err.Error())
	}

	if status != 202 {
		return fmt.Errorf("failed to update baseline state; status: %v", status)
	}

	return nil
}
