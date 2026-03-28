package commit

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyDescription = errors.New("a descricao do commit e obrigatoria")

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
	description = limitHeaderDescription(commitType, scope, description)
	if scope == "" {
		return fmt.Sprintf("%s: %s", commitType, description)
	}

	return fmt.Sprintf("%s(%s): %s", commitType, scope, description)
}

func normalizeType(commitType Type) Type {
	switch commitType {
	case TypeFeat, TypeFix, TypeRefactor, TypeChore, TypeDocs, TypeTest:
		return commitType
	default:
		return TypeChore
	}
}

func limitHeaderDescription(commitType Type, scope string, description string) string {
	const headerLimit = 72

	prefixLength := len(commitType) + 2
	if scope != "" {
		prefixLength += len(scope) + 2
	}

	maxDescriptionLength := headerLimit - prefixLength
	if maxDescriptionLength <= 0 || len(description) <= maxDescriptionLength {
		return description
	}

	truncatedDescription := strings.TrimSpace(description[:maxDescriptionLength])
	lastSpaceIndex := strings.LastIndex(truncatedDescription, " ")
	if lastSpaceIndex <= 0 {
		return truncatedDescription
	}

	return strings.TrimSpace(truncatedDescription[:lastSpaceIndex])
}
