package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/thatisuday/commando"
)

func main() {
	commandoConfig()
}

func commandoConfig() {
	commando.
		SetExecutableName("syno-thumb").
		SetVersion("1.0.0").
		SetDescription("Creates thumbnails ready to be used by a Synology system")

	commando.
		Register(nil).
		AddFlag("dir", "scans local path for images", commando.String, "./").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			directory, _ := flags["dir"].GetString()
			_, err := ProcessDirectory(directory)

			if err != nil {
				log.Fatal(err)
			}
		})

	commando.Parse(nil)
}

func ProcessDirectory(directory string) ([]string, error) {
	var matches []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match("*.[jJ][pP][gG]", filepath.Base(path))
		if err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
			err = ProcessFile(path)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func ProcessFile(imagePath string) error {
	fmt.Println("Processing", imagePath)

	path, filename := filepath.Split(imagePath)
	thumbDir := filepath.Join(path, "@eaDir", filename)

	if _, err := os.Stat(thumbDir); os.IsNotExist(err) {
		err := os.MkdirAll(thumbDir, 0700)
		if err != nil {
			return err
		}
	}

	err := CreateThumbnails(thumbDir, imagePath)
	if err != nil {
		return err
	}

	fmt.Println("Thumbnails generated for", imagePath)

	return nil
}

func CreateThumbnails(thumbDir string, imagePath string) error {
	return nil
}
