package commit

type Type string

const (
	TypeFeat     Type = "feat"
	TypeFix      Type = "fix"
	TypeRefactor Type = "refactor"
	TypeChore    Type = "chore"
)

type Model struct {
	Type        Type
	Scope       string
	Description string
	Body        string
}
