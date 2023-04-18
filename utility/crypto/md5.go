package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

/**
 * MD5 Hash算法
 */
func Md5Hash(data string) string {
	Md5 := md5.New()
	Md5.Write([]byte(data))
	databytes := Md5.Sum(nil)
	return hex.EncodeToString(databytes)
}
