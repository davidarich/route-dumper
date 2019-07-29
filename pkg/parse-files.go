package pkg

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"io/ioutil"
	"os"
)

func ParseFiles(files []string) []string {
	for _, file := range files {
		showFileList := false
		if showFileList {
			fmt.Println(file)
		}

		src, _ := ioutil.ReadFile(file)

		parser := php7.NewParser(src)
		parser.Parse()
		for _, e := range parser.GetErrors() {
			fmt.Println(e)
		}
		fmt.Println(src)

		nsResolver := visitor.NewNamespaceResolver()

		annotationRouteDumper := AnnotationRouteDumper{
			Writer:     os.Stdout,
			NsResolver: nsResolver,
		}

		rootNode := parser.GetRootNode()
		rootNode.Walk(&annotationRouteDumper)

	}
	return files
}
