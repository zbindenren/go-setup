package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	var tt = []struct {
		name        string
		gopath      string
		packagePath string
		src         string
		err         error
	}{
		{
			"valid",
			"/tmp/go-setup/test",
			"github.com/zbindenren/go-setup",
			".",
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := setup(tc.gopath, tc.packagePath, tc.src)
			if tc.err != nil {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.True(t, len(msg) > 0)
			fmt.Println(msg)
			dst, err := os.Readlink("/tmp/go-setup/test/src/github.com/zbindenren")
			require.NoError(t, err)
			assert.Equal(t, path.Join(os.Getenv("GOPATH"), "src/github.com/zbindenren"), dst)
			// create again
			_, err = os.Readlink("/tmp/go-setup/test/src/github.com/zbindenren")
			require.NoError(t, err)
			err = cleanUp(tc.gopath, tc.packagePath)
			assert.NoError(t, err)
		})
	}
}
