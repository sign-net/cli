package xfilepath_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

func TestJoin(t *testing.T) {
	retriever := xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	path, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), path)

	retriever = xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}

func TestJoinFromHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	retriever := xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	path, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		home,
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), path)

	retriever = xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}

func TestList(t *testing.T) {
	retriever := xfilepath.List()
	list, err := retriever()
	require.NoError(t, err)
	require.Equal(t, []string(nil), list)

	retriever1 := xfilepath.Join(
		xfilepath.Path("foo/bar"),
	)
	retriever2 := xfilepath.Join(
		xfilepath.Path("bar/foo"),
	)
	retriever = xfilepath.List(retriever1, retriever2)
	list, err = retriever()
	require.NoError(t, err)
	require.Equal(t, []string{
		filepath.Join("foo", "bar"),
		filepath.Join("bar", "foo"),
	}, list)

	retrieverError := xfilepath.PathWithError("foo", errors.New("foo"))
	retriever = xfilepath.List(retriever1, retrieverError, retriever2)
	_, err = retriever()
	require.Error(t, err)
}

func TestIsDirEmpty(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want bool
		err  error
	}{
		{
			name: "current dir",
			dir:  "./",
			want: false,
		},
		{
			name: "temp dir",
			dir:  t.TempDir(),
			want: true,
		},
		{
			name: "parent dir",
			dir:  "../",
			want: false,
		},
		{
			name: "invalid dir",
			dir:  "())*",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xfilepath.IsDirEmpty(tt.dir)
			require.Equal(t, tt.want, got)
		})
	}
}
