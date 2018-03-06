package watch

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// ScanTime is time duration for scanning files
var ScanTime = time.Millisecond * 500

type (
	// Event is event of file changes
	Event struct {
		File string `json:"file"`
		Op   string `json:"op"`
	}
	// EventHandler is handler to resolve events
	EventHandler func([]*Event)
)

// Watch watch dirs to compare and emmit events
func Watch(dirs []string, fn EventHandler) {
	log.Info().Strs("dirs", dirs).Float64("freq", ScanTime.Seconds()).Msg("watch-dirs")

	ticker := time.NewTicker(ScanTime)
	defer ticker.Stop()

	var lastDirs map[string]int64

	for {
		nowDirs := scanDirs(dirs)
		if len(lastDirs) > 0 {
			events := compareDirs(lastDirs, nowDirs)
			if len(events) > 0 && fn != nil {
				log.Info().Strs("dirs", dirs).Interface("changes", events).Msg("watch-dirs-ok")
				go fn(events)
			}
		}
		lastDirs = nowDirs
		// log.Debug().Strs("dirs", dirs).Msg("scan-dirs-ok")
		<-ticker.C
	}
}

func compareDirs(lastDirs, nowDirs map[string]int64) []*Event {
	var events []*Event
	for nowFile, nowT := range nowDirs {
		lastT := lastDirs[nowFile]
		if lastT == 0 {
			events = append(events, &Event{
				File: nowFile,
				Op:   "created",
			})
			continue
		}
		if nowT != lastT {
			events = append(events, &Event{
				File: nowFile,
				Op:   "modified",
			})
			continue
		}
	}
	for lastFile := range lastDirs {
		nowT := nowDirs[lastFile]
		if nowT == 0 {
			events = append(events, &Event{
				File: lastFile,
				Op:   "deleted",
			})
		}
	}
	return events
}

func scanDirs(dirs []string) map[string]int64 {
	scanFiles := make(map[string]int64)
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			scanFiles[fpath] = info.ModTime().Unix()
			return nil
		})
		if err != nil {
			log.Warn().Err(err).Str("dir", dir).Msg("watch-dir-error")
		}
	}
	return scanFiles
}
