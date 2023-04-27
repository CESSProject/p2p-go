package core

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	fileTransferProto = "/kldr/sft/1"
)

func (n *Node) StarFileTransferProtocol() {
	n.host.SetStreamHandler(fileTransferProto, func(s network.Stream) {
		// Create a buffer stream for non blocking read and write
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go n.recvFile(rw)
	})
}

func (n *Node) SendFile(ctx context.Context, id peer.ID, path string) error {
	// Open stream from this host to the target host
	s, err := n.host.NewStream(ctx, id, fileTransferProto)
	if err != nil {
		return err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Send filename
	filename := filepath.Base(path)
	_, err = rw.Write([]byte(fmt.Sprintf("%s\n", filename)))
	if err != nil {
		return err
	}
	err = rw.Flush()
	if err != nil {
		return err
	}

	// Send file size
	fileStats, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = rw.Write([]byte(fmt.Sprintf("%d\n", fileStats.Size())))
	if err != nil {
		return err
	}
	err = rw.Flush()
	if err != nil {
		return err
	}

	// Send file chunks
	hashFunc := sha256.New()
	buffer := make([]byte, BufferSize)
	for {
		_, err = file.Read(buffer)
		if err != nil {
			// if err is not EOF, error during reading, else stop
			if err != io.EOF {
				return err
			} else {
				break
			}
		}

		_, err = rw.Write(buffer)
		if err != nil {
			return err
		}
		err = rw.Flush()
		if err != nil {
			return err
		}

		hashFunc.Write(buffer)
	}

	// Send file hash
	hash := hashFunc.Sum(nil)

	_, err = rw.Write(append(hash, '\n'))
	if err != nil {
		return err
	}
	err = rw.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) recvFile(rw *bufio.ReadWriter) {
	// Recv filename
	filename, err := rw.ReadString('\n')
	filename = strings.Replace(filename, "\n", "", -1)
	if err != nil {
		return
	}
	savepath := filepath.Join(n.workspace, uuid.NewString())
	file, err := os.Create(savepath)
	if err != nil {
		return
	}
	defer os.Remove(savepath)
	defer file.Close()

	// Recv file size
	str, err := rw.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1)
	if err != nil {
		return
	}
	fileSize, err := strconv.Atoi(str)
	if err != nil {
		return
	}

	// Recv file chunks
	hashFunc := sha256.New()
	buf := make([]byte, BufferSize)
	for bytesRead := 0; bytesRead < fileSize; {
		// Receive bytes
		n, err := rw.Read(buf)
		if err != nil {
			return
		}

		bytesRead += n

		// Write bytes to the file
		file.Write(buf)

		// Add to hash
		hashFunc.Write(buf)
	}

	// Recv file hash
	hashRecv, err := rw.ReadBytes('\n')
	if err != nil {
		return
	}
	hashRecv = bytes.Trim(hashRecv, "\x00") // remove leading or trailing 0x00
	hashRecv = hashRecv[:len(hashRecv)-1]   // remove trailing \n

	hash := hashFunc.Sum(nil)
	if string(hash) == string(hashRecv) {
		os.Rename(savepath, filepath.Join(n.workspace, string(hash)))
	}
}
