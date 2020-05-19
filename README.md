# goose

<a href=""><img height="190px" src="logo.svg"></a>

**goose** is a lightweight web framework in Go.

<details>
<summary><strong>A hello world example</strong></summary>

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	g.Run(":8080")
}
```

</details>