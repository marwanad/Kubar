package kubar

import (
	"fmt"
	"os"

	"github.com/marwanad/kubar/export"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [output directory]\n", os.Args[0])
		return
	}
	path := os.Args[1]
	export.Export(path)
}
