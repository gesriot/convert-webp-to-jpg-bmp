package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

func convertImage(input_path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(input_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := webp.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	output_path_bmp := filepath.Dir(input_path) + "/" + filepath.Base(input_path[:len(input_path)-len(filepath.Ext(input_path))]) + ".bmp"

	out_bmp, err := os.Create(output_path_bmp)
	if err != nil {
		log.Fatal(err)
	}
	defer out_bmp.Close()

	err = bmp.Encode(out_bmp, img)
	if err != nil {
		log.Fatal(err)
	}

	output_path_jpg := filepath.Dir(input_path) + "/" + filepath.Base(input_path[:len(input_path)-len(filepath.Ext(input_path))]) + ".jpg"

	err = imaging.Save(img, output_path_jpg, imaging.JPEGQuality(100))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Укажите путь к изображению(ям)")
	}

	var wg sync.WaitGroup

	for _, input_path := range os.Args[1:] {
		wg.Add(1)
		go convertImage(input_path, &wg)
	}

	wg.Wait()
}
