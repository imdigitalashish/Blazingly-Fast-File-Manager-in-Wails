package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func getAllFilesAndFolders(path string) string {
	allFiles := []map[string]string{}

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return "Error reading directory"
	}

	for _, entry := range entries {
		// fullPath := filepath.Join(path, entry.Name())
		folder_file_type := make(map[string]string)

		if entry.IsDir() {
			folder_file_type[entry.Name()] = "folder"
		} else {
			folder_file_type[entry.Name()] = "file"
		}
		allFiles = append(allFiles, folder_file_type)
	}

	jsonData, err := json.Marshal(allFiles)

	return string(jsonData)
}

func getDrives() string {

	var drives = make(map[string]map[string]uint64)

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			driveSize, err := getSizeDrive(string(drive) + ":")
			if err == nil {
				drives[string(drive)] = map[string]uint64{
					"total_space": driveSize[0] / 1024 / 1024 / 1024,
					"space_left":  driveSize[1] / 1024 / 1024 / 1024, // Assuming we can't get free space here (modify if possible)
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
