package batch

import (
	types "colexecdb/pkg/query_engine/a_types"
	vector "colexecdb/pkg/query_engine/b_vector"
	"fmt"
)

func MockBatch(colCnt int, rowCnt int) *Batch {
	bat := NewWithSize(colCnt)
	bat.rowCount = rowCnt

	for i := 0; i < colCnt; i++ {
		bat.Attrs[i] = fmt.Sprintf("%s%d", "mock_", i)

		switch i % 20 {
		case 0:
			bat.Vecs[i] = vector.NewVec(types.T_int32.ToType())
			for j := 0; j < rowCnt; j++ {
				_ = vector.Append[int32](bat.Vecs[i], int32(-j), false)
			}
		case 1:
			bat.Vecs[i] = vector.NewVec(types.T_int64.ToType())
			for j := 0; j < rowCnt; j++ {
				_ = vector.Append[int64](bat.Vecs[i], int64(-j), false)
			}
		}

	}
	return bat
}
