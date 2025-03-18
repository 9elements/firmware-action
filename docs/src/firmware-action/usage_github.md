# Github CI

You can use `firmware-action` as any other GitHub action. The action requires a configuration file (`config`) and target (`target`) to build, with additional options for artifact handling and caching to improve build times.

## Artifacts and Outputs

The action can automatically upload build artifacts when the `auto-artifact-upload` option is enabled (disabled by default). When enabled, the action provides the following outputs:

- `artifact-name`: Name of the uploaded artifact
- `artifact-id`: GitHub ID of the artifact (useful for REST API)
- `artifact-url`: Direct download URL for the artifact (requires GitHub login)
- `artifact-digest`: SHA-256 digest of the artifact

Internally, `firmware-action` uses the [actions/upload-artifact@v4](https://github.com/actions/upload-artifact)@v4 action to handle artifact uploads and [actions/download-artifact@v4](https://github.com/actions/download-artifact) to handle artifact downloads.

```admonish example
You can use these outputs in subsequent steps:
~~~yaml
{{#include ../firmware-action-example/.github/workflows/coreboot-example.yml:CorebootExample}}
~~~
```

You can configure artifact handling with these options:
- `artifact-if-no-files-found`: Behavior if no files are found ('warn', 'error', or 'ignore', default: 'warn')
- `artifact-compression-level`: Compression level for artifacts (0-9, default: '6')
- `artifact-overwrite`: Whether to overwrite existing artifacts with the same name (default: 'false')

## Build Caching

The action can cache build artifacts between runs to speed up builds, but this feature must be explicitly enabled with the `auto-cache` input (disabled by default).

When enabled, the cache is:
- Keyed by the config file contents and commit SHA
- Restored at the start of each run
- Saved after each run, even if the build fails

You can also use the `recursive` option to build a target with all its dependencies, and the `debug` option for increased verbosity when troubleshooting.

```admonish tip
You still might want to cache other files and directories as `firmware-action` caches only outputs and its temporary files.
```

## Artifact Downloading

When `auto-artifact-download` is enabled (disabled by default), the action automatically downloads all artifacts from the current workflow run. This feature is particularly useful in workflows with multiple jobs that depend on each other, as it can save a lot of copy-pasting between workflow steps.

```admonish note
Due to limitations in GitHub Actions, it's not possible for the action to download only relevant artifacts. When enabled, it downloads all artifacts from the current workflow run regardless of whether they're needed for the current build.
```

```admonish warning
When using `auto-artifact-download`, be careful to avoid naming conflicts in the `output_dir` values across different targets. Since all artifacts are downloaded and merged into the same workspace, files with the same paths will overwrite each other.
```

```admonish warning
At the time of writing, GitHub-hosted runners have only 14GB of disk space. If your artifacts are large and/or you have many matrix combinations, the workflow could run out of disk space. Monitor your disk usage when using this feature with large builds.
```

## Complete Configuration Example

For a complete example with all options:

```admonish example
~~~yaml
{{#include ../firmware-action-example/.github/workflows/linuxboot-example.yml:AllFeatures}}
~~~
```

## Parametric builds with environment variables

To take advantage of matrix builds in GitHub, it is possible to use environment variables inside the JSON configuration file.

```admonish example
For example let's make `RELEASE_TYPE` environment variable which will hold either `DEBUG` or `RELEASE`.

JSON would look like this:
~~~json
{
  "u-root": {
    "uroot-example-${RELEASE_TYPE}": {
      ...
      "output_dir": "output-linuxboot-uroot-${RELEASE_TYPE}/",
      ...
    }
  }
}
~~~

YAML would look like this:
~~~yaml
{{#include ../firmware-action-example/.github/workflows/linuxboot-example-separate-jobs.yml:Parametric}}
~~~
```
