package pixelcut

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type UpscalerResp struct {
	ResultUrl string `json:"result_url"`
}

const UPSCALER_APIURL = "https://api2.pixelcut.app/image/upscale/v1"

func UpscalerFiles(client *http.Client, files []string) error {
	for _, filePath := range files {
		err := UpscalerFile(client, filePath)
		if err != nil {
			//dont return error, just log it
			log.Printf("error: %v", err)
			continue
		} else {
			log.Printf("save success at:%s", filePath)
		}
	}
	return nil
}
func UpscalerFilesByPool(clientPool *ClientPool, files []string) error {
	var wg sync.WaitGroup
	for _, filePath := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := clientPool.Get()
			urlInfo := clientPool.GetURL(client)
			log.Printf("start:%s,by clinet:%v", filepath.Base(filePath), urlInfo)
			defer clientPool.Put(client)
			err := UpscalerFile(client, filePath)
			if err != nil {
				//dont return error, just log it
				log.Printf("error: %v,by clinet:%v", err, urlInfo)
			} else {
				log.Printf("save success at:%s,by clinet:%v", filePath, urlInfo)
			}
		}()
	}
	wg.Wait()
	return nil
}
func UpscalerFile(client *http.Client, filePath string) error {
	return UpscalerPostData(client, filePath)
}

func UpscalerPostData(client *http.Client, filePath string) error {
	// 创建一个新的缓冲区和 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 添加图像文件字段
	imageFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()
	part, err := writer.CreateFormFile("image", "blob")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, imageFile)
	if err != nil {
		return err
	}
	// 结束 multipart 请求
	writer.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", UPSCALER_APIURL, body)
	if err != nil {
		return err
	}
	// 添加请求头
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
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(responseBody))
	}
	contentType := resp.Header.Get("content-type")
	switch contentType {
	case "image/jpeg":
		return os.WriteFile(filePath, responseBody, 0644)
	case "application/json":
		var upscalerResp UpscalerResp
		json.Unmarshal(responseBody, &upscalerResp)
		return downloadImage(client, upscalerResp.ResultUrl, filePath)
	}
	return errors.New("unknow format: " + contentType)
}

func downloadImage(client *http.Client, imgUrl, savePath string) error {
	// 发送 HTTP 请求获取图片内容
	resp, err := client.Get(imgUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	// 创建文件
	os.WriteFile(savePath, data, 0644)
	return nil
}
