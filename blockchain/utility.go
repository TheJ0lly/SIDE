package blockchain

import (
	"github.com/TheJ0lly/GoChain/osspecifics"
	"io/fs"
	"os"
)

// clearFolder - will clear a folder.
func clearFolder(dbLoc string, files []fs.DirEntry) error {
	for _, f := range files {
		fileName := osspecifics.CreatePath(dbLoc, f.Name())
		err := os.Remove(fileName)

		if err != nil {
			return err
		}
	}

	return nil
}
