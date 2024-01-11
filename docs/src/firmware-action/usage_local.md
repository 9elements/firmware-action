# Local system

To get firmware-action, there are few options:
- [firmware-action](https://aur.archlinux.org/packages/firmware-action) Arch Linux AUR package
- build from source with [taskfile](https://taskfile.dev/)

To build from source:
```
git clone https://github.com/9elements/firmware-action.git
cd firmware-action
task build-go-binary
```

To run:
```
./firmware-action build --config=./firmware-action-config.json --target=my-target
```

