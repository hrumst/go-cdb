package cli

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"runtime"
	"testing"
	"time"

	db "github.com/hrumst/go-cdb/internal/database"
	"github.com/hrumst/go-cdb/internal/database/compute/parser"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestIntegrationREPL(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	input := bytes.NewBuffer(make([]byte, 0))
	output := bytes.NewBuffer(make([]byte, 0))

	testDb := db.NewDatabase(
		engine.NewInMemoryEngine(),
		parser.NewCommandExecParserPlain(),
		tools.NewAppLogger(zap.NewNop()),
	)

	input.Write([]byte(fmt.Sprintf("some invalid data \n")))
	input.Write([]byte(fmt.Sprintf("GEtt %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("GET %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("SET %s %s \n", testKey, testVal)))
	input.Write([]byte(fmt.Sprintf("GET %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("DEL %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("GET %s \n", testKey)))

	testRepl := NewREPL("testREPL", 100, testDb, input, output)
	go func() {
		_ = testRepl.Run()
	}()

	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	outb := make([]byte, output.Len())
	_, err := output.Read(outb)
	if err != nil {
		t.Fatal(err)
	}

	expectOutput := `testREPL > testREPL Error: parse error: invalid token, at position: 0, invalid input: 'some invalid data ' 
testREPL > testREPL Error: parse error: invalid token, at position: 0, invalid input: 'GEtt key1 ' 
testREPL > testREPL Error: execution error: inMemoryEngine.Get error: not found key: key1 
testREPL > testREPL Ok:  
testREPL > testREPL Ok: val1 
testREPL > testREPL Ok:  
testREPL > testREPL Error: execution error: inMemoryEngine.Get error: not found key: key1 
testREPL > Bye!
`

	assert.NoError(t, err)
	assert.Equal(t, expectOutput, string(outb))
}

func TestIntegrationREPLWithExit(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	input := bytes.NewBuffer(make([]byte, 0))
	output := bytes.NewBuffer(make([]byte, 0))

	testDb := db.NewDatabase(
		engine.NewInMemoryEngine(),
		parser.NewCommandExecParserPlain(),
		tools.NewAppLogger(zap.NewNop()),
	)

	input.Write([]byte(fmt.Sprintf("some invalid data \n")))
	input.Write([]byte(fmt.Sprintf("GEtt %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("GET %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("SET %s %s \n", testKey, testVal)))
	input.Write([]byte(fmt.Sprintf(" .exit \n")))
	input.Write([]byte(fmt.Sprintf("DEL %s \n", testKey)))
	input.Write([]byte(fmt.Sprintf("GET %s \n", testKey)))

	testRepl := NewREPL("testREPL", 100, testDb, input, output)
	go func() {
		_ = testRepl.Run()
	}()

	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	outb := make([]byte, output.Len())
	_, err := output.Read(outb)
	if err != nil {
		t.Fatal(err)
	}

	expectOutput := `testREPL > testREPL Error: parse error: invalid token, at position: 0, invalid input: 'some invalid data ' 
testREPL > testREPL Error: parse error: invalid token, at position: 0, invalid input: 'GEtt key1 ' 
testREPL > testREPL Error: execution error: inMemoryEngine.Get error: not found key: key1 
testREPL > testREPL Ok:  
testREPL > Bye!
`

	assert.NoError(t, err)
	assert.Equal(t, expectOutput, string(outb))
}
