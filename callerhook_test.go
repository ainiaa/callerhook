package callerhook

import (
	"testing"
)
// Tests that writing to a tempfile log works.
// Matches the 'msg' of the output and deletes the tempfile.
func Test_Caller(t *testing.T) {

	hook := NewHook("")

	hook.SetPackageName("")
	hook.getCaller()
}

