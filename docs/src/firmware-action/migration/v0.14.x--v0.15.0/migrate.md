# Migration guide from v0.14.x to v0.15.0

Drop-in replacement, should work out of the box. But there are many simplifications and quality of life improvements you might want to check out.

Most important being the optional automatic handling of caches and artifacts. This can greatly simplify your workflows.

```patch
{{#include ./workflow.yml.patch}}
```
