//go:build windows
// +build windows

package agent

import (
	// "bytes"

	// "golang.org/x/text/encoding/unicode"

	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

func DecodeUtf16(b []byte, order binary.ByteOrder) (string, error) {
	ints := make([]uint16, len(b)/2)
	if err := binary.Read(bytes.NewReader(b), order, &ints); err != nil {
		return "", err
	}
	return string(utf16.Decode(ints)), nil
}

func DecodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}

func decodeBytes(b []byte) string {
	// if b.Len()%2 != 0 {
	//     return b.String()
	// }
	// found := false
	// for _, v := range b.Bytes() {
	//     if v == 0x00 {
	//         found = true
	//         break
	//     }
	// }
	// if !found {
	//     return b.String()
	// }
	// enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// enc := japanese.ShiftJIS
	// bb, _, err := transform.Bytes(enc.NewDecoder(), b.Bytes())
	// bb, _, err := transform.Bytes(enc.NewDecoder(), b)
	// log.Infof("DECODE B %v", string(b))
	// log.Infof("DECODE BB %v", string(bb))
	bb2, err := DecodeUTF16(b)
	// log.Infof("DECODE BB2 %v", string(bb2))

	log.Debug("command out ", bb2)
	log.Debug("command err ", err)
	if err != nil {
		return string(b)
	}
	return bb2
}
