package executor

import (
	batch "colexecdb/pkg/query_engine/b_batch"
	"context"
	"fmt"
)

type sqlExecutor struct {
}

// Result exec sql result
type Result struct {
	AffectedRows uint64
	Batches      []*batch.Batch
}

func (s *sqlExecutor) Exec(ctx context.Context, sql string) (Result, error) {
	var res Result
	err := s.ExecTxn(
		ctx,
		func(exec TxnExecutor) error {
			v, err := exec.Exec(sql)
			res = v
			return err
		})
	if err != nil {
		return Result{}, err
	}
	return res, nil
}

func (s *sqlExecutor) ExecTxn(
	ctx context.Context,
	execFunc func(TxnExecutor) error) error {
	exec, err := newTxnExecutor(ctx, s)
	if err != nil {
		return err
	}
	err = execFunc(exec)
	if err != nil {
		fmt.Errorf("internal sql executor error: %v", err)
		return exec.rollback(err)
	}
	if err = exec.commit(); err != nil {
		return err
	}
	return nil
}
