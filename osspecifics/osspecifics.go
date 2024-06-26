package osspecifics

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
)

var PathSep string

func init() {
	if runtime.GOOS == "windows" {
		PathSep = "\\"
	} else {
		PathSep = "/"
	}
}

// CreatePath - will concatenate different strings to form a file path
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

func LockFile(fp string) error {
	exp, err := os.Executable()

	if err != nil {
		return err
	}

	expdir := filepath.Dir(exp)

	lockf := fmt.Sprintf("%s_open", fp)

	err = os.WriteFile(CreatePath(expdir, lockf), nil, 0666)

	if err != nil {
		return err
	}

	return nil
}

func UnlockFile(fp string) {
	exp, err := os.Executable()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
	}

	expdir := filepath.Dir(exp)

	lockf := fmt.Sprintf("%s_open", fp)

	//If it cannot find the file, then that's it.
	_ = os.Remove(CreatePath(expdir, lockf))
}

func IsLocked(fp string) bool {
	exp, err := os.Executable()

	if err != nil {
		return false
	}

	expdir := filepath.Dir(exp)

	lockf := fmt.Sprintf("%s_open", fp)

	files, err := os.ReadDir(expdir)

	if err != nil {
		return false
	}

	for _, f := range files {
		if f.Name() == lockf {
			return true
		}
	}

	return false
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
