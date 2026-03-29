package mocks

type GitRepository struct {
	GetDiffFunc          func(paths ...string) (string, error)
	IsRepositoryFunc     func() (bool, error)
	ListStagedFilesFunc  func() ([]string, error)
	ListChangedFilesFunc func() ([]string, error)
	StageFilesFunc       func(paths []string) error
	CommitFunc           func(message string) error
	CommitPathsFunc      func(message string, paths []string) error
	CreateBranchFunc     func(name string) error

	GetDiffCalls      [][]string
	StageFilesCalls   [][]string
	CommitCalls       []string
	CommitPathsCalls  []CommitCall
	CreateBranchCalls []string
}

type CommitCall struct {
	Message string
	Paths   []string
}

func (mock *GitRepository) GetDiff(paths ...string) (string, error) {
	mock.GetDiffCalls = append(mock.GetDiffCalls, append([]string(nil), paths...))
	if mock.GetDiffFunc == nil {
		return "", nil
	}

	return mock.GetDiffFunc(paths...)
}

func (mock *GitRepository) IsRepository() (bool, error) {
	if mock.IsRepositoryFunc == nil {
		return true, nil
	}

	return mock.IsRepositoryFunc()
}

func (mock *GitRepository) ListStagedFiles() ([]string, error) {
	if mock.ListStagedFilesFunc == nil {
		return nil, nil
	}

	return mock.ListStagedFilesFunc()
}

func (mock *GitRepository) ListChangedFiles() ([]string, error) {
	if mock.ListChangedFilesFunc == nil {
		return nil, nil
	}

	return mock.ListChangedFilesFunc()
}

func (mock *GitRepository) StageFiles(paths []string) error {
	mock.StageFilesCalls = append(mock.StageFilesCalls, append([]string(nil), paths...))
	if mock.StageFilesFunc == nil {
		return nil
	}

	return mock.StageFilesFunc(paths)
}

func (mock *GitRepository) Commit(message string) error {
	mock.CommitCalls = append(mock.CommitCalls, message)
	if mock.CommitFunc == nil {
		return nil
	}

	return mock.CommitFunc(message)
}

func (mock *GitRepository) CommitPaths(message string, paths []string) error {
	mock.CommitPathsCalls = append(mock.CommitPathsCalls, CommitCall{
		Message: message,
		Paths:   append([]string(nil), paths...),
	})
	if mock.CommitPathsFunc == nil {
		return nil
	}

	return mock.CommitPathsFunc(message, paths)
}

func (mock *GitRepository) CreateBranch(name string) error {
	mock.CreateBranchCalls = append(mock.CreateBranchCalls, name)
	if mock.CreateBranchFunc == nil {
		return nil
	}

	return mock.CreateBranchFunc(name)
}
