package mocks

type GitRepository struct {
	GetDiffFunc      func() (string, error)
	CommitFunc       func(message string) error
	CreateBranchFunc func(name string) error

	CommitCalls       []string
	CreateBranchCalls []string
}

func (mock *GitRepository) GetDiff() (string, error) {
	if mock.GetDiffFunc == nil {
		return "", nil
	}

	return mock.GetDiffFunc()
}

func (mock *GitRepository) Commit(message string) error {
	mock.CommitCalls = append(mock.CommitCalls, message)
	if mock.CommitFunc == nil {
		return nil
	}

	return mock.CommitFunc(message)
}

func (mock *GitRepository) CreateBranch(name string) error {
	mock.CreateBranchCalls = append(mock.CreateBranchCalls, name)
	if mock.CreateBranchFunc == nil {
		return nil
	}

	return mock.CreateBranchFunc(name)
}
