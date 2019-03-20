package types

// ProjectSourceFile a source code file of a project
type ProjectSourceFile struct {
	Path      string `json:"path"`
	Body      string `json:"body"`
	IsHandler bool   `json:"is_handler"`
}

// Project a project
type Project struct {
	Name     string              `json:"name"`
	Language string              `json:"language"`
	Files    []ProjectSourceFile `json:"files"`
}
