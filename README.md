# ipmidump

Wrapper library for the [u-root](https://github.com/u-root/u-root) project's [ipmi package](https://github.com/u-root/u-root/tree/main/pkg/ipmi)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jak3kaj/ipmidump"
)

func main() {
	dump := ipmidump.Dump()
	if json, err := json.MarshalIndent(dump, "", "  "); err == nil {
		fmt.Println(string(json))
	}
}
```
