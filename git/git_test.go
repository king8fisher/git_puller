package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitAPI(t *testing.T) {
	commit, err := LatestCommit()
	require.Nil(t, err)
	require.NotEmpty(t, commit)
}
