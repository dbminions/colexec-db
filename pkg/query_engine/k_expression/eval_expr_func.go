package expression

import (
	types "colexecdb/pkg/query_engine/a_types"
	vector "colexecdb/pkg/query_engine/b_vector"
	batch "colexecdb/pkg/query_engine/c_batch"
	process "colexecdb/pkg/query_engine/e_process"
)

type FunctionExpressionExecutor struct {
	ResultVector *vector.Vector
	// parameters related
	parameterResults  []*vector.Vector
	parameterExecutor []ExpressionExecutor

	evalFn func(
		params []*vector.Vector,
		result *vector.Vector,
		proc *process.Process,
		length int) error
}

func (expr *FunctionExpressionExecutor) Init(
	_ *process.Process,
	parameterNum int,
	retType types.Type,
	fn func(
		params []*vector.Vector,
		result *vector.Vector,
		proc *process.Process,
		length int) error,
) (err error) {

	expr.evalFn = fn
	expr.parameterResults = make([]*vector.Vector, parameterNum)
	expr.parameterExecutor = make([]ExpressionExecutor, parameterNum)

	expr.ResultVector = vector.NewVec(retType)
	return err
}

func (expr *FunctionExpressionExecutor) Eval(proc *process.Process, batches []*batch.Batch) (*vector.Vector, error) {
	var err error
	for i := range expr.parameterExecutor {
		expr.parameterResults[i], err = expr.parameterExecutor[i].Eval(proc, batches)
		if err != nil {
			return nil, err
		}
	}

	if err = expr.evalFn(expr.parameterResults, expr.ResultVector, proc, batches[0].GetRowCount()); err != nil {
		return nil, err
	}
	return expr.ResultVector, nil
}

func (expr *FunctionExpressionExecutor) Free() {
	for _, p := range expr.parameterExecutor {
		p.Free()
	}
}

func (expr *FunctionExpressionExecutor) SetParameter(index int, executor ExpressionExecutor) {
	expr.parameterExecutor[index] = executor
}
