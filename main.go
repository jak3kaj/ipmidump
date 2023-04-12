package main

import (
	"encoding/json"
	"fmt"
	"github.com/jak3kaj/ipmidump"
)

func main() {
	dump := ipmi_dump.Dump()
	if json, err := json.MarshalIndent(dump, "", "  "); err == nil {
		fmt.Println(string(json))
	}
}
