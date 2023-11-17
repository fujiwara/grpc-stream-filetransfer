package grpcp_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fujiwara/grpcp"
)

var (
	testPort = 18022
	testHost = "127.0.0.1"
)

func runServer() {
	ctx := context.Background()
	opt := &grpcp.ServerOption{
		Port:   testPort,
		Listen: testHost,
	}
	go func() {
		err := grpcp.RunServer(context.Background(), opt)
		if err != nil {
			panic("failed to run grpcp server:" + err.Error())
		}
	}()
	client := grpcp.NewClient(&grpcp.ClientOption{Port: testPort})
	for i := 0; i < 3; i++ {
		_, err := client.Ping(ctx, testHost)
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	panic("failed to run grpcp server")
}

func generateRandomBytes(t *testing.T) []byte {
	size := grpcp.StreamBufferSize*2 + 1234
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("failed to generate random bytes: %s", err)
	}
	return b
}

func TestMain(m *testing.M) {
	grpcp.StreamBufferSize = 4096 // for test
	runServer()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestLocalToRemote(t *testing.T) {
	dir := t.TempDir()
	testLocal := filepath.Join(dir, "local.txt")
	content := generateRandomBytes(t)
	if err := os.WriteFile(testLocal, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}
	testRemote := filepath.Join(dir, "remote.txt")

	opt := &grpcp.ClientOption{
		Port:  testPort,
		Quiet: true,
	}
	client := grpcp.NewClient(opt)
	err := client.Copy(
		context.Background(),
		testLocal,
		testHost+":"+testRemote,
	)
	if err != nil {
		t.Fatalf("failed to run grpcp client: %s", err)
	}
	remoteContent, err := os.ReadFile(testRemote)
	if err != nil {
		t.Fatalf("failed to read remote file: %s", err)
	}
	if !bytes.Equal(content, remoteContent) {
		t.Fatalf("content mismatch: expected %d bytes, got %d bytes", len(content), len(remoteContent))
	}
}

func TestRemoteToLocal(t *testing.T) {
	dir := t.TempDir()
	testRemote := filepath.Join(dir, "remote.txt")
	content := generateRandomBytes(t)
	if err := os.WriteFile(testRemote, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}
	testLocal := filepath.Join(dir, "local.txt")
	opt := &grpcp.ClientOption{
		Port:  testPort,
		Quiet: true,
	}
	client := grpcp.NewClient(opt)
	err := client.Copy(
		context.Background(),
		testHost+":"+testRemote,
		testLocal,
	)
	if err != nil {
		t.Fatalf("failed to run grpcp client: %s", err)
	}
	remoteContent, err := os.ReadFile(testRemote)
	if err != nil {
		t.Fatalf("failed to read remote file: %s", err)
	}
	if !bytes.Equal(content, remoteContent) {
		t.Fatalf("content mismatch: expected %d bytes, got %d bytes", len(content), len(remoteContent))
	}

}
