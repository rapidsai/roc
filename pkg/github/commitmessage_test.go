package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveHTMLComments(t *testing.T) {
	examples := []string{
		"hello world<!-- test  -->",
		"hello world<!-- test  -->more words",
		"hello world<!-- test",
		`my fancy commit
message on multiple lines
<!--
please don't forget to sign the CLA
-->`,
	}
	expected := []string{
		"hello world",
		"hello worldmore words",
		"hello world<!-- test",
		`my fancy commit
message on multiple lines
`,
	}
	for i, example := range examples {
		uncommented := RemoveHTMLComments(example)
		require.Equal(t, uncommented, expected[i])
	}
}
