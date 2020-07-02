// +build !windows

package agent

func decodeBytes(b []byte) string {
	return string(b)
}
