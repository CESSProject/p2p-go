/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// CalcPathSHA256 is used to calculate the sha256 value
// of a file with a given path.
func CalcPathSHA256(fpath string) (string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return CalcFileSHA256(f)
}

// CalcFileSHA256 is used to calculate the sha256 value
// of the file type.
func CalcFileSHA256(f *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
