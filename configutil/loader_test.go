package configutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type config struct {
	A string
	B string
}

func TestLoader_Load(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		c := new(config)
		_ = os.Setenv("TEST_B", "baz")
		_ = os.Setenv("B", "bla")

		err := Loader{
			File:      "testdata/config.yml",
			EnvPrefix: "TEST",
		}.Load(c)

		require.NoError(t, err)
		require.Equal(t, "foo", c.A)
		require.Equal(t, "baz", c.B)
	})

	t.Run("no file", func(t *testing.T) {
		c := new(config)
		_ = os.Setenv("TEST2_B", "baz")
		_ = os.Setenv("B", "bla")

		err := Loader{
			EnvPrefix: "TEST2",
		}.Load(c)

		require.NoError(t, err)
		require.Equal(t, "", c.A)
		require.Equal(t, "baz", c.B)
	})
}
