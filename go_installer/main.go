package main

import (
	"os"

	// "github.com/flocknroll/azga/go_installer/add_txt_content_closure"
	"github.com/flocknroll/azga/go_installer/add_txt_content_go"
)

func main() {
	// fmt.Println(checkContent("test_data/src_found.txt", "test_data/dest.txt"))
	// fmt.Println(checkContent("test_data/src_not_found.txt", "test_data/dest.txt"))

	add_txt_content_go.AddContent(os.Args[1], os.Args[2])
}
