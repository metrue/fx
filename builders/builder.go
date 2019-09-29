package builders

// Builder image builder
type Builder interface {
	Build(workdir string, name string) error
}
