package runner

import (
	"context"
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

func checkIfTypeSame(result interface{}, output interface{}) error {
	if reflect.TypeOf(output).Elem().Name() != reflect.TypeOf(result).Name() {
		return fmt.Errorf(
			"type of output is %s, but the result is %s",
			reflect.TypeOf(output).Elem().Name(),
			reflect.TypeOf(result).Name(),
		)
	}
	return nil
}

func convertUserData(userData *lua.LUserData, output interface{}) error {
	if err := checkIfTypeSame(userData.Value, output); err != nil {
		return err
	}
	reflect.ValueOf(output).Elem().Set(reflect.ValueOf(userData.Value))
	return nil
}

func convertString(luastring lua.LString, output interface{}) error {
	if err := checkIfTypeSame(string(luastring), output); err != nil {
		return err
	}
	reflect.ValueOf(output).Elem().Set(reflect.ValueOf(string(luastring)))
	return nil
}

func convertNumber(luanumber lua.LNumber, output interface{}) error {
	if err := checkIfTypeSame(float64(luanumber), output); err != nil {
		return err
	}
	reflect.ValueOf(output).Elem().Set(reflect.ValueOf(float64(luanumber)))
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
			return convertString(lv.(lua.LString), output)
		case lua.LTNumber:
			return convertNumber(lv.(lua.LNumber), output)
		case lua.LTUserData:
			return convertUserData(lv.(*lua.LUserData), output)
		default:
			return fmt.Errorf("output type %s not supported", lv.Type())
		}

	}
	return nil
}
