package interfaces

type GitRepository interface {
	GetDiff(paths ...string) (string, error)
	IsRepository() (bool, error)
	ListStagedFiles() ([]string, error)
	ListChangedFiles() ([]string, error)
	StageFiles(paths []string) error
	Commit(message string) error
	CommitPaths(message string, paths []string) error
	CreateBranch(name string) error
}
