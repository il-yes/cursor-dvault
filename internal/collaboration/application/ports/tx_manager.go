package collaboration_ports

import "context"

type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type NopTransactionManager struct{}

func (n *NopTransactionManager) ExecuteInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}
