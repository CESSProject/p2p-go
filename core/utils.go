/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

// CalcSHA256 is used to calculate the sha256 value
// of the data.
func CalcSHA256(data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("data is nil")
	}
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Get local ip address
func GetLocalIp() ([]string, error) {
	var result = make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return result, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	return result, nil
}

// Get external ip address
func GetExternalIp() (string, error) {
	var (
		err        error
		externalIp string
	)

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	resp, err := client.Get("http://myexternalip.com/raw")
	if err == nil {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		externalIp = fmt.Sprintf("%s", string(b))
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx1, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err := exec.CommandContext(ctx1, "bash", "-c", "curl ifconfig.co").Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err = exec.CommandContext(ctx2, "bash", "-c", "curl cip.cc | grep  IP | awk '{print $3;}'").Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx3, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err = exec.CommandContext(ctx3, "bash", "-c", `curl ipinfo.io | grep \"ip\" | awk '{print $2;}'`).Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\"", "")
		externalIp = strings.ReplaceAll(externalIp, ",", "")
		externalIp = strings.ReplaceAll(externalIp, "\n", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	return "", errors.New("Please check your network status")
}

func CheckExternalIpv4(ip string) bool {
	if isIPv4(ip) {
		values := strings.Split(ip, ".")
		if len(values) < 2 {
			return false
		}
		if values[0] == "10" || values[0] == "127" {
			return false
		}
		if values[0] == "192" && values[1] == "168" {
			return false
		}
		if values[0] == "169" && values[1] == "254" {
			return false
		}
		if values[0] == "172" {
			number, err := strconv.Atoi(values[1])
			if err != nil {
				return false
			}
			if number >= 16 && number <= 31 {
				return false
			}
		}
		return true
	}
	return false
}
