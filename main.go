package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
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
			err := filepath.Walk(directory, ProcessFile)

			if err != nil {
				log.Fatal(err)
			}
		})

	commando.Parse(nil)
}

func ProcessFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.Contains(path, "@eaDir") {
		return nil
	}

	matched, err := filepath.Match("*.[jJ][pP]*[gG]", filepath.Base(path))
	if err != nil {
		return err
	} else if matched {
		err = ProcessImage(path)

		if err != nil {
			fmt.Println("ERROR: Image not processed", path)
		}
	}

	return nil
}

func ProcessImage(imagePath string) error {
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

	fullImagePath := filepath.Join(thumbDir, imagePath)
	fmt.Println("Thumbnails generated for", fullImagePath)

	return nil
}

func ImageThumbnailSizesAndWidths() map[string]int {
	return map[string]int{
		"SYNOPHOTO_THUMB_XL.jpg": 1280,
		"SYNOPHOTO_THUMB_SM.jpg": 320,
		"SYNOPHOTO_THUMB_M.jpg":  240,
	}
}

func VideoThumbnailSizesAndWidths() map[string]int {
	return map[string]int{
		"SYNOPHOTO_THUMB_XL.jpg": 1280,
		"SYNOPHOTO_THUMB_SM.jpg": 320,
		"SYNOPHOTO_THUMB_M.jpg":  240,
	}
}

func CreateThumbnails(thumbDir string, imagePath string) error {
	src, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	for dstFilename, dstMaxSize := range ImageThumbnailSizesAndWidths() {
		var m *image.NRGBA

		if src.Bounds().Size().X > src.Bounds().Size().Y {
			m = imaging.Resize(src, dstMaxSize, 0, imaging.Lanczos)
		} else {
			m = imaging.Resize(src, 0, dstMaxSize, imaging.Lanczos)
		}

		dst := imaging.New(m.Bounds().Size().X, m.Bounds().Size().Y, color.NRGBA{0, 0, 0, 0})
		dst = imaging.Paste(dst, m, image.Pt(0, 0))

		err = imaging.Save(dst, filepath.Join(thumbDir, dstFilename))
		if err != nil {
			return err
		}
	}

	return nil
}
