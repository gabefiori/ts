package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generatePath(dirs map[string]struct{}, elem ...string) {
	dirs[filepath.Join(elem...)] = struct{}{}
}

func TestFind(t *testing.T) {
	tempDir := t.TempDir()
	tempDir2 := t.TempDir()

	dirs := make(map[string]struct{})

	generatePath(dirs, tempDir, "sub_dir_1")
	generatePath(dirs, tempDir, "sub_dir_2")
	generatePath(dirs, tempDir2, "sub_dir_3", "sub_dir_4")
	generatePath(dirs, tempDir2, "sub_dir_3", "sub_dir_5")

	for dir := range dirs {
		assert.NoError(t, os.MkdirAll(dir, 0755))
		delete(dirs, dir)

		newDir := strings.Replace(dir, tempDir2, tempDir, 1)
		dirs[strings.Replace(newDir, "sub_dir_3", "symlink", 1)] = struct{}{}
	}

	generatePath(dirs, tempDir)

	symlinkPath := filepath.Join(tempDir, "symlink")
	assert.NoError(t, os.Symlink(filepath.Join(tempDir2, "sub_dir_3"), symlinkPath))

	generatePath(dirs, symlinkPath)

	// Depth 0
	foundDirs, err := Find(tempDir, 0)
	assert.NoError(t, err)


	assert.Len(t, foundDirs, 1)

	for _, v := range foundDirs {
		_, ok := dirs[v]
		assert.True(t, ok)
	}

	// Depth 1
	foundDirs, err = Find(tempDir, 1)
	assert.NoError(t, err)

	assert.Len(t, foundDirs, 4)

	for _, v := range foundDirs {
		_, ok := dirs[v]
		assert.True(t, ok)
	}

	// Full depth with symlinks
	foundDirs, err = Find(tempDir, 3)
	assert.NoError(t, err)

	assert.Len(t, foundDirs, len(dirs))

	for _, v := range foundDirs {
		_, ok := dirs[v]
		assert.True(t, ok)
	}
}

func BenchmarkFind(b *testing.B) {
	tempDir := b.TempDir()
	tempDir2 := b.TempDir()

	dirs := make(map[string]struct{})

	for i := 0; i < 100; i++ {
		generatePath(dirs, tempDir, "sub_dir_1", fmt.Sprintf("sub_dir_%v", i))
		generatePath(dirs, tempDir2, "sub_dir_3", fmt.Sprintf("sub_dir_%v", i))
	}

	for dir := range dirs {
		_ = os.MkdirAll(dir, 0755)
	}

	symlinkPath := filepath.Join(tempDir, "symlink")
	_ = os.Symlink(filepath.Join(tempDir2, "sub_dir_3"), symlinkPath)

	generatePath(dirs, symlinkPath)

	b.Run("Depth 0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Depth 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Full Depth", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 3)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
