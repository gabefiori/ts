package sessionizer

import (
	"bytes"
	"errors"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/gabefiori/gotmux"
	"github.com/gabefiori/ts/config"
	"github.com/gabefiori/ts/internal/errutil"
	"github.com/gabefiori/ts/internal/selector"
	"github.com/gabefiori/ts/internal/targets"
)

func Run(cfg *config.Config) error {
	allTargets, err := generateTargets(cfg.Targets)

	if err != nil {
		return errutil.NewError(errutil.SessionizerErr, err)
	}

	if cfg.List {
		if cfg.Filter != "" {
			filterTargets(cfg.Filter, &allTargets)
		}

		PrintList(allTargets)
		return nil
	}

	if cfg.Filter != "" {
		// In this case, we delegate the responsibility of filtering to the selector.
		// This way, we avoid losing any targets.
		cfg.Selector = append(cfg.Selector, "--query="+cfg.Filter)
	}

	selected, err := selector.Run(allTargets, cfg.Selector)

	if err != nil {
		return errutil.NewError(errutil.SelectorErr, err)
	}

	if selected == "" {
		return nil
	}

	if err := runTmux(selected); err != nil {
		return errutil.NewError(errutil.TmuxErr, err)
	}

	return nil
}

func RunSingle(target string) error {
	if err := targets.FindSingle(target); err != nil {
		return errutil.NewError(errutil.SessionizerErr, err)
	}

	if err := runTmux(target); err != nil {
		return errutil.NewError(errutil.TmuxErr, err)
	}

	return nil
}

func PrintList(allTargets []string) {
	var result []byte

	for i, target := range allTargets {
		result = append(result, target...)

		if i < len(allTargets)-1 {
			result = append(result, '\n')
		}
	}

	io.Copy(os.Stdout, bytes.NewBuffer(result))
}

func filterTargets(filter string, targets *[]string) {
	filtered := []string{}

	for _, t := range *targets {
		if strings.Contains(t, filter) {
			filtered = append(filtered, t)
		}
	}

	*targets = filtered
}

func runTmux(target string) error {
	sessionName := strings.TrimPrefix(filepath.Base(target), ".")

	if gotmux.HasSession(sessionName) {
		if err := gotmux.AttachOrSwitchTo(sessionName); err != nil {
			return err
		}

		return nil
	}

	session, err := gotmux.NewSession(&gotmux.SessionConfig{
		Name: sessionName,
		Dir:  target,
	})

	if err != nil {
		return err
	}

	if err := session.AttachOrSwitch(); err != nil {
		return err
	}

	return nil
}

func generateTargets(configTargets []config.Target) ([]string, error) {
	targetCh := make(chan string, len(configTargets)*10)
	errCh := make(chan error, len(configTargets))

	var wg sync.WaitGroup

	for _, configTarget := range configTargets {
		wg.Add(1)

		go func(ct config.Target) {
			defer wg.Done()
			foundTargets, err := targets.Find(ct.Path, ct.Depth)

			if err != nil {
				errCh <- err
				return
			}

			for _, fd := range foundTargets {
				targetCh <- fd
			}

		}(configTarget)
	}

	go func() {
		wg.Wait()

		close(targetCh)
		close(errCh)
	}()

	targetMap := make(map[string]struct{})
	var errs []error

	for targetCh != nil && errCh != nil {
		select {
		case target, ok := <-targetCh:
			if !ok {
				targetCh = nil
				break
			}

			targetMap[target] = struct{}{}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				break
			}

			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	allTargets := slices.Collect(maps.Keys(targetMap))

	sort.Slice(allTargets, func(i, j int) bool {
		return allTargets[i] > allTargets[j]
	})

	return allTargets, nil
}
