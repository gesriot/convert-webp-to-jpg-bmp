package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

func convertImage(inputPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(inputPath)
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v", inputPath, err)
		return
	}
	defer file.Close()

	img, err := webp.Decode(file)
	if err != nil {
		log.Printf("Ошибка декодирования WebP %s: %v", inputPath, err)
		return
	}

	baseDir := filepath.Dir(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	outputPathBmp := filepath.Join(baseDir, baseName+".bmp")
	outBmp, err := os.Create(outputPathBmp)
	if err != nil {
		log.Printf("Ошибка создания BMP файла %s: %v", outputPathBmp, err)
		return
	}
	defer outBmp.Close()

	err = bmp.Encode(outBmp, img)
	if err != nil {
		log.Printf("Ошибка кодирования BMP %s: %v", outputPathBmp, err)
		return
	}

	outputPathJpg := filepath.Join(baseDir, baseName+".jpg")
	err = imaging.Save(img, outputPathJpg, imaging.JPEGQuality(100))
	if err != nil {
		log.Printf("Ошибка сохранения JPG %s: %v", outputPathJpg, err)
		return
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Укажите путь к изображению(ям)")
	}

	var wg sync.WaitGroup
	for _, inputPath := range os.Args[1:] {
		wg.Add(1)
		go convertImage(inputPath, &wg)
	}
	wg.Wait()
}
