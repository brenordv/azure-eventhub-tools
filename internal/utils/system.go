package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

// GetAppDir returns the path of this application. if it's not loaded yet, will save it to a global variable.
// Will panic if it cannot get the application path.
//
// Parameters:
//  None.
//
// Returns:
//  String containing the path of the application.
func GetAppDir() string {
	if d.AppDir == "" {
		execPath, err := os.Executable()
		h.HandleError("Failed to get application directory.", err, true)
		d.AppDir = path.Dir(execPath)
	}
	return d.AppDir
}

// Exists returns whether the given file or directory exists.
// Will panic in case of failure.
//
// Parameters:
//  path: path that will be checked.
//
// Returns:
//  true if the path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)
}

// EnsureExists if the path does not exist, tries to create it.
//
// Parameters:
//  path: path that will be checked.
//
// Returns:
//  error if any happens or nil if all is ok.
func EnsureExists(path string) error {
	if Exists(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to create path '%s': %v", path, err)
}

// MustGetFileStat gets the FileInfo for a given file.
// Will panic on error.
//
// Parameters:
//   p: target file
//
// Returns:
//   FileInfo for p.
func MustGetFileStat(p string) os.FileInfo {
	fi, err := os.Stat(p)
	h.HandleError(fmt.Sprintf("Failed to get file stats for '%s'.", p), err, true)
	return fi
}

// TODO: add summary
func IsFile(p string) bool {
	fi := MustGetFileStat(p)
	return !fi.IsDir()
}

// TODO: add summary
func LoadRuntimeConfig(cfgFile string, validator func()) {
	if !Exists(cfgFile) {
		h.HandleError(
			"Config file validation failed!",
			fmt.Errorf("file '%s' does not exist or is unaccessible", cfgFile),
			true)
	}

	file, err := os.Open(cfgFile)
	h.HandleError(fmt.Sprintf("Failed to open Config file: %s", cfgFile), err, true)

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&d.CurrentConfig)
	h.HandleError(fmt.Sprintf("Failed to read Config file: %s", cfgFile), err, true)

	validator()
}

// ReadTextFile will read a text file and return it's content. If something goes wrong, will explode.
// Will panic read fails.
//
// Parameters:
//  f: path to the file that will be read.
//
// Returns:
//  String content of the file.
func ReadTextFile(f string) string {

	file, err := os.Open(f)
	h.HandleError(fmt.Sprintf("Failed to open file '%s'", f), err, true)
	defer h.CloseWithErrorHandling(file.Close, fmt.Sprintf("Failed to close file '%s'", f), true)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 64)
	var read []byte

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		h.HandleError(fmt.Sprintf("Error reading file '%s'.", f), err, true)
		read = append(read, buffer[0:n]...)
	}

	return string(read)
}

// GetSubFolderBasedOnTime returns a string containing a folder
// Will panic on error.
//
// Parameters:
//   p: base folder.
//   t: time that will be used to create the sub-folder
//
// Returns:
//   concatenated path of p and t (yyyy-mm-dd).
func GetSubFolderBasedOnTime(p string, t time.Time) string {
	if pi := MustGetFileStat(p); !pi.IsDir() {
		h.HandleError("Informed parameter is a file.",
			fmt.Errorf("parameter '%s' must be a folder", p), true)
	}

	dir := filepath.Join(p, t.Format("2006-01-02"))
	err := EnsureExists(dir)
	h.HandleError("Failed to create sub-folder based on time.", err, true)
	return dir
}


// PutFileInSubFolderBasedOnTime generates a path+filename based on current time and the filename.
// Will panic in case of failure.
//
// Parameters:
//   baseFolder: folder to be used as base to construct full path.
//   filename: desired filename.
//   t: time to be used as reference.
//
// Returns:
//   filename that will be used to dump an eventhub message.
func PutFileInSubFolderBasedOnTime(baseFolder string, filename string, t time.Time) string {
	dir := GetSubFolderBasedOnTime(baseFolder, t)
	return filepath.Join(dir, fmt.Sprintf("%s--%s", t.Format("2006-01-02T15-04-05.00"), filename))
}