package selector

import (
	"maps"
	"strings"

	fzf "github.com/junegunn/fzf/src"
)

func Run(items []string, opts []string) (string, error) {
	input := make(chan string)
	output := make(chan string)

	var result string

	go func() {
		for _, i := range items {
			input <- i
		}

		close(input)
	}()

	go func() {
		for o := range output {
			result = o
		}
	}()

	options, err := fzf.ParseOptions(true, mergeConfigs(opts))

	if err != nil {
		return "", err
	}

	options.Input = input
	options.Output = output

	_, err = fzf.Run(options)

	close(output)

	return result, err
}

func mergeConfigs(config []string) []string {
	cm := configToMap(DefaultOptions())
	maps.Copy(cm, configToMap(config))

	result := make([]string, 0, len(cm))

	for k, v := range cm {
		cfg := k

		if v != "" {
			cfg = cfg + "=" + v
		}

		result = append(result, cfg)
	}

	return result
}

func configToMap(config []string) map[string]string {
	cm := make(map[string]string)

	for _, c := range config {
		splitted := strings.SplitN(c, "=", 2)
		val := ""

		if len(splitted) > 1 {
			val = splitted[1]
		}

		cm[splitted[0]] = val
	}

	return cm
}

func DefaultOptions() []string {
	return []string{
		"--border",
		"--border-label= Tmux Sessionizer ",
	}
}
