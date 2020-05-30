package sessions

import (
	"context"
	"io"

	"github.com/escaletech/tog-go/flags"
)

type flagLister interface {
	io.Closer
	ListFlags(context.Context, string) ([]flags.Flag, error)
}
