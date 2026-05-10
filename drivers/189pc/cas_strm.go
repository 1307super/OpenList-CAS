package _189pc

import (
	"context"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/cmd/flags"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func (y *Cloud189PC) generateStrm(ctx context.Context, casFileName string) {
	if !y.GenerateSTRM || y.OpenListURL == "" {
		return
	}

	dstDirPath := driver.DstDirPathFromCtx(ctx)
	if dstDirPath == "" {
		return
	}

	mountPath := utils.GetActualMountPath(y.GetStorage().MountPath)
	casCloudPath := stdpath.Join(mountPath, dstDirPath, casFileName)

	strmLocalPath := filepath.Join(flags.DataDir, "strm", strings.TrimPrefix(casCloudPath, "/"))
	strmLocalPath = strings.TrimSuffix(strmLocalPath, ".cas") + ".strm"

	strmContent := strings.TrimRight(y.OpenListURL, "/") + "/d/" + strings.TrimPrefix(casCloudPath, "/")

	strmDir := filepath.Dir(strmLocalPath)
	if err := os.MkdirAll(strmDir, 0755); err != nil {
		log.Errorf("casStrmMkdirError: %s", err)
		return
	}
	if err := os.WriteFile(strmLocalPath, []byte(strmContent), 0644); err != nil {
		log.Errorf("casStrmWriteError: %s", err)
	} else {
		log.Infof("cas_strm: generated %s -> %s", strmLocalPath, strmContent)
	}
}
