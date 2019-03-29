package runtime

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aegistrace/aegistrace/checkpoint"
	"github.com/aegistrace/aegistrace/journal"
	"github.com/aegistrace/aegistrace/provider"
	"github.com/aegistrace/aegistrace/workflow"
)

type Approver interface {
	Approve(context.Context, workflow.Step, map[string]any) error
}
type ApproverFunc func(context.Context, workflow.Step, map[string]any) error

func (f ApproverFunc) Approve(ctx context.Context, step workflow.Step, input map[string]any) error {
	return f(ctx, step, input)
}

type Observer interface{ Event(string, map[string]any) }
type ObserverFunc func(string, map[string]any)

