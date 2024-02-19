package actions

import "context"

// Executor is an interface to an action executor
type Executor interface {
	// Run will execute an action with optional arguments
	Run(context.Context, map[string]interface{}) error
}
