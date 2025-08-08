package patterns

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

func ListTemplates() {
	x, err := GetTemplateCtrl()
	if err != nil {
		fmt.Println("Error getting templates:", err)
		return
	}
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.Indent = "  "
	if err := encoder.Encode(x); err != nil {
		fmt.Println("Error encoding templates:", err)
		return
	}

	fmt.Println(buf.String())
}
