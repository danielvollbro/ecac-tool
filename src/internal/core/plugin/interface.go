package plugin

import "context"

type HCLPlugin interface {
	Schema() any
	Validate(ctx context.Context) error
	Run(ctx context.Context) (string, error)
}
