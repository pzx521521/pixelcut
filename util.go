package pixelcut

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func DecodeConfig(filePath string) (*image.Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	return &config, err
}

// Calculate16_9 计算将图片转换为 w:h 比例时需要的上下左右的拓宽量
// 如转换为宽图片 16:9 w=16 h=9
func Calculate(filePath string, w, h int) (top int, left int) {
	config, err := DecodeConfig(filePath)
	if err != nil || config == nil {
		return 0, 0 // 若读取失败，返回零
	}
	width := config.Width
	height := config.Height

	// 计算目标的 16:9 高度和宽度
	targetWidth := height * w / h
	left = (targetWidth - width) / 2
	return top, left
}

func CreatDir(saveDir string) error {
	if FileExists(saveDir) {
		return nil
	}
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func GetAllFilesByExts(dirPth string, exts []string) ([]string, error) {
	var files []string
	skipOutpaint := true
	if filepath.Base(dirPth) == saveDirName {
		skipOutpaint = false
	}
	err := filepath.WalkDir(dirPth, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//不进行递归
		if d.IsDir() && d.Name() == saveDirName && skipOutpaint {
			return filepath.SkipDir
		}

		if slices.Contains(exts, filepath.Ext(path)) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// FileExists 判断文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func ChangePng2Jpg(dirPth string) error {
	files, err := GetAllFilesByExts(dirPth, []string{".png"})
	if err != nil {
		return err
	}
	for _, file := range files {
		outputFile := strings.TrimSuffix(file, ".png") + ".jpg"
		if FileExists(outputFile) {
			return nil
		}
		err := Png2Jpg(file, outputFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func Png2Jpg(inputFile, outputFile string) error {
	// 打开输入文件
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码 PNG 图像
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 编码为 JPG 格式
	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer dstFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件时出错: %v", err)
	}

	return nil
}
func filterFile(dirPath string) []string {
	files, _ := GetAllFilesByExts(dirPath, []string{".jpg", ".jpeg"})
	ret := []string{}
	for _, file := range files {
		config, err := DecodeConfig(file)
		if err != nil {
			continue
		}
		if config.Width < 2000 && config.Height < 1000 {
			ret = append(ret, file)
		}
	}
	return ret
}
