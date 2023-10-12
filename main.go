package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type MediaFile struct {
	Status  string
	Src     string
	Title   string
	Link    string
	Alt     string
	Caption string
	Thumb   string
}

type GalleryFile struct {
	GalleryItems map[string]MediaFile `json:"gallery"`
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Path canont be empty\n")
		os.Exit(1)
	}

	// check if the path exsts
	path := os.Args[1]
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Given target path does not exists: %s\n", path)
		os.Exit(1)
	}

	// get files
	files, err := filepath.Glob(filepath.Join(path, "*.json"))
	if err != nil {
		fmt.Printf("Something went wrong while loading files... %s\n", err)
		os.Exit(1)
	}

	total := 0
	bar := pb.StartNew(len(files))

	for _, file := range files {
		outputPath := strings.TrimSuffix(filepath.Base(file), ".json")

		// make output dir
		os.MkdirAll(filepath.Join("./output", outputPath), 0755)

		// load file and parse JSON
		body, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file: %s .. %s\n", file, err)
			os.Exit(1)
		}

		data := GalleryFile{}
		json.Unmarshal(body, &data)

		bar2 := pb.StartNew(len(data.GalleryItems))
		for _, media := range data.GalleryItems {
			if len(media.Src) == 0 {
				fmt.Printf("Skip a record because of a missing link")
				continue
			}
			fn := filepath.Base(media.Src)

			DownloadFile(media.Link, filepath.Join("./output", outputPath, fn))

			bar2.Increment()
			total += 1
		}
		bar2.Finish()

		bar.Increment()
	}

	bar.Finish()

	fmt.Printf("Finished download in total %d files..\n", total)
}

func DownloadFile(link string, target string) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("Couldnt download file: %s .. %s\n", link, err)
		os.Exit(1)
	}

	out, err := os.Create(target)
	if err != nil {
		fmt.Printf("Cannot create target: %s .. %s\n", target, err)
		os.Exit(1)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Download failed: %s .. %s\n", link, err)
	}

}
