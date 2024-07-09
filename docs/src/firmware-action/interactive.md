# Interactive mode


```admonish note title="A little bit of backstory"
While I was playing around with firmware-action I found early on that debugging what is going on inside the docker container is rather lengthy and annoying process. This was the moment when the idea of some interactive option was born.
```

```admonish bug title="Issue [#109](https://github.com/9elements/firmware-action/issues/109)"
```

```admonish done title="Pull request [#147](https://github.com/9elements/firmware-action/pull/147)"
```

On build failure open `ssh` server in the container and let user connect into it to debug the problem. To enable this feature user has to pass argument `--interactive`. User can ssh into the container with a randomly generated password.

The container will keep running until user presses `ENTER` key.


```admonish attention
The container is launched in the interactive mode before the failed command was started.

This reverting is out of technical necessity.
```

```admonish note
The containers in dagger (at the time of writing) are single-use non-interactive containers. Dagger has a pipeline (command queue for each container) which starts executing only when specific functions such as [Sync()](https://pkg.go.dev/dagger.io/dagger#Container.Sync) are called which trigger evaluation of the pipeline inside the container.

To start a `ssh` service and wait for user to log-in, the container has to be converted into a [service](https://pkg.go.dev/dagger.io/dagger#Container.AsService) which also forces evaluation of the pipeline. And if any of the commands should fail, it would fail to start the `service` container.

As a workaround, when the evaluation of pipeline fails in the container, the container from previous step is converted into a `service` container with everything as it was just before the failing command was executed. In essence, when you connect, you end up in pristine environment.

~~~go
{{#include ../../../action/container/ssh.go:ContainerAsService}}
~~~
```

