package runner

import (
	"context"
	"errors"
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

func convertUserData(userData *lua.LUserData, output interface{}) error {
	if reflect.TypeOf(output).Elem().Name() != reflect.TypeOf(userData.Value).Name() {
		return errors.New("type of output is not same")
	}
	reflect.ValueOf(output).Elem().Set(reflect.ValueOf(userData.Value))
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
		case lua.LTUserData:
			return convertUserData(lv.(*lua.LUserData), output)
		}

	}
	return nil
}
