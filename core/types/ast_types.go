package types

type Language string

const (
	LangTypeScript Language = "typescript"
	LangJavaScript Language = "javascript"
	LangUnknown    Language = "unknown"
)

type Symbol struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Line    int    `json:"line"`
	EndLine int    `json:"end_line,omitempty"`
	Doc     string `json:"doc,omitempty"`
}

type FileNode struct {
	Path     string   `json:"path"`
	Language Language `json:"language"`
	Lines    int      `json:"lines"`
	Symbols  []Symbol `json:"symbols,omitempty"`
}

type CodebaseAST struct {
	Root  string     `json:"root"`
	Files []FileNode `json:"files"`
}

type IndexCodebaseInput struct {
	Path string `json:"path"`
}
