package cli

import (
	"os"
	"runtime/debug"

	"github.com/gabefiori/ts/config"
	"github.com/gabefiori/ts/internal/sessionizer"
	"github.com/gabefiori/ts/internal/utils"
	"github.com/urfave/cli/v2"
)

func Run() error {
	var path string
	var filter string
	var target string
	var list bool

	app := &cli.App{
		Name:    "ts",
		Usage:   "Tmux Sessionizer is a tool for navigating through folders and projects as tmux sessions.",
		Version: getVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `file`",
				Value:       "~/.config/ts/config.json",
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "filter",
				Aliases:     []string{"f"},
				Usage:       "Specify a filter to narrow down the results displayed in the selector",
				Value:       "",
				Destination: &filter,
			},
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"t"},
				Usage:       "Specify a target (e.g., path) to switch or attach to",
				Value:       "",
				Destination: &target,
			},
			&cli.BoolFlag{
				Name:        "list",
				Aliases:     []string{"l"},
				Usage:       "List of all discovered targets",
				Value:       false,
				Destination: &list,
			},
		},
		Action: func(ctx *cli.Context) error {
			if target != "" {
				return sessionizer.RunSingle(target)
			}

			cliCfg := config.Cli{
				Path:   path,
				Filter: filter,
				List:   list,
			}

			cfg, err := config.Load(cliCfg)

			if err != nil {
				return utils.NewErrorWithPrefix("Config", err)
			}

			return sessionizer.Run(cfg)
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return "unknown"
}
