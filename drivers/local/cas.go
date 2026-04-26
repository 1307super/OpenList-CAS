package local

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
)

type localCASPayload struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	MD5        string `json:"md5"`
	SliceMD5   string `json:"sliceMd5"`
	CreateTime string `json:"create_time"`
}

func localIsCASName(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".cas")
}

func (d *Local) uploadCAS(ctx context.Context, dstDir model.Obj, srcObj model.Obj) error {
	if !d.GenerateCAS || srcObj == nil || localIsCASName(srcObj.GetName()) {
		return nil
	}
	file, err := os.Open(srcObj.GetPath())
	if err != nil {
		return err
	}
	defer file.Close()
	md5, err := utils.HashFile(utils.MD5, file)
	if err != nil {
		return err
	}
	md5 = strings.ToUpper(md5)
	content, err := utils.Json.Marshal(localCASPayload{
		Name:       srcObj.GetName(),
		Size:       srcObj.GetSize(),
		MD5:        md5,
		SliceMD5:   md5,
		CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err != nil {
		return err
	}
	content = []byte(base64.StdEncoding.EncodeToString(content))
	casPath := filepath.Join(dstDir.GetPath(), srcObj.GetName()+".cas")
	if err = os.WriteFile(casPath, content, 0o666); err != nil {
		return err
	}
	if d.DeleteSource {
		return d.Remove(ctx, srcObj)
	}
	if d.directoryMap.Has(dstDir.GetPath()) {
		d.directoryMap.UpdateDirSize(dstDir.GetPath())
		d.directoryMap.UpdateDirParents(dstDir.GetPath())
	}
	return nil
}
