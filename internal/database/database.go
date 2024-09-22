package database

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	pc "github.com/hrumst/go-cdb/internal/database/compute"
)

var (
	unsupportedCommandErr = errors.New("unsupported command in request")
)

//go:generate mockgen -source=database.go -package=database -destination=mock.go
type commandExecParser interface {
	Parse(input string) (*pc.CommandExec, error)
}

type storage interface {
	Set(ctx context.Context, key string, val string) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
}

type logger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, fields ...zapcore.Field)
}

type database struct {
	storage storage
	cep     commandExecParser
	logger  logger
}

func NewDatabase(
	storage storage,
	cep commandExecParser,
	logger logger,
) *database {
	return &database{
		storage: storage,
		cep:     cep,
		logger:  logger,
	}
}

func formatErrorResult(err error) string {
	return "Error: " + err.Error()
}

func formatOkResult(res string) string {
	if len(res) == 0 {
		return "Ok: <empty>"
	}
	return "Ok: " + res
}

func (d *database) Execute(ctx context.Context, input string) string {
	d.logger.Debug(
		ctx,
		"incoming db query",
		zap.String("input", input),
	)

	cmdExec, err := d.cep.Parse(input)
	if err != nil {
		d.logger.Debug(
			ctx,
			"db query error",
			zap.Any("error", err),
		)

		return formatErrorResult(err)
	}

	d.logger.Debug(
		ctx,
		"db query parse",
		zap.Any("query", cmdExec),
	)

	var (
		resultVal string
		execErr   error
	)
	switch cmdExec.Command {
	case pc.CommandTypeGet:
		resultVal, execErr = d.storage.Get(ctx, cmdExec.Key)
	case pc.CommandTypeSet:
		execErr = d.storage.Set(ctx, cmdExec.Key, cmdExec.Val)
	case pc.CommandTypeDel:
		execErr = d.storage.Del(ctx, cmdExec.Key)
	default:
		execErr = fmt.Errorf("%w: %d", unsupportedCommandErr, cmdExec.Command)
	}

	if execErr != nil {
		d.logger.Debug(
			ctx,
			"db query error",
			zap.Any("error", execErr),
		)
		return formatErrorResult(execErr)
	}
	d.logger.Debug(
		ctx,
		"db query result",
		zap.Any("result", resultVal),
	)
	return formatOkResult(resultVal)
}
