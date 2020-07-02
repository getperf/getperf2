// +build !windows

package common

import (
	"fmt"
	"syscall"
)

// CheckDiskFree は指定したディレクトリのディスク使用量[%]を取得します。
func CheckDiskFree(dir string) (DiskStatus, error) {
	disk := DiskStatus{}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(dir, &fs)
	if err != nil {
		return disk, fmt.Errorf("check disk free %s", err)
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return disk, nil
}
