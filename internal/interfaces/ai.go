package interfaces

type AIProvider interface {
	GenerateCommit(diff string) (string, error)
}
