package util

import (
	"fmt"
	"os"
	"time"
)

func Exit(error string,t time.Duration)  {
	fmt.Println(error)
	if t > 0 {
		time.Sleep(t)
		os.Exit(1)
	}
}

func ChangeSize(size int) string{
	if size < 1024 {
		//return strconv.FormatInt(size, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else { //if size < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
}

func IsExist(path string) bool{
	existed := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		existed = false
	}
	return existed
}
