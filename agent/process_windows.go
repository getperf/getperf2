// +build windows

package agent

import (
    // "bytes"

    // "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"

    log "github.com/sirupsen/logrus"
)

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
    enc := japanese.ShiftJIS
    // bb, _, err := transform.Bytes(enc.NewDecoder(), b.Bytes())
    bb, _, err := transform.Bytes(enc.NewDecoder(), b)
    log.Debug("command out ", string(bb))
    log.Debug("command err ", err)
    if err != nil {
        return string(b)
    }
    return string(bb)
}
