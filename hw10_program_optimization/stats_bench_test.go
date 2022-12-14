package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	require.Equal(b, 1, len(r.File))

	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		require.NoError(b, err)
		_, err = GetDomainStat(data, "biz")
		require.NoError(b, err)
		err = data.Close()
		require.NoError(b, err)
	}
	err = r.Close()
	require.NoError(b, err)
}
