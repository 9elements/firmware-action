# Run `firmware-action` in GitHub CI

Now that we have `firmware-action` working on local system. Let's set up CI.

```admonish example title=".github/workflows/example.yml"
~~~yaml
{{#include ../../firmware-action-example/.github/workflows/coreboot-example.yml}}
~~~
```

Commit, push and watch. And that is it.

The action automatically handles artifact uploading and caching. See [GitHub CI Usage](../usage_github.md) for details about artifacts, outputs and caching.

Now you should be able to build coreboot in CI and on your local machine.

