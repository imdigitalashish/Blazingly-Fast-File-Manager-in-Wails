package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type Bytex struct {
	QuadPart int64
}

var (
	kernel32dll           = syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceExPtr = kernel32dll.NewProc("GetDiskFreeSpaceExW")
)

func getSizeDrive(driveLetter string) ([]uint64, error) {
	var freeBytes Bytex
	var totalBytes Bytex
	var totalFreeBytes Bytex

	drivePath := syscall.StringToUTF16Ptr(driveLetter + "\\")
	ret, _, err := getDiskFreeSpaceExPtr.Call(uintptr(unsafe.Pointer(drivePath)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)))

	if ret == 0 {
		return nil, fmt.Errorf("Failed %w", err)
	}

	return []uint64{uint64(totalBytes.QuadPart), uint64(totalFreeBytes.QuadPart)}, nil
}

func getDrives() string {

	var drives = make(map[string]map[string]uint64)

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			driveSize, err := getSizeDrive(string(drive) + ":")
			if err == nil {
				drives[string(drive)] = map[string]uint64{
					"total space": driveSize[0] / 1024 / 1024 / 1024,
					"space left":  driveSize[1] / 1024 / 1024 / 1024, // Assuming we can't get free space here (modify if possible)
				}
			}
		}
	}

	jsonData, err := json.Marshal(drives)
	if err != nil {
		fmt.Println("Error Marshalling map to JSON: ", err)
	}

	return string(jsonData)
}

func getAllFilesAndFolders(path string) ([]string, error) {
	var allFiles []string

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			allFiles = append(allFiles, fullPath) // Add directory path
		} else {
			allFiles = append(allFiles, fullPath) // Add file path
		}
	}

	return allFiles, nil
}

func main() {
	driveLetter := "D:\\" // Adjust drive letter as needed
	folders, err := getAllFilesAndFolders(driveLetter)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(folders) == 0 {
		fmt.Println("No folders found on drive", driveLetter)
	} else {
		fmt.Println("Folders in", driveLetter, "drive:")
		for _, folder := range folders {
			fmt.Println("-", folder)
		}
	}
	fmt.Println(getDrives())
}
