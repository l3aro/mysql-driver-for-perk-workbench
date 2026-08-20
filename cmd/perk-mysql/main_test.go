package main

import (
	"bufio"
	"encoding/json"
	"io"
	"testing"

	"github.com/l3aro/mysql-driver-for-perk-workbench/mysql"
	"github.com/l3aro/perk-workbench-plugin-sdk-go/server"
)

func TestLifecycleTransport(t *testing.T) {
	inputReader, inputWriter := io.Pipe()
	outputReader, outputWriter := io.Pipe()
	done := make(chan error, 1)
	go func() { done <- server.Run(inputReader, outputWriter, mysql.Factory{}) }()
	readResponse := func() map[string]any {
		t.Helper()
		line, err := bufio.NewReader(outputReader).ReadBytes('\n')
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]any
		if err := json.Unmarshal(line, &response); err != nil {
			t.Fatal(err)
		}
		return response
	}
	writeRequest := func(value string) {
		t.Helper()
		if _, err := io.WriteString(inputWriter, value+"\n"); err != nil {
			t.Fatal(err)
		}
	}
	writeRequest(`{"jsonrpc":"2.0","id":1,"method":"perk/v1/initialize","params":{"protocol_version":1,"workbench_version":"test"}}`)
	initialize := readResponse()
	capabilities := initialize["result"].(map[string]any)["capabilities"].(map[string]any)
	if capabilities["name"] != "mysql" {
		t.Fatalf("capabilities name = %v, want mysql", capabilities["name"])
	}
	writeRequest(`{"jsonrpc":"2.0","id":2,"method":"perk/v1/build_target","params":{"host":"db","port":"3306","user":"alice","pass":"secret","database":"app","tls":"false"}}`)
	target := readResponse()
	if target["result"].(map[string]any)["target"] != "alice:secret@tcp(db:3306)/app?tls=false" {
		t.Fatalf("build_target response = %#v", target)
	}
	if err := inputWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
