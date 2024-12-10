# pixelcut 的扩图接口
自己用stable diffusion 扩图,生成效果不理想.而且很慢.发现一个可以白嫖的扩图接口



### 用法
```bash
go get github.com/pzx521521/pixelcut
```


单文件
```go
package main
import (
    "github.com/pzx521521/pixelcut"
	"log"
	"net/http"
	"path/filepath"
)

func main()  {
	filePath := "./test.jpg"
	savePath := "./outpaint/test.jpg"
	pixelcut.OutPaintFile(http.DefaultClient, filePath, savePath)
}

```

### 文件夹

```go
package main

import (
	"github.com/pzx521521/pixelcut"
	"log"
	"net/http"
)

func main() {
	dirPath := "/Users/parapeng/Downloads/pinterest"
	err := pixelcut.OutPaintDir(http.DefaultClient, dirPath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}

```

### 多ip多线程
因为有ip限制,使用的[hy2魔改的](https://github.com/pzx521521/hysteria/r)  
hy2本身是不支持同一网址不同ip并发访问的
```go

package main

import (
	"github.com/pzx521521/pixelcut"
	"log"
	"fmt"
)

func main() {
	dirPath := "/Users/parapeng/Downloads/pinterest"
	//本地ip
	proxies := []string{"http://127.0.0.1:8888"}
	//通过其他方式的ip
	for i := 0; i < 6; i++ {
		proxies = append(proxies, fmt.Sprintf("http://127.0.0.1:%d", i+7000))
	}
	//每个ip进行几个并发 这里没有做并发
	pool := pixelcut.NewClientPool(proxies, 1)
	err := pixelcut.OutPaintDirByPool(pool, dirPath)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}	
}
```
效果比自己的sd强:  
原图:  
![原图](https://gd-hbimg.huaban.com/8a755b1b31d1a6fc19187aa993a75b64019fd02d2480d-7inhux_fw658webp)
扩展后:  
![扩展后](https://gd-hbimg.huaban.com/eb99b9509db1d2c9b4166901af5d41e4e126af415ade5-y5qpks_fw658webp)
![扩展后](https://gd-hbimg.huaban.com/da12f7c053ec4b8500c6815e15c7f21fca04c2235343d-39dRid_fw658webp)
![扩展后](https://gd-hbimg.huaban.com/7b9d2321b1b8ebb57287f1fed0984ce63d9ca8f31d00c9-0GhCGV_fw658webp)
![扩展后](https://gd-hbimg.huaban.com/d98b6c090a0de28d52c53a3cd1b5a737224a59951c947d-YkadJp_fw658webp)
[更多扩图的效果](https://paral.us.kg/)可以看一下
