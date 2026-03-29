package semantic

import (
	"path/filepath"
	"sort"
	"strings"
)

type ChangedFile struct {
	Path   string
	Status string
}

type CommitContext struct {
	Files []ChangedFile
	Diff  string
	Tags  []string
}

type ChangeIntent struct {
	Type        string
	Scope       string
	Intent      string
	Description string
}

type CommitPreview struct {
	FilesChanged int
	Additions    int
	Deletions    int
}

type CommitQuality struct {
	Score   int
	Reasons []string
}

func NewCommitContext(diff string) CommitContext {
	files := extractChangedFiles(diff)
	return CommitContext{
		Files: files,
		Diff:  diff,
		Tags:  detectTags(diff, files),
	}
}

func BuildPreview(context CommitContext) CommitPreview {
	additions := 0
	deletions := 0

	for _, line := range strings.Split(context.Diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+++"), strings.HasPrefix(line, "---"):
			continue
		case strings.HasPrefix(line, "+"):
			additions++
		case strings.HasPrefix(line, "-"):
			deletions++
		}
	}

	return CommitPreview{
		FilesChanged: len(context.Files),
		Additions:    additions,
		Deletions:    deletions,
	}
}

func BuildGroupingKey(commitType string, context CommitContext) string {
	return strings.Join([]string{
		commitType,
		NormalizeScopeFromFiles(context.Files),
		normalizePackageFamily(context.Files),
	}, "|")
}

func extractChangedFiles(diff string) []ChangedFile {
	uniqueFiles := map[string]ChangedFile{}

	lines := strings.Split(diff, "\n")
	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}

		path := extractPath(line)
		if path == "" {
			continue
		}

		uniqueFiles[path] = ChangedFile{
			Path:   path,
			Status: detectStatus(lines, index),
		}
	}

	paths := make([]string, 0, len(uniqueFiles))
	for path := range uniqueFiles {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	files := make([]ChangedFile, 0, len(paths))
	for _, path := range paths {
		files = append(files, uniqueFiles[path])
	}

	return files
}

func extractPath(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return ""
	}

	return strings.TrimPrefix(fields[3], "b/")
}

func detectStatus(lines []string, start int) string {
	if start+1 >= len(lines) {
		return "atualizado"
	}

	nextLine := strings.TrimSpace(lines[start+1])
	switch {
	case strings.HasPrefix(nextLine, "new file mode"):
		return "adicionado"
	case strings.HasPrefix(nextLine, "deleted file mode"):
		return "removido"
	default:
		return "atualizado"
	}
}

func detectTags(diff string, files []ChangedFile) []string {
	tagSet := map[string]bool{}
	normalizedDiff := strings.ToLower(diff)

	for _, file := range files {
		baseName := strings.ToLower(filepath.Base(file.Path))
		name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		for _, part := range strings.FieldsFunc(name, splitToken) {
			if len(part) < 3 {
				continue
			}
			tagSet[part] = true
		}
	}

	for _, keyword := range []string{
		"commit", "cli", "preview", "strict", "score", "suggest", "prompt",
		"output", "readme", "deps", "build", "config", "git", "test", "ui",
	} {
		if strings.Contains(normalizedDiff, keyword) {
			tagSet[keyword] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func normalizePackageFamily(files []ChangedFile) string {
	if len(files) == 0 {
		return "root"
	}

	votes := map[string]int{}
	for _, file := range files {
		votes[packageFamily(file.Path)]++
	}

	return mostCommon(votes)
}

func packageFamily(path string) string {
	directory := filepath.Dir(path)
	if directory == "." || directory == "" {
		return NormalizeScope(path)
	}

	segments := strings.Split(directory, "/")
	if len(segments) >= 3 && segments[0] == "internal" {
		return strings.Join(segments[:3], "/")
	}
	if len(segments) >= 2 {
		return strings.Join(segments[:2], "/")
	}

	return directory
}

func mostCommon(votes map[string]int) string {
	selectedValue := ""
	selectedVotes := 0

	for value, voteCount := range votes {
		if voteCount > selectedVotes || (voteCount == selectedVotes && value < selectedValue) {
			selectedValue = value
			selectedVotes = voteCount
		}
	}

	return selectedValue
}

func splitToken(r rune) bool {
	return r == '_' || r == '-' || r == '.' || r == '/' || r == ' '
}
