# Local system

To get `firmware-action` look into [Get firmware-action](get_started/04_get_firmware_action.md) section.

## Run
```
firmware-action build --config=<path-to-JSON-config> --target=<my-target>
```

## Help

```
Usage: firmware-action --config=firmware-action.json <command> [flags]

Utility to create firmware images for several open source firmware solutions. Source code at 'https://github.com/9elements/firmware-action'

Flags:
  -h, --help                           Show context-sensitive help.
      --json                           switch to JSON stdout and stderr output
      --indent                         enable indentation for JSON output
      --debug                          increase verbosity
      --config=firmware-action.json    Path to configuration file, supports multiple flags to use multiple configuration files

Commands:
  build --config=firmware-action.json --target=STRING [flags]
    Build a target defined in configuration file. For interactive debugging preface the command with 'dagger run --interactive', for example
    'dagger run --interactive $(which firmware-action) build --config=...'. To install dagger follow instructions at https://dagger.io/

  generate-config --config=firmware-action.json [flags]
    Generate empty configuration file

  version --config=firmware-action.json [flags]
    Print version and exit

Run "firmware-action <command> --help" for more information on a command.
```
