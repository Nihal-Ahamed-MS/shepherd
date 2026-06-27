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

func walkTree(cursor *tree_sitter.TreeCursor, content []byte, filePath string, allowedKinds map[string]bool, codebase *types.CodebaseAST, mut *sync.Mutex) {
	node := cursor.Node()

	if node.IsNamed() && allowedKinds[node.Kind()] {
		chunk := types.Chunk{
			SourceCode: string(content[node.StartByte():node.EndByte()]),
			StartLine:  int(node.StartPosition().Row),
			EndLine:    int(node.EndPosition().Row),
			FilePath:   []string{filePath},
		}
		mut.Lock()
		codebase.Chunks = append(codebase.Chunks, chunk)
		mut.Unlock()
	}

	if cursor.GotoFirstChild() {
		walkTree(cursor, content, filePath, allowedKinds, codebase, mut)
		cursor.GotoParent()
	}

	if cursor.GotoNextSibling() {
		walkTree(cursor, content, filePath, allowedKinds, codebase, mut)
	}
}

func parseFile(currentPath string, wg *sync.WaitGroup, codebase *types.CodebaseAST, mut *sync.Mutex, sem chan struct{}) {
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

	// fmt.Println(string(content))

	parser := parserPool.Get().(*tree_sitter.Parser)
	defer parserPool.Put(parser)

	tree := parser.Parse(content, nil)
	defer tree.Close()

	allowedKinds := make(map[string]bool)
	for _, kind := range utils.NodeTypes[string(lang)] {
		allowedKinds[kind] = true
	}

	cursor := tree.Walk()
	defer cursor.Close()
	walkTree(cursor, content, currentPath, allowedKinds, codebase, mut)
}

func ParseCodebase(rootPath string) (*types.CodebaseAST, error) {
	codebase := types.CodebaseAST{Root: rootPath}

	var (
		wg  sync.WaitGroup
		mut sync.Mutex
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
		go parseFile(currentPath, &wg, &codebase, &mut, sem)

		return nil
	})

	wg.Wait()
	log.Println(len(codebase.Chunks))
	log.Printf("Execution took %s, %s", time.Since(start), codebase.Chunks[5000].SourceCode)

	return &codebase, err
}
