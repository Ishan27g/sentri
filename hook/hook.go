package hook

import (
	"context"
	"io"

	"github.com/DavidGamba/dgtools/run"
	"github.com/Ishan27g/sentri/internal"
)

func Hook(appName string, w io.Writer) io.Writer {
	pr, pw := io.Pipe()
	internal.InitPool(context.Background())
	internal.Ready(appName, pr)
	w = io.MultiWriter(w, pw)
	return w
}
func Cmd(appName string, w io.Writer, run *run.RunInfo) {
	w = Hook(appName, w)
	_ = run.Run(w)
}
