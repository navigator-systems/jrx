package cmd

import (
	"fmt"
	"jrx/ops"
)

func InfoCmd(path string, osv bool) {
	ops.GetFileInfo(path)
	x, _ := ops.ReadGoSum(path)
	if osv {
		fmt.Println("Checking for vulnerabilities...")
		ops.CheckVulnerabilities(x)
	}

}
