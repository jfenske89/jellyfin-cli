package main

import (
	"fmt"
	"os"

	"codeberg.org/jfenske/jellyfin-cli/cmd/jellyfin-cli/actions"
)

func usage() {
	available := []string{
		actions.ListSessions + ": list active sessions",
		actions.ListLibraryFolders + ": list library folders",
	}

	fmt.Printf("%s <action> [options]\n", os.Args[0])

	for i := range available {
		fmt.Printf(" - %s\n", available[i])
	}
}
