package interfaces

type GitRepository interface {
	GetDiff() (string, error)
	Commit(message string) error
	CreateBranch(name string) error
}
