package patterns

import (
	"fmt"

	"github.com/navigator-systems/jrx/ops"
)

func InitTemplates() {
	fmt.Println("Downloading templates...")
	fmt.Println("Reading JRX config toml file...")
	jrxConfig, err := ops.ReadJRXConfig()
	if err != nil {
		fmt.Println("Error reading JRX config:", err)
		return
	}
	fmt.Println("JRX config loaded successfully")
	CloneTemplate(jrxConfig.TemplatesRepo, jrxConfig.TemplatesBranch, jrxConfig.SshKeyPath, jrxConfig.SshKeyPassphrase)

}
