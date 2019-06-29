package runner

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func convertTable(table *lua.LTable, output interface{}) error {
	if err := gluamapper.Map(table, &output); err != nil {
		return err
	}
	return nil
}

func setValueToOutput(value interface{}, output interface{}) error {
	if reflect.ValueOf(output).Type().Kind() != reflect.Ptr {
		return errors.New("output is not a pointer")
	}
	if reflect.TypeOf(output).Elem().Name() != reflect.TypeOf(value).Name() {
		return fmt.Errorf(
			"type of output is %s, but the result is %s",
			reflect.TypeOf(output).Elem().Name(),
			reflect.TypeOf(value).Name(),
		)
	}
	reflect.ValueOf(output).Elem().Set(reflect.ValueOf(value))
	return nil
}

func RunScript(context context.Context, script string, input interface{}, output interface{}) error {
	L := lua.NewState()
	defer L.Close()
	L.SetContext(context)
	L.SetGlobal("input", luar.New(L, input))
	if err := L.DoString(script); err != nil {
		return err
	}

	if output != nil {
		lv := L.Get(-1)
		switch lv.Type() {
		case lua.LTTable:
			return convertTable(lv.(*lua.LTable), output)
		case lua.LTString:
			return setValueToOutput(string(lv.(lua.LString)), output)
		case lua.LTNumber:
			return setValueToOutput(float64(lv.(lua.LNumber)), output)
		case lua.LTBool:
			return setValueToOutput(bool(lv.(lua.LBool)), output)
		case lua.LTUserData:
			return setValueToOutput(lv.(*lua.LUserData).Value, output)
		default:
			return fmt.Errorf("output type %s not supported", lv.Type())
		}

	}
	return nil
}
