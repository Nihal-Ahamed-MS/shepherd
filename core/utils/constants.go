package utils

// TODO: load the data from .gitignore
var SkipDirs = map[string]bool{
	".git": true, "node_modules": true, ".next": true, "build": true,
	"vendor": true, "dist": true, ".cache": true, "bin": true,
}

var NodeTypes = map[string][]string{
	"javascript": {
		"function_declaration",
		"method_definition",
		"class_declaration",
		"import_statement",
		"export_statement",
		"lexical_declaration",
	},
	"typescript": {
		"function_declaration",
		"method_definition",
		"class_declaration",
		"import_statement",
		"export_statement",
		"lexical_declaration",
		"interface_declaration",
		"type_alias_declaration",
		"enum_declaration",
	},
}