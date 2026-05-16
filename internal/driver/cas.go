package driver

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
)

type CASPreviewNamer interface {
	CASPreviewName(ctx context.Context, file model.Obj) (string, error)
}

type CASDownloadRestorer interface {
	CASDownloadRestoreName(ctx context.Context, file model.Obj) (string, error)
}

type ctxKey struct{}

var DstDirPathKey = ctxKey{}

func DstDirPathFromCtx(ctx context.Context) string {
	s, _ := ctx.Value(DstDirPathKey).(string)
	return s
}
