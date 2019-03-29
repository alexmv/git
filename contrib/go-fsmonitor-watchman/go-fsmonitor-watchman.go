package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cstrahan/go-watchman/cmd"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments")
	}

	version := os.Args[1]
	if version != "1" {
		log.Fatalf("Unsupported fsmonitor hook version %s", version)
	}

	time_ns, err := strconv.ParseInt(os.Args[2], 0, 64)
	if err != nil {
		log.Fatalf("Time argument (%s) cannot be parsed: %s", os.Args[2], err)
	}
	time_seconds := int(time_ns) / 1e9

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get working directory: %s", err)
	}

	query := []interface{}{
		"query",
		dir,
		map[string]interface{}{
			"fields": []interface{}{"name"},
			"since":  time_seconds,
		},
	}

	res, err := cmd.Command("watchman", query)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unable to resolve root") &&
			strings.HasSuffix(err.Error(), "is not watched") {
			fmt.Fprintf(os.Stderr, "Adding %s to watchman's watch list\n", dir)
			command := []interface{}{
				"watch-project",
				dir,
			}
			_, err = cmd.Command("watchman", command)
			if err != nil {
				log.Fatalf("Failed to make watchman watch %s: %s", dir, err)
			}
			// The first call would always return all files, so just emulate that by
			// telling git that everything is dirty.
			fmt.Print("/")
			return
		} else {
			log.Fatalf("Unknown watchman error: %s", err)
		}
	}

	file_interfaces := res.(map[string]interface{})["files"].([]interface{})
	files := make([]string, len(file_interfaces))
	for i, v := range file_interfaces {
		files[i] = v.(string)
	}

	fmt.Print(strings.Join(files, "\x00"))
}
