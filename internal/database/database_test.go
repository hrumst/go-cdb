package database

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestDatabase(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	type testCase struct {
		input            string
		setupStorageMock func(*Mockstorage)
		setupCepMock     func(*MockcommandExecParser)
		expectResult     string
	}

	tcs := []testCase{
		{
			input: "GET " + testKey,
			setupCepMock: func(mcep *MockcommandExecParser) {
				mcep.EXPECT().
					Parse("GET "+testKey).
					Return(
						&compute.CommandExec{
							Command: compute.CommandTypeGet,
							Key:     testKey,
						}, nil,
					)
			},
			setupStorageMock: func(ms *Mockstorage) {
				ms.EXPECT().
					Get(context.Background(), testKey).
					Return("", engine.NotFoundKey).
					Times(1)
			},
			expectResult: formatErrorResult(engine.NotFoundKey),
		}, {
			input: "SET " + testKey + " " + testVal,
			setupCepMock: func(mcep *MockcommandExecParser) {
				mcep.EXPECT().
					Parse("SET "+testKey+" "+testVal).
					Return(
						&compute.CommandExec{
							Command: compute.CommandTypeSet,
							Key:     testKey,
							Val:     testVal,
						}, nil,
					)
			},
			setupStorageMock: func(ms *Mockstorage) {
				ms.EXPECT().
					Set(context.Background(), testKey, testVal).
					Return(nil).
					Times(1)
			},
			expectResult: formatOkResult(""),
		}, {
			input: "GET " + testKey,
			setupCepMock: func(mcep *MockcommandExecParser) {
				mcep.EXPECT().
					Parse("GET "+testKey).
					Return(
						&compute.CommandExec{
							Command: compute.CommandTypeGet,
							Key:     testKey,
						}, nil,
					)
			},
			setupStorageMock: func(ms *Mockstorage) {
				ms.EXPECT().
					Get(context.Background(), testKey).
					Return(testVal, nil).
					Times(1)
			},
			expectResult: formatOkResult(testVal),
		}, {
			input: "DEL " + testKey,
			setupCepMock: func(mcep *MockcommandExecParser) {
				mcep.EXPECT().
					Parse("DEL "+testKey).
					Return(
						&compute.CommandExec{
							Command: compute.CommandTypeDel,
							Key:     testKey,
						}, nil,
					)
			},
			setupStorageMock: func(ms *Mockstorage) {
				ms.EXPECT().
					Del(context.Background(), testKey).
					Return(nil).
					Times(1)
			},
			expectResult: formatOkResult(""),
		}, {
			input: "DEL " + testKey,
			setupCepMock: func(mcep *MockcommandExecParser) {
				mcep.EXPECT().
					Parse("DEL "+testKey).
					Return(
						&compute.CommandExec{
							Command: 1000,
							Key:     testKey,
						}, nil,
					)
			},
			expectResult: formatErrorResult(unsupportedCommandErr) + ": 1000",
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("TestDatabase_case_%d", i),
			func(t *testing.T) {
				ms := NewMockstorage(cntrl)
				if tc.setupStorageMock != nil {
					tc.setupStorageMock(ms)
				}
				mcep := NewMockcommandExecParser(cntrl)
				if tc.setupCepMock != nil {
					tc.setupCepMock(mcep)
				}
				db := NewDatabase(ms, mcep, tools.NewAppLogger(zap.NewNop()))
				result := db.Execute(
					context.Background(),
					tc.input,
				)
				assert.Equal(t, tc.expectResult, result)
			},
		)
	}
}
