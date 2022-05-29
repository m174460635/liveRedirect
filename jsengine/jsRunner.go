package jsengine

import (
	"fmt"
	"github.com/dop251/goja"
	//"rogchap.com/v8go"
	"strings"
)

//func RunJSFunc(jsContext, jsFunc string, args ...string) (string, error) {
//	ctx := v8go.NewContext()
//	ctx.RunScript(jsContext, "test.js")
//	strArgs := []string{}
//	for _, arg := range args {
//		newFormat := fmt.Sprintf("'%s'", arg)
//		strArgs = append(strArgs, newFormat)
//	}
//	callfunc := jsFunc + "(" + strings.Join(strArgs, ",") + ")"
//	result, err := ctx.RunScript(callfunc, "test.js")
//	if err != nil {
//		fmt.Println(err.Error())
//		return "", err
//	}
//	return result.String(), err
//}
func RunJSFunc2(jsContext, jsFunc string, args ...string) (string, error) {
	vm := goja.New()
	_, err := vm.RunString(jsContext)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	strArgs := []string{}
	for _, arg := range args {
		newFormat := fmt.Sprintf("'%s'", arg)
		strArgs = append(strArgs, newFormat)
	}
	callfunc := jsFunc + "(" + strings.Join(strArgs, ",") + ")"
	result, err := vm.RunString(callfunc)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return result.String(), err
}
