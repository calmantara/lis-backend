package mysql

import (
	"context"

	"github.com/Calmantara/lis-backend/internal/core/ports"
)

type transactionImpl struct {
	mysql MySql
}

func NewTransaction(mysql MySql) ports.Transaction {
	return &transactionImpl{mysql: mysql}
}

func (t *transactionImpl) Begin(ctx context.Context) context.Context {
	txn := t.mysql.Master().Begin()
	ctx, _ = begin(ctx, txn)

	return ctx
}

func (t *transactionImpl) End(ctx context.Context, err *error) error {
	if err != nil && *err != nil {
		return rollback(ctx)
	}

	return commit(ctx)

}
