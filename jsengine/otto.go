// +build !darwin !darwin_amd64

package jsengine

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"strings"
)

func RunJSFunc(jsContext, jsFunc string, args ...string) (string, error) {
	vm := otto.New()
	strArgs := []string{}
	for _, arg := range args {
		newFormat := fmt.Sprintf("'%s'", arg)
		strArgs = append(strArgs, newFormat)
	}
	callfunc := "result = " + jsFunc + "(" + strings.Join(strArgs, ",") + ")"
	vm.Run(fmt.Sprintf("%s\n%s", jsContext, callfunc))
	value, err := vm.Get("result")
	if err != nil {
		return "", err
	}
	return value.ToString()
}
