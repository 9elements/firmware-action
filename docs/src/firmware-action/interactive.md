# Interactive debugging

```admonish note title="A little bit of backstory"
While I was playing around with `firmware-action` I found early on that debugging what is going on inside the docker container is rather lengthy and annoying process. This was the moment when the idea of some interactive option was born.
```

```admonish done collapsible=true title="Dropping the SSH feature in favor of Dagger build-in debugging"
Dagger [since v0.12](https://dagger.io/blog/dagger-0-12) supports new built-in interactive debugging.

~~We are already planning to re-write `firmware-action` to use this new feature instead of the `ssh` solution we are currently using. For more details see issue [269](https://github.com/9elements/firmware-action/issues/269).~~

**UPDATE:** It is possible now to use the new and shiny feature of dagger for interactive debugging! As a result we have dropped the SSH feature.

Related:
- Issue [#269](https://github.com/9elements/firmware-action/issues/269)
- Pull Request [#522](https://github.com/9elements/firmware-action/pull/522)

Supplementary dagger documentation:
- Blog post [Dagger 0.12: interactive debugging](https://dagger.io/blog/dagger-0-12)
- Documentation for [Interactive Debugging](https://docs.dagger.io/features/debugging)
- Documentation for [Custom applications](https://docs.dagger.io/api/sdk/#custom-applications)
- Documentation for [Interactive Terminal](https://docs.dagger.io/api/terminal/)
```

To leverage the use of interactive debugging, you have to install [dagger CLI](https://docs.dagger.io/install).

Then when using `firmware-action`, simply prepend the command with `dagger run --interactive`.

Instead of:
```bash
firmware-action build --config=firmware-action.json --target=coreboot-example
```
 Execute this:
```bash
dagger run --interactive firmware-action build --config=firmware-action.json --target=coreboot-example
```

If you are using `Taskfile` to abstract away some of the complexity that comes with larger projects, simply prepend the whole `Taskfile` command.

Instead of:
```bash
task build:coreboot-example
```

 Execute this:
```bash
dagger run --interactive task build:coreboot-example
```


On build failure you will be dropped into the container and can debug the problem.

To exit the container run command `exit` or press `CTRL+D`.
