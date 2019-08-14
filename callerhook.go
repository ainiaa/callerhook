package callerhook

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
)

var (
// Used for caller information initialisation
callerInitOnce sync.Once
	logrusPackage string
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)


// CallerHook is a hook to handle wrapper log's caller.
type CallerHook struct {
	PackageName string
	MinimumCallerDepth int
	MaximumCallerDepth int
}

// NewHook returns new CallerHook.
func NewHook(packageName string) *CallerHook {
	hook := &CallerHook{
		PackageName:    packageName,
		MaximumCallerDepth:maximumCallerDepth,
		MinimumCallerDepth:knownLogrusFrames,
	}

	return hook
}

func (hook *CallerHook) SetPackageName(packageName string) *CallerHook{
	hook.PackageName = packageName
	return hook
}

func (hook *CallerHook) SetMinimumCallerDepth(minimumCallerDepth int) *CallerHook{
	hook.MinimumCallerDepth = minimumCallerDepth
	return hook
}

func (hook *CallerHook) SetMaximumCallerDepth(maximumCallerDepth int) *CallerHook{
	hook.MaximumCallerDepth = maximumCallerDepth
	return hook
}

// GetPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func GetPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func (hook *CallerHook)getCaller() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, hook.MaximumCallerDepth)
	depth := runtime.Callers(hook.MinimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[0:depth])

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		if hook.PackageName != "" {
			logrusPackage = hook.PackageName
		} else {
			logrusPackage = GetPackageName(runtime.FuncForPC(pcs[0]).Name())
		}

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary store an entry in a logger interface
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := GetPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// Fire override entry's caller
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Caller = hook.getCaller()
	return nil
}
