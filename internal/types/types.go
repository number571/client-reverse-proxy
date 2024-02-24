package types

import "context"

type IRunner interface {
	Run(ctx context.Context) error
}
