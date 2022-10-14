package laclong

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GetFileName(contestantId string) string {
	fileName := fmt.Sprintf("%s_ticks.csv", strings.Replace(contestantId, "/", "_", -1))
	fileName = filepath.Join("/data", fileName)
	return fileName
}
