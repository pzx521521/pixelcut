package pixelcut

import (
	"log"
	"net/http"
	"net/url"
)

type ClientPool struct {
	mapClients map[*http.Client]*url.URL
	clients    chan *http.Client
	validURLs  []*url.URL
}

func (p *ClientPool) AddURL(c *http.Client, u *url.URL) {
	if _, ok := p.mapClients[c]; !ok {
		p.mapClients[c] = u
	}
}
func (p *ClientPool) GetURL(c *http.Client) *url.URL {
	if _, ok := p.mapClients[c]; ok {
		return p.mapClients[c]
	}
	return nil
}

// 初始化客户端池
func NewClientPool(proxyURLs []string, batchSize int) *ClientPool {
	validURLs := checkProxyURLs(proxyURLs)
	chanSize := batchSize * len(validURLs)
	pool := &ClientPool{
		clients:    make(chan *http.Client, chanSize),
		validURLs:  validURLs,
		mapClients: make(map[*http.Client]*url.URL),
	}
	// 初始化池，填充最大数量的客户端
	for i := 0; i < batchSize; i++ {
		for j := 0; j < len(validURLs); j++ {
			client, _ := newProxyClient(validURLs[j])
			if i == 0 {
				pool.AddURL(client, validURLs[j])
			}
			pool.clients <- client
		}
	}
	return pool
}

// 获取一个代理客户端
func (p *ClientPool) Get() *http.Client {
	return <-p.clients
}

// 将客户端放回池中
func (p *ClientPool) Put(client *http.Client) {
	p.clients <- client
}
func NewProxyClientByUrl(porxyUrl string) *http.Client {
	validURLs := checkProxyURLs([]string{porxyUrl})
	if len(validURLs) < 1 {
		return nil
	}
	client, err := newProxyClient(validURLs[0])
	if err != nil {
		return nil
	}
	return client
}
func newProxyClient(porxyUrl *url.URL) (*http.Client, error) {
	// 创建一个带有代理的 Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(porxyUrl),
	}
	// 创建一个带有自定义 Transport 的 Client
	client := &http.Client{
		Transport: transport,
	}
	return client, nil
} // 创建一个新的代理 HTTP 客户端
func checkProxyURLs(proxyURLs []string) []*url.URL {
	var validURLs []*url.URL
	for _, proxyURL := range proxyURLs {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			log.Printf("无效的 URL: %s, 错误: %v", proxyURL, err)
			continue
		}
		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			log.Printf("无效的 URL: %s, 缺少 scheme 或 host", proxyURL)
			continue
		}
		validURLs = append(validURLs, parsedURL)
	}
	return validURLs
}
