package callerhook

import (
	"fmt"
	"runtime"
	"testing"
)
// Tests that writing to a tempfile log works.
// Matches the 'msg' of the output and deletes the tempfile.
func Test_Caller(t *testing.T) {

	var pkg string
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(0, pcs)
	//frames := runtime.CallersFrames(pcs[1:depth]) //psc[0]为runtime 需要过滤掉

	pkg = GetPackageName(runtime.FuncForPC(pcs[1]).Name())
	fmt.Printf("pkg:%s cnt:%d\n", pkg, depth)
	for i:=0;i < depth;i++ {
		pkg2 := GetPackageName(runtime.FuncForPC(pcs[i]).Name())
		fmt.Printf("index:%d pkg:%s cnt:%d\n",i, pkg2, depth)
	}

	hook := NewHook("")

	hook.SetPackageName("")
	hook.getCaller()
}

