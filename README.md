# Tmux Sessionizer

Tmux Sessionizer is a tool for navigating through folders and projects as tmux sessions.
Inspired by ThePrimeagen's [tmux-sessionizer](https://github.com/ThePrimeagen/.dotfiles/blob/master/bin/.local/scripts/tmux-sessionizer), this version has been modified to fit my preferences.

## Features
- **navigate directories**: switch between directories and projects as tmux sessions.
- **custom configuration**: configure targets and session options through a JSON file.
- **fzf integration**: use [fzf](https://github.com/junegunn/fzf) options to customize session selection.

## Requirements
- **tmux**: ensure [tmux](https://github.com/tmux/tmux) is installed on your system.

## Installation
```sh
go install github.com/gabefiori/ts@latest
```

## Configuration
Create a configuration file at `~/.config/ts/config.json`:

```json
{
   "targets":[
      {
         "path":"~/your/path",
         "depth":1
      },
      {
         "path":"/home/you/your_other/path",
         "depth":0
      }
   ],
   "selector":[
      "--height=60%"
   ]
}
```

- **targets**: List of directories to be navigated, with path specifying the directory and depth determining the level of subdirectories to display.
- **selector**: Options passed to fzf for customizing the selection interface. (Optional)

## Usage 
To start the sessionizer, run the following command:
```sh
ts
```

For more information about command-line options, use:
```sh
ts --help
```

## Adding a shortcut to tmux
To bind a key to create a new window and run the `ts` command, add the following line to your `.tmux.conf` file:

```bash
bind-key -r f run-shell "tmux neww ts"
```

The new window will close automatically after the command completes.
