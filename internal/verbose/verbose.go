package verbose

import (
	"fmt"
	"io"
	"time"
)

const timestampLayout = "2006-01-02 15:04:05"

func PrefixAt(t time.Time, model string) string {
	return fmt.Sprintf("%s %s", t.Format(timestampLayout), model)
}

func Printf(w io.Writer, now func() time.Time, model, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "%s %s\n", PrefixAt(now(), model), fmt.Sprintf(format, args...))
}
