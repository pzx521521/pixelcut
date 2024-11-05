package pixelcut

import (
	"fmt"
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
	dirPath := "/Users/parapeng/Downloads/pinterest"
	client := NewProxyClientByUrl("http://127.0.0.1:8888")
	err := OutPaintDir(client, dirPath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
func TestOutPaintDirWithPoolCharles(t *testing.T) {
	dirPath := "/Users/parapeng/Downloads/pinterest"
	proxies := []string{"http://127.0.0.1:8888"}
	//proxies := []string{}
	//for i := 0; i < 6; i++ {
	//	proxies = append(proxies, fmt.Sprintf("http://127.0.0.1:%d", i+7000))
	//}
	pool := NewClientPool(proxies, 1)
	err := OutPaintDirByPool(dirPath, pool)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
func TestOutPaintDirWithPool(t *testing.T) {
	dirPath := "/Users/parapeng/Downloads/pinterest"
	proxies := []string{"http://127.0.0.1:8888"}
	for i := 0; i < 6; i++ {
		proxies = append(proxies, fmt.Sprintf("http://127.0.0.1:%d", i+7000))
	}
	pool := NewClientPool(proxies, 1)
	err := OutPaintDirByPool(dirPath, pool)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
