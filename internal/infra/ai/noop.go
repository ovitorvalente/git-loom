package ai

type NoopProvider struct{}

func NewNoopProvider() NoopProvider {
	return NoopProvider{}
}

func (NoopProvider) GenerateCommit(diff string) (string, error) {
	return "", nil
}
