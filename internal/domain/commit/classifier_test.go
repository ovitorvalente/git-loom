package commit

import "testing"

func TestClassifyCommit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		diff     string
		expected Type
	}{
		{
			name:     "returns feat for added behavior",
			diff:     "add support for branch automation",
			expected: TypeFeat,
		},
		{
			name:     "returns fix for bug related diff",
			diff:     "fix timeout error in git command",
			expected: TypeFix,
		},
		{
			name:     "returns refactor for internal cleanup",
			diff:     "refactor commit generator to simplify logic",
			expected: TypeRefactor,
		},
		{
			name:     "returns chore for unmatched diff",
			diff:     "update dependency metadata",
			expected: TypeChore,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := ClassifyCommit(testCase.diff)
			if result != testCase.expected {
				t.Fatalf("expected %s, got %s", testCase.expected, result)
			}
		})
	}
}
