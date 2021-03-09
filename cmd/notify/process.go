package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

func terminateDuplicateProcess() {
	myPID := os.Getpid()

	for _, pid := range getDuplicateProcessIDs(myPID) {
		if pid == myPID {
			continue
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			log.Fatalf("Unable to find process : %v", err)
		}

		err = process.Kill()
		if err != nil {
			log.Fatalf("Unable to kill process : %v", err)
		}
	}
}

func getDuplicateProcessIDs(myPID int) []int {
	handle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		log.Fatalf("Unable to get process handle: %v", err)
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = syscall.Process32First(handle, &entry)
	if err != nil {
		log.Fatalf("Unable to get first process entry: %v", err)
	}

	myExeFile := ``
	processMap := make(map[string][]int, 1)

	for {
		exeFile := getExeFile(entry)

		processMap[exeFile] = append(processMap[exeFile], int(entry.ProcessID))
		if int(entry.ProcessID) == myPID {
			myExeFile = exeFile
		}

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			log.Fatalf("Unable to get next process entry: %v", err)
		}
	}

	return processMap[myExeFile]
}

func getExeFile(entry syscall.ProcessEntry32) string {
	end := 0
	for {
		if entry.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return syscall.UTF16ToString(entry.ExeFile[:end])
}
