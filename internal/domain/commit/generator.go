package commit

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyDescription = errors.New("commit description is required")

func GenerateMessage(model Model) (string, error) {
	description := strings.TrimSpace(model.Description)
	if description == "" {
		return "", ErrEmptyDescription
	}

	commitType := normalizeType(model.Type)
	scope := strings.TrimSpace(model.Scope)
	body := strings.TrimSpace(model.Body)
	header := formatHeader(commitType, scope, description)
	if body == "" {
		return header, nil
	}

	return fmt.Sprintf("%s\n\n%s", header, body), nil
}

func formatHeader(commitType Type, scope string, description string) string {
	if scope == "" {
		return fmt.Sprintf("%s: %s", commitType, description)
	}

	return fmt.Sprintf("%s(%s): %s", commitType, scope, description)
}

func normalizeType(commitType Type) Type {
	switch commitType {
	case TypeFeat, TypeFix, TypeRefactor, TypeChore:
		return commitType
	default:
		return TypeChore
	}
}
