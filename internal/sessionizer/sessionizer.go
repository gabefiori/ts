package sessionizer

import (
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/gabefiori/gotmux"
	"github.com/gabefiori/ts/config"
	"github.com/gabefiori/ts/internal/selector"
	"github.com/gabefiori/ts/internal/targets"
	"github.com/gabefiori/ts/internal/utils"
)

func Run(cfg *config.Config) error {
	allTargets, err := generateTargets(cfg.Targets)

	if err != nil {
		return utils.NewErrorWithPrefix("Tmux Sessionizer", err)
	}

	if cfg.List {
		if cfg.Filter != "" {
			filterTargets(cfg.Filter, &allTargets)
		}

		fmt.Println(strings.Join(allTargets, "\n"))
		return nil
	}

	if cfg.Filter != "" {
		// In this case, we delegate the responsibility of filtering to the selector.
		// This way, we avoid losing any targets.
		cfg.Selector = append(cfg.Selector, "--query="+cfg.Filter)
	}

	selected, err := selector.Run(allTargets, cfg.Selector)

	if err != nil {
		return utils.NewErrorWithPrefix("Selector", err)
	}

	if selected == "" {
		return nil
	}

	if err := runTmux(selected); err != nil {
		return utils.NewErrorWithPrefix("Tmux", err)
	}

	return nil
}

func RunSingle(target string) error {
	if err := targets.FindSingle(target); err != nil {
		return utils.NewErrorWithPrefix("Tmux Sessionizer", err)
	}

	if err := runTmux(target); err != nil {
		return utils.NewErrorWithPrefix("Tmux", err)
	}

	return nil
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
	var (
		targetMap = make(map[string]struct{})
		errs      []error
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for configTarget := range slices.Values(configTargets) {
		wg.Add(1)

		go func(ct config.Target) {
			defer wg.Done()

			foundTargets, err := targets.Find(ct.Path, ct.Depth)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}

			for fd := range slices.Values(foundTargets) {
				targetMap[fd] = struct{}{}
			}

		}(configTarget)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	allTargets := slices.Collect(maps.Keys(targetMap))

	sort.Slice(allTargets, func(i, j int) bool {
		return allTargets[i] > allTargets[j]
	})

	return allTargets, nil
}
