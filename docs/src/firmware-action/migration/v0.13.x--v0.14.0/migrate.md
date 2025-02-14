# Migration guide from v0.13.x to v0.14.0

The handling of `coreboot` blobs have been refactored.

`coreboot` can have far more blobs that we supported and the setup we had would not scale. We remove all of the hard-coded stuff and replace it with much more flexible setup where use has to define key-value map for the blobs.

The old way of defining blobs:
~~~json
blobs: {
    "payload_file_path": "./my-payload.bin"
}
~~~

The new way:
~~~json
blobs: {
  "CONFIG_PAYLOAD_FILE": "./my-payload.bin"
}
~~~

We have made a script that will allow you to migrate:
~~~bash
{{#include migrate-config.sh}}
~~~
