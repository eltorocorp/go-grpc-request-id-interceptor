package panichandler

import (
	"fmt"
	"os"
	"runtime/debug"

	"golang.org/x/net/context"
)

var _ PanicHandler = LogPanicDump

func LogPanicDump(ctx context.Context, r any) {
	_, _ = fmt.Fprint(os.Stderr, string(debug.Stack()))
}
