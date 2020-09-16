// +build !linux

package jsengine

import (
	"fmt"
	"rogchap.com/v8go"
	"strings"
)

func RunJSFunc(jsContext, jsFunc string, args ...string) (string, error) {
	ctx, _ := v8go.NewContext(nil)
	ctx.RunScript(jsContext, "test.js")
	strArgs := []string{}
	for _, arg := range args {
		newFormat := fmt.Sprintf("'%s'", arg)
		strArgs = append(strArgs, newFormat)
	}
	callfunc := jsFunc + "(" + strings.Join(strArgs, ",") + ")"
	result, err := ctx.RunScript(callfunc, "test.js")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return result.String(), err
}
