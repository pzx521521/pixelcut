package pixelcut

import (
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

func GetIP(pool *ClientPool) error {
	// 使用客户端池
	client := pool.Get()
	defer pool.Put(client)

	//resp, err := client.Get("https://www.pixelcut.ai/t/uncrop")
	resp, err := client.Get("https://directory.cookieyes.com/api/v1/ip")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	//模拟长连接
	time.Sleep(1 * time.Second)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client:%v, err:%v", client.Transport, err)
	}
	log.Printf("resp.Body: %s", body)
	return nil
}
func TestPoolHttpByBalance(t *testing.T) {
	proxies := []string{"http://127.0.0.1:8888"}
	var wg sync.WaitGroup
	pool := NewClientPool(proxies, 5)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := GetIP(pool)
			if err != nil {
				log.Printf("err:%v", err)
				return
			}
		}()
	}
	wg.Wait()
}
func TestPoolHttpByLinxDo(t *testing.T) {
	proxies := []string{"https://api-proxy.me/anthropic"}
	var wg sync.WaitGroup
	pool := NewClientPool(proxies, 1)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := GetIP(pool)
			if err != nil {
				log.Printf("err:%v", err)
				return
			}
		}()
	}
	wg.Wait()
}
