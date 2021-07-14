// Windows filesystem monitoring implementation based on
// golang.org/x/exp/winfsnotify
// (specifically
// https://github.com/golang/exp/tree/c84be7c6d1cd7b6a43fd7101daaf2dc35ded445f/winfsnotify),
// but modified to remove import path enforcement, increase
// ReadDirectoryChangesW buffer size, support recursive watching, use more
// idiomatic filesystem path joins, and remove test logging.
//
// The original code license:
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// The original license header inside the code itself:
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package winfsnotify

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func expect(t *testing.T, eventstream <-chan *Event, name string, mask uint32) {
	select {
	case event := <-eventstream:
		if event == nil {
			t.Fatal("nil event received")
		}
		if event.Name != name || event.Mask != mask {
			t.Fatal("did not receive expected event")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestNotifyEvents(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher() failed: %s", err)
	}

	testDir := "TestNotifyEvents.testdirectory"
	testFile := filepath.Join(testDir, "TestNotifyEvents.testfile")
	testFile2 := testFile + ".new"
	const mask = FS_ALL_EVENTS & ^(FS_ATTRIB|FS_CLOSE) | FS_IGNORED

	// Add a watch for testDir
	os.RemoveAll(testDir)
	if err = os.Mkdir(testDir, 0777); err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}
	defer os.RemoveAll(testDir)
	err = watcher.AddWatch(testDir, mask)
	if err != nil {
		t.Fatalf("Watcher.Watch() failed: %s", err)
	}

	// Create a file
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("creating test file failed: %s", err)
	}
	expect(t, watcher.Event, testFile, FS_CREATE)

	err = watcher.AddWatch(testFile, mask)
	if err != nil {
		t.Fatalf("Watcher.Watch() failed: %s", err)
	}

	if _, err = file.WriteString("hello, world"); err != nil {
		t.Fatalf("failed to write to test file: %s", err)
	}
	if err = file.Close(); err != nil {
		t.Fatalf("failed to close test file: %s", err)
	}
	expect(t, watcher.Event, testFile, FS_MODIFY)
	expect(t, watcher.Event, testFile, FS_MODIFY)

	if err = os.Rename(testFile, testFile2); err != nil {
		t.Fatalf("failed to rename test file: %s", err)
	}
	expect(t, watcher.Event, testFile, FS_MOVED_FROM)
	expect(t, watcher.Event, testFile2, FS_MOVED_TO)
	expect(t, watcher.Event, testFile, FS_MOVE_SELF)

	if err = os.RemoveAll(testDir); err != nil {
		t.Fatalf("failed to remove test directory: %s", err)
	}
	expect(t, watcher.Event, testFile2, FS_DELETE_SELF)
	expect(t, watcher.Event, testFile2, FS_IGNORED)
	expect(t, watcher.Event, testFile2, FS_DELETE)
	expect(t, watcher.Event, testDir, FS_DELETE_SELF)
	expect(t, watcher.Event, testDir, FS_IGNORED)

	if err = watcher.Close(); err != nil {
		t.Fatalf("failed to close watcher: %s", err)
	}

	// Check for errors
	if err := <-watcher.Error; err != nil {
		t.Fatalf("error received: %s", err)
	}
}

func TestNotifyClose(t *testing.T) {
	watcher, _ := NewWatcher()
	watcher.Close()

	var done int32
	go func() {
		watcher.Close()
		atomic.StoreInt32(&done, 1)
	}()

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&done) == 0 {
		t.Fatal("double Close() test failed: second Close() call didn't return")
	}

	err := watcher.Watch(t.TempDir())
	if err == nil {
		t.Fatal("expected error on Watch() after Close(), got nil")
	}
}

func TestWatchLongPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "watch-extended-path")
	if err != nil {
		t.Fatalf("TempDir failed: %s", err)
	}
	defer os.RemoveAll(dir)
	// Create a path longer than syscall.MAX_PATH*2
	path := dir
	for {
		if len(path) > syscall.MAX_PATH*2 {
			break
		}
		path = filepath.Join(path, "another-dir")
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatalf("MkdirAll(%s) failed: %s", path, err)
	}

	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher() failed: %s", err)
	}

	// Receive errors on the error channel on a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		for err := range watcher.Error {
			t.Fatalf("error received: %s", err)
		}
	}()
	defer watcher.Close()

	err = watcher.AddWatch(dir, FS_ALL_EVENTS)
	if err != nil {
		t.Fatalf("Watcher.Watch() failed: %s", err)
	}
	newFile := filepath.Join(path, "new-file")
	err = ioutil.WriteFile(newFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %s", err)
	}
	expect(t, watcher.Event, newFile, FS_CREATE)
}
