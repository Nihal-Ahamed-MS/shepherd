package types

type Language string

const (
	LangTypeScript Language = "typescript"
	LangJavaScript Language = "javascript"
	LangUnknown    Language = "unknown"
)

type Chunk struct {
	SourceCode string   `json:"sourceCode"`
	StartLine  int      `json:"startLine"`
	EndLine    int      `json:"endLine"`
	FilePath   []string `json:"filePath"`
}

type CodebaseAST struct {
	Root   string  `json:"root"`
	Chunks []Chunk `json:"chunks"`
}

type IndexCodebaseInput struct {
	Path string `json:"path"`
}
