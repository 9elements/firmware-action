# Tips and Tricks

## Python2 for building Intel FSP

Intel FSP can be build in the EDK2 containers. However the containers often have `Python3` as default.

Most `EDK2` containers have `Python2` installed and contain script `switch-to-python2` (`/bin/switch-to-python2`) which will let you easily switch to `Python2` as default.

```admonish tip
To see which python version are installed and which python version is used as default, look into our [compose.yaml](https://github.com/9elements/firmware-action/blob/main/docker/compose.yaml). Specifically, look into `edk2` containers and their arguments `PYTHON_PACKAGES` (which python versions are installed) and `PYTHON_VERSION` (which python version is the default).
```
