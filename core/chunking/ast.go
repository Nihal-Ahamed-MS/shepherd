package chunking

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"shepherd/core/types"
	"shepherd/core/utils"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tsjs "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var (
	tsLanguage = tree_sitter.NewLanguage(tsjs.Language())
	parserPool = sync.Pool{
		New: func() any {
			p := tree_sitter.NewParser()
			p.SetLanguage(tsLanguage)
			return p
		},
	}
)

func parseFile(currentPath, rootPath string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	lang := utils.DetectLanguage(currentPath)
	if lang == types.LangUnknown {
		return
	}

	content, err := os.ReadFile(currentPath)
	if err != nil {
		return
	}

	parser := parserPool.Get().(*tree_sitter.Parser)
	defer parserPool.Put(parser)

	tree := parser.Parse(content, nil)
	defer tree.Close()

	// rootNode := tree.RootNode()
	// fmt.Println("kind", rootNode.Child(1))
	// fmt.Println("startpos", rootNode.StartPosition())
	// fmt.Println("endpos", rootNode.EndPosition())
}

func ParseCodebase(rootPath string) (*types.CodebaseAST, error) {
	codebase := &types.CodebaseAST{Root: rootPath}

	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, runtime.NumCPU())
	)

	start := time.Now()
	err := filepath.WalkDir(rootPath, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if utils.SkipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}
		go parseFile(currentPath, rootPath, &wg, sem)

		return nil
	})

	wg.Wait()
	log.Printf("Execution took %s", time.Since(start))

	return codebase, err
}
