package packer

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// Asset writes bytes data to one file with base64 encoding
func Asset(toFile string, bytesData map[string][]byte) (int, error) {
	fileContents := `package asset
	
var Files = make(map[string]string)
var FilesTime = "` + time.Now().Format("2006/01/02 15:04:05") + `"` + "\n" + `

func init(){` + "\n"
	for file, b := range bytesData {
		fileContents += "\t" + `Files["` + file + `"] = "` + base64Encode(b) + `"` + "\n\n"
	}
	fileContents += `
}
	`
	os.MkdirAll(filepath.Dir(toFile), 0755)
	return len(fileContents), ioutil.WriteFile(toFile, []byte(fileContents), 0644)
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

// Unasset unpack files from source code
func Unasset(fileStrings map[string]string) error {
	for file, data := range fileStrings {
		gzipBytes, err := base64Decode(data)
		if err != nil {
			log.Warn().Str("file", file).Err(err).Msg("read-file-error")
			continue
		}
		rawBytes, err := Ungzip(gzipBytes)
		if err != nil {
			log.Warn().Str("file", file).Err(err).Msg("ungzip-file-error")
			continue
		}
		os.MkdirAll(filepath.Dir(file), 0644)
		if err := ioutil.WriteFile(file, rawBytes, 0644); err != nil {
			log.Warn().Str("file", file).Err(err).Msg("write-file-error")
			continue
		}
		log.Debug().Str("file", file).Msg("load-file")
	}
	return nil
}
