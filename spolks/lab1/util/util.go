package util

import (
	"fmt"
	"os"
)

func ParseArgs() (string, error) {
	args := os.Args

	if len(args) == 1 {
		return "", fmt.Errorf("need minimum 2 arguments (provide host:port)")
	}

	port := args[1]

	return port, nil
}