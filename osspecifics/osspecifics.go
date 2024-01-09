package osspecifics

import (
	"github.com/TheJ0lly/GoChain/generalerrors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

var PathSep string

func init() {
	if runtime.GOOS == "windows" {
		PathSep = "\\"
	} else {
		PathSep = "/"
	}
}

func CreatePath(format ...string) string {
	var stringToReturn string

	for i := 0; i < len(format); i++ {
		stringToReturn += format[i]

		if i != len(format)-1 {
			stringToReturn += PathSep
		}
	}
	return stringToReturn
}

func GetFileName(filepath string) string {

	//_, filename, ok := strings.Cut(filepath, PathSep)

	var filename []byte

	for i := len(filepath) - 1; i >= 0; i-- {
		if filepath[i] == PathSep[0] {
			break
		} else {
			filename = slices.Insert(filename, 0, filepath[i])
		}
	}

	return string(filename)
}

func ClearFolder(folder string) error {
	files, err := os.ReadDir(folder)

	if err != nil {
		return &generalerrors.ReadDirFailed{Dir: folder}
	}

	for _, f := range files {
		path := CreatePath(folder, f.Name())
		err = os.Remove(path)

		if err != nil {
			return &generalerrors.RemoveFileFailed{File: path}
		}
	}

	return nil
}

func IsExecutable(filepath string) bool {
	if runtime.GOOS == "windows" {
		return strings.Contains(filepath, ".exe")
	} else {
		fi, err := os.Stat(filepath)

		if err != nil {
			log.Printf("Error when finding if %s is executable: %s\n", GetFileName(filepath), err.Error())
			return false
		}
		//Check permission bits to tell if a file is executable for any
		return fi.Mode()&0111 != 0
	}
}

func RemoveUninstaller(dir string) error {
	if runtime.GOOS == "windows" {
		uninstallerPath := CreatePath(dir, "SIDE_Uninstaller.exe")
		removeUninstaller := "/c start timeout /t 1 /NOBREAK > NUL && del " + uninstallerPath

		s := exec.Command("cmd.exe", removeUninstaller)
		s.Stdout = os.Stdout
		s.Stdin = os.Stdin
		s.Stderr = os.Stderr

		if err := s.Start(); err != nil {
			return err
		}
		return nil
	} else {
		uninstallerPath := CreatePath(dir, "SIDE_Uninstaller")
		removeUninstaller := " sleep 1 && rm " + uninstallerPath

		s := exec.Command("bash", "-c", removeUninstaller)
		s.Stdout = os.Stdout
		s.Stdin = os.Stdin
		s.Stderr = os.Stderr

		if err := s.Start(); err != nil {
			return err
		}
		return nil
	}
}

func GetFullPathFromArg(pathArg string) string {
	fullpath, err := filepath.Abs(pathArg)

	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}

	return fullpath
}
