package pixelcut

import (
	"log"
	"net/http"
	"path/filepath"
	"testing"
)

func TestCalculate(t *testing.T) {
	top, left := Calculate("/Users/parapeng/Downloads/pinterest/灰橙黑aa1816.jpg", 16, 9)
	log.Printf("top:%d, left:%d\n", top, left)
}

func TestOutPaintFile(t *testing.T) {
	filePath := "./test.jpg"
	savePath := filepath.Join(filepath.Dir(filePath), saveDirName, filepath.Base(filePath))
	err := OutPaintFile(http.DefaultClient, filePath, savePath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func TestOutPaintDir(t *testing.T) {
	dirPath := "/Users/parapeng/Downloads/wait"
	client := NewProxyClientByUrl("http://127.0.0.1:8888")
	err := OutPaintDir(client, dirPath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
func TestUpscalerFile(t *testing.T) {
	files := FilterFile("/Users/parapeng/Downloads/outpaint", 2000, 1000)
	client := NewProxyClientByUrl("http://127.0.0.1:7897")
	newFiles := files[:1]
	err := UpscaleFiles(client, newFiles)
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func TestOutPaintDirWithPool(t *testing.T) {
	dirPath := "/Users/parapeng/Downloads"
	proxies := []string{"http://127.0.0.1:7897"}
	pool := NewClientPool(proxies, 5)
	err := OutPaintDirByPool(pool, dirPath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func TestUpscaleFilesByPool(t *testing.T) {
	files := FilterFile("/Users/parapeng/Downloads/outpaint", 2000, 1000)
	proxies := []string{"http://127.0.0.1:7897"}
	pool := NewClientPool(proxies, 5)
	err := UpscaleFilesByPool(pool, files)
	if err != nil {
		log.Printf("%v\n", err)
	}
}
