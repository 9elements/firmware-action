# Local system

To get firmware-action, there are few options:
- download compiled binary executable from [releases](https://github.com/9elements/firmware-action/releases)
- build from source with [taskfile](https://taskfile.dev/)
- `firmware-action` [Arch Linux AUR package](https://aur.archlinux.org/packages/firmware-action)


## Build from source
```
git clone https://github.com/9elements/firmware-action.git
cd firmware-action
task build-go-binary
```


## Run
```
./firmware-action build --config=<path-to-JSON-config> --target=<my-target>
```


## Help

```
Usage: firmware-action --config="firmware-action.json" <command> [flags]

Utility to create firmware images for several open source firmware solutions

Flags:
  -h, --help                             Show context-sensitive help.
      --json                             switch to JSON stdout and stderr output
      --indent                           enable indentation for JSON output
      --debug                            increase verbosity
      --config="firmware-action.json"    Path to configuration file

Commands:
  build --config="firmware-action.json" --target=STRING [flags]
    Build a target defined in configuration file

  generate-config --config="firmware-action.json" [flags]
    Generate empty configuration file

  version --config="firmware-action.json" [flags]
    Print version and exit

Run "firmware-action <command> --help" for more information on a command.
```

