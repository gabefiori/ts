package targets

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charlievieth/fastwalk"
	hd "github.com/mitchellh/go-homedir"
)

func Find(rootDir string, maxDepth uint8) ([]string, error) {
	homeDir, err := hd.Dir()

	if err != nil {
		return nil, err
	}

	rootDir, err = hd.Expand(rootDir)

	if err != nil {
		return nil, err
	}

	isDir, err := isDirectory(rootDir)

	if err != nil {
		return nil, err
	}

	// PERF: Depth 0 and depth 1 are handled separately to improve performance (approximately 90% faster).
	// Depth 0
	if maxDepth == 0 && isDir {
		return []string{unexpandPath(rootDir, homeDir)}, nil
	}

	// Depth 1
	if maxDepth == 1 && isDir {
		return findDepthOne(rootDir, homeDir)
	}

	// Depth > 1
	return findGreaterDepth(rootDir, homeDir, maxDepth)
}

func findDepthOne(rootDir, homeDir string) ([]string, error) {
	entries, err := os.ReadDir(rootDir)

	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(entries)+1)

	for _, entry := range entries {
		entryPath := filepath.Join(rootDir, entry.Name())
		info, err := os.Stat(entryPath)

		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			paths = append(paths, unexpandPath(entryPath, homeDir))
		}
	}

	paths = append(paths, unexpandPath(rootDir, homeDir))
	return paths, nil
}

func findGreaterDepth(rootDir, homeDir string, maxDepth uint8) ([]string, error) {
	errStopWalking := errors.New("stop walking")

	rootDir = strings.TrimSuffix(rootDir, "/")
	paths := make([]string, 0)

	cfg := fastwalk.Config{Follow: true}

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if currentDepth(path, rootDir) > maxDepth {
			return errStopWalking
		}

		if err != nil {
			return err
		}

		paths = append(paths, unexpandPath(path, homeDir))
		return nil
	}

	err := fastwalk.Walk(&cfg, rootDir, fastwalk.IgnoreDuplicateDirs(walkFn))

	if err != nil && err != errStopWalking {
		return nil, err
	}

	return paths, nil
}

func FindSingle(path string) error {
	target, err := hd.Expand(path)

	if err != nil {
		return err
	}

	isValidTarget, err := isDirectory(target)

	if err != nil {
		return err
	}

	if !isValidTarget {
		return errors.New("Invalid target. Not a directory.")
	}

	return nil
}

func unexpandPath(path, homeDir string) string {
	return strings.Replace(path, homeDir, "~", 1)
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, nil
}

func currentDepth(path string, rootDir string) uint8 {
	if path == rootDir {
		return 0
	}

	subPath := strings.Replace(path, rootDir, "", 1)
	return uint8(strings.Count(subPath, "/"))
}
