package pixelcut

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const saveDirName = "outpaint"
const creativity = .0
const OUTPAINT_APIURL = "https://api2.pixelcut.app/image/outpaint/v1"

// 处理非jpg图片,跳过已经处理的图片
func PreDo(dirPath string) (map[string]string, error) {
	saveDir := filepath.Join(dirPath, saveDirName)
	err := CreatDir(saveDir)
	if err != nil {
		return nil, err
	}
	err = ChangePng2Jpg(dirPath)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	files, _ := GetAllFilesByExts(dirPath, []string{".jpg", ".jpeg"})
	for _, filePath := range files {
		savePath := filepath.Join(saveDir, filepath.Base(filePath))
		if FileExists(savePath) {
			_, err = DecodeConfig(savePath)
			if err == nil {
				//mybe return "The server encountered a temporary error and could not complete your request."
				//os.Remove(savePath)
				continue
			}
		}
		ret[filePath] = savePath
	}
	return ret, nil
}
func OutPaintDir(client *http.Client, dirPath string) error {
	remainFiles, err := PreDo(dirPath)
	if err != nil {
		return err
	}
	for filePath, savePath := range remainFiles {
		err = OutPaintFile(client, filePath, savePath)
		if err != nil {
			//dont return error, just log it
			log.Printf("error: %v", err)
			continue
		} else {
			log.Printf("save success at:%s", savePath)
		}
	}
	return nil
}
func OutPaintDirByPool(clientPool *ClientPool, dirPath string) error {
	remainFiles, err := PreDo(dirPath)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for filePath, savePath := range remainFiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := clientPool.Get()
			urlInfo := clientPool.GetURL(client)
			log.Printf("start:%s,by clinet:%v", filepath.Base(filePath), urlInfo)
			defer clientPool.Put(client)
			err = OutPaintFile(client, filePath, savePath)
			if err != nil {
				//dont return error, just log it
				log.Printf("error: %v,by clinet:%v, file:%v", err, urlInfo, filepath.Base(filePath))
			} else {
				log.Printf("save success at:%s,by clinet:%v", savePath, urlInfo)

			}
		}()
	}
	wg.Wait()
	return nil
}
func OutPaintFile(client *http.Client, filePath, savePath string) error {
	top, left := Calculate(filePath, 16, 9)
	err := OutPaintPostData(client, filePath, savePath, left, top, left, top, creativity)
	if err != nil {
		return err
	}
	return nil
}
func maxPaint(nums ...*int) {
	for _, num := range nums {
		if *num > 2000 {
			*num = 2000
		}
	}
}
func OutPaintPostData(client *http.Client, filePath, savePath string, left, top, right, bottom int, creativity float64) error {
	if top+left+right+bottom == 0 {
		return errors.New("top = left = 0")
	}
	maxPaint(&left, &top, &right, &bottom)
	// 创建一个新的缓冲区和 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加图像文件字段
	imageFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()
	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, imageFile)
	if err != nil {
		return err
	}
	// 添加其他字段
	writer.WriteField("left", strconv.Itoa(left))
	writer.WriteField("top", strconv.Itoa(top))
	writer.WriteField("right", strconv.Itoa(right))
	writer.WriteField("bottom", strconv.Itoa(bottom))
	writer.WriteField("creativity", fmt.Sprintf("%g", creativity))
	writer.WriteField("output_format", "jpeg")

	// 结束 multipart 请求
	writer.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", OUTPAINT_APIURL, body)
	if err != nil {
		return err
	}

	// 添加请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", writer.FormDataContentType())
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	req.Header.Set("origin", "https://www.pixelcut.ai")
	req.Header.Set("x-client-version", "web")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cross")
	req.Header.Set("sec-fetch-dest", "empty")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取并输出响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, _, err = image.DecodeConfig(bytes.NewBuffer(responseBody))
	if err != nil {
		return errors.New(string(responseBody))
	}
	err = os.WriteFile(savePath, responseBody, 0644)
	return nil
}
