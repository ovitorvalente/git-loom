package mocks

type AIProvider struct {
	GenerateCommitFunc  func(diff string) (string, error)
	GenerateCommitCalls []string
}

func (mock *AIProvider) GenerateCommit(diff string) (string, error) {
	mock.GenerateCommitCalls = append(mock.GenerateCommitCalls, diff)
	if mock.GenerateCommitFunc == nil {
		return "", nil
	}

	return mock.GenerateCommitFunc(diff)
}
