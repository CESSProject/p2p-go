/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ds "github.com/ipfs/go-datastore"
	query "github.com/ipfs/go-datastore/query"
)

var data = ".data"

// Datastore uses a uses a file per key to store values.
type Datastore struct {
	path string
}

var _ ds.Datastore = (*Datastore)(nil)
var _ ds.Batching = (*Datastore)(nil)
var _ ds.PersistentDatastore = (*Datastore)(nil)

// NewDatastore returns a new fs Datastore at given `path`
func NewDatastore(path string) (ds.Datastore, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	if !isDir(path) {
		return nil, fmt.Errorf("failed to find directory at: %v (file? perms?)", path)
	}

	return &Datastore{path: path}, nil
}

// KeyFilename returns the filename associated with `key`
func (d *Datastore) KeyFilename(key ds.Key) string {
	return filepath.Join(d.path, key.Name(), data)
}

// Put stores the given value.
func (d *Datastore) Put(ctx context.Context, key ds.Key, value []byte) (err error) {
	fn := d.KeyFilename(key)
	// mkdirall above.
	err = os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(fn, value, 0666)
}

// Sync would ensure that any previous Puts under the prefix are written to disk.
// However, they already are.
func (d *Datastore) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

// Get returns the value for given key
func (d *Datastore) Get(ctx context.Context, key ds.Key) (value []byte, err error) {
	fn := d.KeyFilename(key)
	if !isFile(fn) {
		return nil, ds.ErrNotFound
	}

	return os.ReadFile(fn)
}

// Has returns whether the datastore has a value for a given key
func (d *Datastore) Has(ctx context.Context, key ds.Key) (exists bool, err error) {
	return ds.GetBackedHas(ctx, d, key)
}

func (d *Datastore) GetSize(ctx context.Context, key ds.Key) (size int, err error) {
	return ds.GetBackedSize(ctx, d, key)
}

// Delete removes the value for given key
func (d *Datastore) Delete(ctx context.Context, key ds.Key) (err error) {
	fn := d.KeyFilename(key)
	if !isFile(fn) {
		return nil
	}

	err = os.Remove(fn)
	if os.IsNotExist(err) {
		err = nil // idempotent
	}
	return err
}

// Query implements Datastore.Query
func (d *Datastore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	results := make(chan query.Result)

	walkFn := func(path string, info os.FileInfo, _ error) error {
		// remove ds path prefix
		relPath, err := filepath.Rel(d.path, path)
		if err == nil {
			path = filepath.ToSlash(relPath)
		}

		if !info.IsDir() {
			path = strings.TrimSuffix(path, data)
			var result query.Result
			key := ds.NewKey(path)
			result.Entry.Key = key.String()
			if !q.KeysOnly {
				result.Entry.Value, result.Error = d.Get(ctx, key)
			}
			results <- result
		}
		return nil
	}

	go func() {
		filepath.Walk(d.path, walkFn)
		close(results)
	}()
	r := query.ResultsWithChan(q, results)
	r = query.NaiveQueryApply(q, r)
	return r, nil
}

// isDir returns whether given path is a directory
func isDir(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return finfo.IsDir()
}

// isFile returns whether given path is a file
func isFile(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !finfo.IsDir()
}

func (d *Datastore) Close() error {
	return nil
}

func (d *Datastore) Batch(ctx context.Context) (ds.Batch, error) {
	return ds.NewBasicBatch(d), nil
}

// DiskUsage returns the disk size used by the datastore in bytes.
func (d *Datastore) DiskUsage(ctx context.Context) (uint64, error) {
	var du uint64
	err := filepath.Walk(d.path, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if f != nil && f.Mode().IsRegular() {
			du += uint64(f.Size())
		}
		return nil
	})
	return du, err
}
