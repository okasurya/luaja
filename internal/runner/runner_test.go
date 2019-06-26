package runner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Person struct {
	Name string
}

func TestPlainLuaHelloWorld(t *testing.T) {
	ctx := context.Background()
	script := `
	print('hello world')
	`
	err := RunScript(ctx, script, nil, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestOutputLua(t *testing.T) {
	ctx := context.Background()
	script := `
		return {
			name= "John Doe"
		}
	`
	var output Person
	err := RunScript(ctx, script, nil, &output)
	if err != nil {
		t.Error(err)
	}

	if output.Name != "John Doe" {
		t.Errorf("failed, expected %s actual %s", "John Doe", output.Name)
	}

}

func TestInputLua(t *testing.T) {
	ctx := context.Background()
	script := `
		return input
	`
	input := Person{Name: "John Doe"}
	var output Person
	err := RunScript(ctx, script, input, &output)
	if err != nil {
		t.Error(err)
	}
	if output.Name != input.Name {
		t.Error("failed, input doesn't have same output")
	}

}

func TestInputDiffTypeFromOutput(t *testing.T) {
	ctx := context.Background()
	script := `
		return input
	`
	input := Person{Name: "John Doe"}
	var output string
	err := RunScript(ctx, script, input, &output)
	require.NotNil(t, err)
}

func TestWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	script := `
		while true do
		end
	`
	err := RunScript(ctx, script, nil, nil)
	require.NotNil(t, err)
}
