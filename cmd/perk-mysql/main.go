package main

import (
	"fmt"
	"os"

	"github.com/l3aro/mysql-driver-for-perk-workbench/mysql"
	"github.com/l3aro/perk-workbench-plugin-sdk-go/server"
)

func main() {
	if err := server.Run(os.Stdin, os.Stdout, mysql.Factory{}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[perk-mysql] fatal: %v\n", err)
		os.Exit(1)
	}
}
