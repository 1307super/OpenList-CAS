package local

import (
	"context"
	"os"
	"path/filepath"

	"github.com/OpenListTeam/OpenList/v4/internal/casmeta"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
)

type localCASInfo = casmeta.Info

func newLocalCASHasherWriter() *casmeta.HasherWriter {
	return casmeta.NewHasherWriter(casmeta.DefaultSliceSize)
}

func localIsCASName(name string) bool {
	return casmeta.IsName(name)
}

func (d *Local) shouldUploadCAS(name string) bool {
	return d.GenerateCAS && !localIsCASName(name) && casmeta.GlobalExtAllowed(name)
}

func (d *Local) uploadCAS(ctx context.Context, dstDir model.Obj, info *localCASInfo) error {
	if info == nil || !d.shouldUploadCAS(info.Name) {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	content, err := casmeta.Encode(info)
	if err != nil {
		return err
	}
	casPath := filepath.Join(dstDir.GetPath(), casmeta.FileName(info.Name))
	if err = os.WriteFile(casPath, content, 0o666); err != nil {
		return err
	}
	if d.directoryMap.Has(dstDir.GetPath()) {
		d.directoryMap.UpdateDirSize(dstDir.GetPath())
		d.directoryMap.UpdateDirParents(dstDir.GetPath())
	}
	return nil
}

func (d *Local) deleteSource(ctx context.Context, fullPath string, info *localCASInfo) error {
	if info == nil || !d.DeleteSource || !d.shouldUploadCAS(info.Name) {
		return nil
	}
	return d.Remove(ctx, &model.Object{
		Path: fullPath,
		Name: info.Name,
		Size: info.Size,
	})
}
