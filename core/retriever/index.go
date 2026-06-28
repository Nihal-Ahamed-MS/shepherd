package retriever

import (
	"fmt"
	"runtime"
	"sync"

	"shepherd/core/types"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/token/camelcase"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
)

const (
	codeAnalyzer = "code"
	batchSize    = 100
)

func buildMapping() (mapping.IndexMapping, error) {
	mapping := bleve.NewIndexMapping()

	err := mapping.AddCustomAnalyzer(codeAnalyzer, map[string]interface{}{
		"type":          custom.Name,
		"tokenizer":     unicode.Name,
		"token_filters": []string{camelcase.Name, lowercase.Name},
	})
	if err != nil {
		return nil, err
	}

	sourceField := bleve.NewTextFieldMapping()
	sourceField.Analyzer = codeAnalyzer

	fileField := bleve.NewTextFieldMapping()
	fileField.Analyzer = "keyword"

	chunkMapping := bleve.NewDocumentMapping()
	chunkMapping.AddFieldMappingsAt("sourceCode", sourceField)
	chunkMapping.AddFieldMappingsAt("filePath", fileField)

	mapping.AddDocumentMapping("chunk", chunkMapping)
	mapping.DefaultAnalyzer = codeAnalyzer

	return mapping, nil
}

func indexBatch(index bleve.Index, batch []types.Chunk, offset int, mu *sync.Mutex, indexErr *error, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	b := index.NewBatch()
	for j, chunk := range batch {
		id := fmt.Sprintf("%d", offset+j)
		if err := b.Index(id, chunk.SourceCode); err != nil {
			mu.Lock()
			*indexErr = fmt.Errorf("index chunk %s: %w", id, err)
			mu.Unlock()
			return
		}
	}

	if err := index.Batch(b); err != nil {
		mu.Lock()
		*indexErr = fmt.Errorf("flush batch at %d: %w", offset, err)
		mu.Unlock()
	}
}

func BuildIndex(chunks []types.Chunk) (bleve.Index, error) {
	mapping, err := buildMapping()
	if err != nil {
		return nil, fmt.Errorf("build mapping: %w", err)
	}

	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, fmt.Errorf("create index: %w", err)
	}

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		indexErr error
		sem      = make(chan struct{}, runtime.NumCPU())
	)

	for i := 0; i < len(chunks); i += batchSize {
		end := min(i+batchSize, len(chunks))
		batch := chunks[i:end]
		offset := i

		wg.Add(1)
		sem <- struct{}{}
		go indexBatch(index, batch, offset, &mu, &indexErr, &wg, sem)
	}

	wg.Wait()
	if indexErr != nil {
		return nil, indexErr
	}

	return index, nil
}
