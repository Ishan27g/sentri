package hook

import (
	"context"
	"io"

	"github.com/Ishan27g/sentri/internal"
)

func Hook(appName string, w io.Writer) io.Writer {
	pr, pw := io.Pipe()
	internal.InitPool(context.Background())
	internal.Ready(appName, pr)
	w = io.MultiWriter(w, pw)
	return w
}
