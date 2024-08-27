## Unreleased (4b9e9d4..4b9e9d4)
#### Bug Fixes
- **(action/linux)** defconfig filename - (4b9e9d4) - AtomicFS

- - -

## v0.6.0 - 2024-08-27
#### Bug Fixes
- **(ci)** consolidate jobs - (a9e6b0d) - AtomicFS
- **(dagger)** missing docker-compose - (6b41c2e) - AtomicFS
- **(megalinter)** fix spelling - (61c2c1e) - AtomicFS
- **(megalinter)** fix spelling - (d614b25) - AtomicFS
#### Build system
- **(deps)** bump oxsecurity/megalinter from 7 to 8 - (e7e593a) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (602dff8) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (bf69e65) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto in /action in the golang group - (6752a89) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (f69338a) - dependabot[bot]
#### Features
- **(ci)** add reminder bot - (b3615d5) - AtomicFS
- **(docker)** build coreboot cross-compilers for all platforms - (3659a03) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.6.0 - (6a4f3d0) - AtomicFS
- **(action)** remove unnecessary input - (2b42c43) - AtomicFS

- - -

## v0.5.0 - 2024-08-27
#### Bug Fixes
- **(action)** remove unnecessary apostrophes - (f715557) - AtomicFS
- **(action)** if statement using compile input - (77403bf) - AtomicFS
- **(ci)** consolidate multiple jobs into few status checks - (6bdd189) - AtomicFS
- **(ci)** jobs are canceled on submitted review - (28516b9) - AtomicFS
- **(lint)** fix broken matching in .editorconfig - (6bc13d3) - AtomicFS
- **(lint)** fixes according to megalinter - (fd42f59) - AtomicFS
- **(megalinter)** fix spelling - (61587d4) - AtomicFS
- **(megalinter)** exclude mdbook theme from linting - (744824c) - AtomicFS
- **(release)** forgot to remove the old version bumper - (a3ab36c) - AtomicFS
#### Documentation
- **(action)** cleanup anchors in recipes - (fdbaa04) - AtomicFS
- add CODEOWNERS file - (0f2a143) - AtomicFS
- remove google analytics from mdbook - (eab70ac) - AtomicFS
- add note about interactive mode - (220baf9) - AtomicFS
- add note about Arch Linux AUR package - (9bc4ff3) - AtomicFS
- add notes about interactive mode - (a857531) - AtomicFS
- add more docs - (5737270) - AtomicFS
- add firmware-action-example as submodule - (3f9f888) - AtomicFS
- split example_config files for more clarity - (499bf53) - AtomicFS
- update README.md - (137915e) - AtomicFS
#### Features
- **(ci)** add labeler bot - (d88ba9f) - AtomicFS
- **(release)** automatically determinate the next version - (2306d2d) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.5.0 - (1dc275f) - AtomicFS
- rename labeler workflow - (be50b25) - AtomicFS

- - -

## v0.4.0 - 2024-08-27
#### Bug Fixes
- **(action)** broken InputDirs - (4ae3cba) - AtomicFS
- **(action)** simplify u-root test - (ed3f45c) - AtomicFS
- **(action)** container.Export returns string instead of boolean - (e76678c) - AtomicFS
- **(ci)** add conditional JIT compilation from examples - (cdf1986) - AtomicFS
- **(ci)** remove go setup from example jobs - (ecb09aa) - AtomicFS
- **(ci)** skip on all reviews except when the release PR - (4faed8d) - AtomicFS
- **(ci/example)** dorny/paths-filter issue - (45ec8c0) - AtomicFS
- **(docker)** add missing dependencies for building edk2 - (6eba755) - AtomicFS
- **(lint)** fix frequent markdown-link-check fail with AUR link - (9d230a6) - AtomicFS
- **(megalinter)** fix spelling - (76e263d) - AtomicFS
- **(megalinter)** fix spelling - (cce075b) - AtomicFS
- **(release)** remove leftover javascript stuff - (73097e7) - AtomicFS
- **(release)** typo - (6b9959f) - AtomicFS
- **(release)** release-prepare missing tags and history - (32c8c6d) - AtomicFS
#### Build system
- **(deps)** bump github.com/heimdalr/dag in /action in the golang group - (fd61b56) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (a0c30f2) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (40d5baa) - dependabot[bot]
- **(deps)** update pytest-flake8 requirement in /.dagger-ci/daggerci - (e2d1237) - dependabot[bot]
- **(deps)** update pytest requirement in /.dagger-ci/daggerci - (86cfa9d) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (96ec784) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (f4cb409) - dependabot[bot]
- **(deps)** update dagger-io requirement - (1162f25) - dependabot[bot]
#### Continuous Integration
- **(action)** download executable instead of JIT compilation - (8f28305) - AtomicFS
- **(docker)** extra step to validate compose.yaml - (e493118) - AtomicFS
- **(docker)** bump python to v3.12 - (757da02) - AtomicFS
- **(docker)** add first u-root container v0.14.0 - (1447da2) - AtomicFS
- **(example)** fix example workflows to use our uroot container - (5c116f6) - AtomicFS
- more tweaks to run conditions to reduce number of CIs - (6927afc) - AtomicFS
- tweak run conditions to reduce number of CIs - (bce7e82) - AtomicFS
#### Documentation
- remove obsolete information - (bbc402b) - AtomicFS
- add notes on linux defconfig - (0cef84f) - AtomicFS
#### Features
- **(action)** add InputFiles and InputDirs options - (778f0bd) - AtomicFS
- **(action)** check if file from needed modules exist - (a13d992) - AtomicFS
- **(action)** print summary overview at the end - (46fd6ff) - AtomicFS
- **(action)** do not override output files - (1040cd5) - AtomicFS
- **(ci)** use cog to bump version - (ab9dcda) - AtomicFS
- **(docker)** add linux 6.9.9 - (a37d9ea) - AtomicFS
- **(test)** add examples for non-Linux systems - (7292ac9) - AtomicFS
- add CHANGELOG.md - (bc0cb58) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.4.0 - (676c7b6) - AtomicFS
- **(action)** bump version to v0.4.0 - (a7dcfe3) - AtomicFS
- **(daggerci)** megalinter fix - (955747d) - AtomicFS
- should not happen duplicate - (72ccaf2) - AtomicFS
- fix cosmetic issues - (a2dc60c) - AtomicFS

- - -

## v0.3.2 - 2024-08-27
#### Features
- **(action)** allow multi-module workspaces for u-root - (f54803d) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.3.2 - (4cac34c) - AtomicFS

- - -

## v0.3.1 - 2024-08-27
#### Bug Fixes
- **(again)** build docker containers on release - (2d33a7e) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.3.1 - (ca343ff) - AtomicFS

- - -

## v0.3.0 - 2024-08-27
#### Bug Fixes
- **(megalinter)** fix spelling - (99a6247) - AtomicFS
- **(typo)** typo - (a5d9fb5) - AtomicFS
- build docker containers on release - (9872682) - AtomicFS
#### Build system
- **(deps)** bump actions/github-script from 6 to 7 - (92851a9) - dependabot[bot]
- **(deps)** bump google.golang.org/grpc in /action - (8adf0f8) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto in /action in the golang group - (6018d7e) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (f2d1f92) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (896d3dc) - dependabot[bot]
- **(deps-dev)** bump mega-linter-runner in the javascript group - (cd0ca04) - dependabot[bot]
#### Continuous Integration
- **(docker)** add coreboot 24.02.01 and 24.05 - (fc4388d) - AtomicFS
- build docker images on tags - (c928e2c) - AtomicFS
#### Features
- **(action)** add support for u-root - (6d4dba7) - AtomicFS
- **(action)** add option to ignore missing blob files - (0697a02) - AtomicFS
- **(action)** notify user that interactivity rolls back state - (bf26e96) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.3.0 - (35ed1fe) - AtomicFS
- **(cosmetic)** run gofumt formatter - (01af31b) - AtomicFS
- add output-* directories into .gitignore - (ea97f8b) - AtomicFS

- - -

## v0.2.1 - 2024-08-27
#### Bug Fixes
- **(action)** fix issue 195 - (d8bc51a) - AtomicFS
- **(action)** fix bin filename in Taskfile.yml - (5b20be4) - AtomicFS
- **(action)** sync versions across files - (2d269f3) - AtomicFS
- **(ci/example)** fix coreboot blobs cache - (b5cd9b2) - AtomicFS
- **(ci/example)** fix coreboot blobs URL and add cache - (0e66cf2) - AtomicFS
- **(megalinter)** disable trivy checks for GitHub specific configs - (271e78b) - AtomicFS
- **(megalinter)** disable v8r yaml linter for while - (15b9aed) - AtomicFS
- release workflow does not triger other workflows - (b8b8376) - AtomicFS
- return not checked in container/ssh.go - (d747064) - AtomicFS
#### Build system
- **(deps)** bump github.com/go-playground/validator/v10 - (f8081de) - dependabot[bot]
- **(deps)** bump github.com/vektah/gqlparser/v2 in /action - (02d56ba) - dependabot[bot]
- **(deps)** run update on npm packages - (fa486d6) - AtomicFS
- **(deps)** bump the golang group across 1 directory with 2 updates - (f4f891a) - dependabot[bot]
- **(deps)** bump goreleaser/goreleaser-action from 5 to 6 - (ca0ebc1) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (fb341e2) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (83e387a) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (f55fd9f) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (044ed91) - dependabot[bot]
- **(deps)** bump the golang group across 1 directory with 2 updates - (2603b82) - dependabot[bot]
- **(deps)** bump golangci/golangci-lint-action from 5 to 6 - (5b2b5bf) - dependabot[bot]
- **(deps)** bump github.com/go-playground/validator/v10 - (404745e) - dependabot[bot]
- **(deps)** update pytest requirement in /.dagger-ci/daggerci - (43c3558) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (fac8a4d) - dependabot[bot]
- **(deps)** bump golangci/golangci-lint-action from 4 to 5 - (35f7407) - dependabot[bot]
- **(deps)** bump golang.org/x/net from 0.21.0 to 0.23.0 in /action - (e1d19f2) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (fc9814f) - dependabot[bot]
- **(deps)** bump the golang group in /action with 2 updates - (72c8636) - dependabot[bot]
- **(deps)** update dagger-io requirement - (8ef2d88) - dependabot[bot]
- **(deps)** bump actions/configure-pages from 4 to 5 - (84b726e) - dependabot[bot]
- **(deps)** update pytest requirement in /.dagger-ci/daggerci - (02fc8f5) - dependabot[bot]
- **(deps)** bump wagoid/commitlint-github-action from 5 to 6 - (2dba929) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (e963114) - dependabot[bot]
- **(deps)** update pytest-cov requirement in /.dagger-ci/daggerci - (af0594d) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (1844bdc) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (072b13a) - dependabot[bot]
- **(deps-dev)** bump mega-linter-runner in the javascript group - (7ceadf2) - dependabot[bot]
- **(deps-dev)** bump ejs from 3.1.9 to 3.1.10 - (24fdff5) - dependabot[bot]
- **(deps-dev)** bump mega-linter-runner in the javascript group - (db36859) - dependabot[bot]
#### Continuous Integration
- **(action)** add goreleaser configuration file - (47460bc) - AtomicFS
- **(action)** add release workflow - (db7ac14) - AtomicFS
- **(action)** add release-prepare workflow - (9ca7f4f) - AtomicFS
- **(dependabot)** increase open-pull-requests-limit - (2a94539) - AtomicFS
- **(golangci)** update config - (a539fe1) - AtomicFS
- fix duplicated workflows on release - (df3f973) - AtomicFS
#### Features
- **(action)** print a message on success - (0de6d54) - AtomicFS
- **(action)** print warning when uninitialized git submodule found - (b18e86e) - AtomicFS
- **(action)** add version sub-command - (bfc46cd) - AtomicFS
- **(action/container)** catch timeout and suggest restart (WIP) - (0845aa8) - AtomicFS
- **(daggerci)** print size of tarball into stdout - (1140359) - AtomicFS
- **(logging)** add custom slog logger - (b9eb652) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.2.1 - (65aa6bd) - AtomicFS
- **(action)** fix according to megaliter complains - (8dd75a2) - AtomicFS
- **(action)** add go-git and run go mod tidy - (261132f) - AtomicFS
- **(action)** fix according to megaliter complains - (88f5def) - AtomicFS
- **(action)** remove obsolete comment - (5aa4a14) - AtomicFS
- **(daggerci)** megalinter fix - (8ac3179) - AtomicFS
- **(daggerci)** run task format - (4850a3a) - AtomicFS
- **(megalinter)** spelling - (003d2a4) - AtomicFS
#### Performance Improvements
- **(docker)** reduce rust size - (92f2e12) - AtomicFS
#### Refactoring
- **(action)** pass config struct around as pointer - (c9fb602) - AtomicFS
- **(action/config)** improve logging with slog, better JSON - (44961d1) - AtomicFS
- **(action/container)** improve logging with slog - (523658b) - AtomicFS
- **(action/coreboot)** improve logging with slog - (6c83930) - AtomicFS
- **(action/edk2)** improve logging with slog - (147302d) - AtomicFS
- **(action/linux)** improve logging with slog - (6024156) - AtomicFS
- **(action/main)** improve logging with slog and fix error handling - (754235b) - AtomicFS
- **(action/recipes)** improve logging with slog - (93b53e7) - AtomicFS
- **(action/stitching)** improve logging with slog - (6c3dba2) - AtomicFS
#### Tests
- **(action)** add slog linting to golangci-lint - (37f0b82) - AtomicFS
- **(megalinter/trivy)** add exception to Dockerfiles - (34893e9) - AtomicFS

- - -

## v0.2.0 - 2024-08-27
#### Build system
- **(deps)** bump the golang group in /action with 1 update - (d7e9dc8) - dependabot[bot]
- **(deps)** update pytest-timeout requirement in /.dagger-ci/daggerci - (11045e2) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (a8270bc) - dependabot[bot]
- **(deps-dev)** bump the javascript group with 1 update - (960b087) - dependabot[bot]
#### Features
- **(action/interactive)** update CLI - (52ebc86) - AtomicFS
- **(action/interactive)** update remaining tests - (1c806f9) - AtomicFS
- **(action/interactive)** open SSH tunnel if build fails - (c1565d3) - AtomicFS
- **(action/interactive)** update recipes - (b17a728) - AtomicFS
- **(action/interactive)** add function to start SSH in container - (037eac7) - AtomicFS

- - -

## v0.1.2 - 2024-08-27
#### Bug Fixes
- **(commitlint)** add config increasing max line length - (38a005e) - AtomicFS
- **(docker)** edk2 repositories were missing files - (4cf4603) - AtomicFS
- **(docker)** missing ifdtool and cbfstool in container - (87ffcf5) - AtomicFS
- **(docs)** fix mdbook and mdbook-graphviz incompatibility - (f3b6739) - AtomicFS
#### Build system
- **(deps)** bump the golang group in /action with 2 updates - (c0c02d2) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (be21ef7) - dependabot[bot]
- **(deps)** bump the python group in /.dagger-ci/daggerci with 1 update - (662bc12) - dependabot[bot]
- **(deps)** bump actions/cache from 3 to 4 - (ec47480) - dependabot[bot]
- **(deps)** bump the golang group in /action with 1 update - (7d52088) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (c234b3a) - dependabot[bot]
- **(deps)** bump actions/upload-artifact from 3 to 4 - (f740513) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (2d6c0d9) - dependabot[bot]
- **(deps)** bump actions/checkout from 3 to 4 - (843e9c1) - dependabot[bot]
- **(deps)** bump actions/setup-go from 4 to 5 - (44a2ac8) - dependabot[bot]
- **(deps)** bump golangci/golangci-lint-action from 3 to 4 - (92faa4e) - dependabot[bot]
- **(deps)** bump actions/setup-python from 4 to 5 - (60601b7) - dependabot[bot]
- **(deps)** bump actions/upload-artifact from 3 to 4 - (61e3df8) - dependabot[bot]
- **(deps)** update pytest requirement in /.dagger-ci/daggerci - (a3f2a78) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (1aad502) - dependabot[bot]
- **(deps)** bump the python group in /.dagger-ci/daggerci with 1 update - (5ed9c54) - dependabot[bot]
- **(deps)** bump the golang group in /action with 3 updates - (404cf82) - dependabot[bot]
- **(deps-dev)** bump mega-linter-runner from 6.22.2 to 7.9.0 - (fe50ea2) - dependabot[bot]
#### Continuous Integration
- **(action)** further reduce load on CI - (1c7ac85) - AtomicFS
- **(dependabot)** tweak config - (48a4b1c) - AtomicFS
- **(docker)** add rust support into edk2 - (10e3eaa) - AtomicFS
- **(docker)** add coreboot 24.02 example - (621a41d) - AtomicFS
- **(docker)** add coreboot 24.02 - (fbc453e) - AtomicFS
- **(docs)** add cache for mdbook packages - (ed7b2ce) - AtomicFS
- **(test)** add example test - (0d58856) - AtomicFS
- **(workflows)** change go version in linting to stable - (8499f36) - AtomicFS
- delete old packages - (d49be74) - AtomicFS
- add cache to reduce stress on git hosting servers - (31daa3c) - AtomicFS
- cosmetic - (c50c4a4) - AtomicFS
- limit Docker building only on change - (7671b9b) - AtomicFS
- add dependabot - (fd729a9) - AtomicFS
#### Features
- **(action)** add support for firmware stitching - (d8d9eb6) - AtomicFS
- **(docker)** prepare container for interactive debugging - (10d7abc) - AtomicFS
#### Miscellaneous Chores
- **(megalinter)** unused variables - (93c760e) - AtomicFS
- **(megalinter)** update npm packages - (76d0adb) - AtomicFS
#### Refactoring
- **(action)** move Arch out of commonOpts - (7567096) - AtomicFS

- - -

## v0.1.1 - 2024-08-27
#### Bug Fixes
- **(action)** naming mistake in JSON config - (58ebf67) - AtomicFS

- - -

## v0.1.0 - 2024-08-27
#### Bug Fixes
- **(docker)** add cleanup commands - (b26acf4) - AtomicFS
- guard queue with mutex in recipes.go - (b160e01) - Marvin Drees
#### Continuous Integration
- run docker build workflows on release - (0a1aadb) - AtomicFS
- add coreboot 4.22.01 - (758be1a) - AtomicFS
#### Miscellaneous Chores
- **(action)** replace word Docker with Container - (fc16825) - AtomicFS
- remove build-coreboot task from taskfile - (fc35d09) - AtomicFS
- update .gitignore - (1b36540) - AtomicFS
#### Refactoring
- **(action)** update example JSON - (6386be8) - AtomicFS
- **(action)** add buildFirmware methods - (7457ad2) - AtomicFS
- **(action)** simplify building of dependency forest - (8fca0aa) - AtomicFS
- **(action)** cleanup interfaces for coreboot, linux and ekd2 - (4092f70) - AtomicFS
- **(action)** refactor edk2 - (b013bb1) - AtomicFS
#### Tests
- **(action)** test for race conditions - (3db33ea) - AtomicFS

- - -

## v0.0.1 - 2024-08-27
#### Bug Fixes
- **(Dockerfile)** add --no-cache-dir - (c116873) - Patrick Rudolph
- **(Dockerfile)** add user to fix linter - (ebca447) - Patrick Rudolph
- **(action)** better architecture handling in edk2 - (934a72f) - AtomicFS
- **(action)** switch from WithMountedDirectory to WithDirectory - (ef048b0) - AtomicFS
- **(action)** create destination directory for mounting - (7100312) - AtomicFS
- **(action)** use relative paths for blobs in container - (944508d) - AtomicFS
- **(action)** fix path to Go project in JavaScript shim - (81c5977) - AtomicFS
- **(action)** bump version of NodeJS to v20 - (3c9de85) - AtomicFS
- **(action)** defconfig path handling - (435891f) - AtomicFS
- **(action)** differentiate between GitHub input and ENV var - (f04b24e) - AtomicFS
- **(action)** fix javascript - (45ea50b) - AtomicFS
- **(action)** cosmetic fixes to edk2 - (7e39c9d) - AtomicFS
- **(action)** fix URL detection for container setup - (d0b8a97) - AtomicFS
- **(action)** exporting directory from container - (6d504ac) - AtomicFS
- **(action)** review corrections - (8966650) - AtomicFS
- **(bash)** fix according to shfmt in megalinter - (c163f47) - AtomicFS
- **(ci)** examples run only on PR or main branch - (0f15874) - AtomicFS
- **(dagger)** update to dagger 0.8.x - (11e4ae8) - AtomicFS
- **(dagger)** publishing - (388ebb4) - AtomicFS
- **(dagger)** fixes according to pylint and mypy - (9c10db1) - AtomicFS
- **(dagger)** fix invalid Python package name - (3a7691d) - AtomicFS
- **(dagger/orchestrator)** arguments for concurrent run - (111fd1c) - AtomicFS
- **(docker)** add missing pkgconf package - (2745ff7) - AtomicFS
- **(docker)** remove WORKDIR for edk2 dockerfile - (f16953f) - AtomicFS
- **(docker)** add nodejs - (4a34186) - AtomicFS
- **(docker)** add nodejs - (ce4953a) - Patrick Rudolph
- **(docker)** vUDK2017 - (e923750) - AtomicFS
- **(dockerfile)** add missing dependency - (ab744f6) - AtomicFS
- **(dockerfiles)** fix according to hadolint in megalinter - (6c43a2a) - AtomicFS
- **(go)** reorder files and update deps - (cc3915f) - Marvin Drees
- **(megalinter)** disable cspell for go.mod and go.sum - (ad3b369) - AtomicFS
- **(megalinter)** temporarily disable go linters - (30256f6) - AtomicFS
- **(megalinter)** disable hadolint DL3008 (pin versions in apt get install) - (d2016a6) - AtomicFS
- **(test)** fix potential problem with edk2 test - (7a9755d) - AtomicFS
- python-black code formatter - (48c7e76) - AtomicFS
- documentation and error handling on docker-compose subprocess call - (29be42d) - AtomicFS
- infinite loop issue #37 - (f353660) - AtomicFS
- typos - (ab85440) - AtomicFS
- remove ubuntu user from Dockerfiles - (a10095a) - AtomicFS
- fix variables in test job - (45c18dd) - AtomicFS
- fix dependencies between jobs - (51057f8) - AtomicFS
- switch to coreboot mirror for buildgcc packages - (9bd8f12) - AtomicFS
- multiple top-level headings - (38f0f10) - Vojtech Vesely
- add .checkov.yml - (d998599) - Patrick Rudolph
- satify js standard linter - (b2f365c) - Marvin Drees
- add missing cspell words - (8615578) - Marvin Drees
- satisfy prettier linter - (b33b3b9) - Marvin Drees
- satify golangci lint - (9458e35) - Marvin Drees
#### Build system
- **(go)** add Taskfile.yml to go-task / task - (060cc09) - AtomicFS
#### Continuous Integration
- **(example)** update workflow - (28c894e) - AtomicFS
- **(linter)** use megalinter cupcake - (dd55906) - Marvin Drees
- **(mega-linter)** configure megalinter and update to cupcake v7 - (942eee6) - AtomicFS
- **(megalinter)** update dictionary - (4f9288f) - AtomicFS
- **(megalinter)** tweak go linting - (38dc89e) - AtomicFS
- **(megalinter)** enable go linting - (7c4bedc) - AtomicFS
- add mdbook workflow - (9776a57) - AtomicFS
- add go setup step - (1552e45) - AtomicFS
- increase timeout for go test - (1fcaf00) - AtomicFS
- add coreboot 4.21 (fix missing 4.20.1) - (9097333) - AtomicFS
- cleanup workflows, fix pytest and add new go-test - (a6802bd) - AtomicFS
- run ci also when PR is merged to main - (82baf4d) - AtomicFS
- more tweak triggers for ci - (92b6004) - AtomicFS
- tweak triggers for ci - (9359daf) - AtomicFS
- make publishing conditional - (d270fb2) - AtomicFS
- rename pytest job to default 'pytest' - (9e07fb4) - AtomicFS
- re-introduce matrix build - (1175025) - AtomicFS
- cleanup tests - (50c778f) - AtomicFS
- reusable matrix definition - (7d190ad) - AtomicFS
- make job continue even on error - (fdd9e74) - AtomicFS
- use matrix strategy - (a2058d9) - AtomicFS
- add workflow to build docker images - (db15702) - AtomicFS
- run linter on every push and pull request - (d10b270) - AtomicFS
- add tests for seabios with coreinfo and nvramcui - (93f275a) - AtomicFS
- use parametric payload approach on coreboot testing - (74372c1) - AtomicFS
- add linters - (2766350) - Marvin Drees
#### Documentation
- **(README.md)** add badge to go tests - (02d8423) - AtomicFS
- **(README.md)** add badges because why not - (8947bcd) - AtomicFS
- **(action)** add anchors into code for documentation - (7c3bd8b) - AtomicFS
- **(action)** add artifacts to example workflow - (e3e9f0a) - AtomicFS
- **(action)** add example workflows - (3f89df6) - AtomicFS
- add link to documentation into main.go - (22f075f) - AtomicFS
- add footnotes - (b393b6f) - AtomicFS
- add mdbook build into taskfile - (a3e063b) - AtomicFS
- add some basic documentation - (89ce9e2) - AtomicFS
- add links to sources - (d09b6ca) - AtomicFS
- add brief README on defconfigs - (c68231b) - AtomicFS
- add CONTRIBUTING.md - (a781122) - Vojtech Vesely
#### Features
- **(Dockerfile)** add UEFI dockerfiles - (f976dbf) - Patrick Rudolph
- **(Dockerfile)** set default toolchain path - (ced2628) - Patrick Rudolph
- **(Dockerfile)** add python3 for FSP builds - (df6a1dd) - Patrick Rudolph
- **(action)** recursive builds - (09c3ffe) - AtomicFS
- **(action)** add new env var to indicate which GCC version to use - (74cf26f) - AtomicFS
- **(action)** add GCC version into edk2 and Linux builds - (f0d62de) - AtomicFS
- **(action)** add blob support into coreboot build - (a3a0dfc) - AtomicFS
- **(action)** add edk2 build function - (a64e305) - AtomicFS
- **(action)** add DirTree function into action/filesystem - (69483c0) - AtomicFS
- **(action)** add option to build container from Dockerfile - (dac6d76) - AtomicFS
- **(action)** add linux build function - (464c049) - AtomicFS
- **(action)** add dagger and container related functions - (63dc7ff) - AtomicFS
- **(action)** add file-system related functions - (eb1cbb7) - AtomicFS
- **(action)** remove custom kconfig code - (d451431) - AtomicFS
- **(action)** implement coreboot target - (0e0cc10) - Patrick Rudolph
- **(action.yml)** define first supported options - (27b1aa7) - Patrick Rudolph
- **(coreboot)** add initial dagger code - (68cba98) - Marvin Drees
- **(dagger)** change requirements.txt version specifier to compatible - (cca0786) - AtomicFS
- **(dagger)** add orchestrator to abstract away all complexity (WIP) - (604b0ad) - AtomicFS
- **(dagger)** add requirements.txt - (a48545d) - AtomicFS
- **(dagger)** add function to read env var with fallback value - (61c0e86) - AtomicFS
- **(dagger)** add git commit sha and describe functions - (dd1f86d) - AtomicFS
- **(dagger)** add functions for files and filesystem - (50cfab1) - AtomicFS
- **(dagger)** add validation to docker-compose file parser - (6ac492a) - AtomicFS
- **(dagger)** add docker-compose file parser - (d485ede) - AtomicFS
- **(dagger)** prepare for rewite into dagger with python - (b2bae63) - AtomicFS
- **(dagger/orchestrator)** improve error handling - (eed7b8e) - AtomicFS
- **(docker)** add cross-compile toolchains - (cbed931) - AtomicFS
- **(docker)** add support for Linux Kernel - (76c2604) - AtomicFS
- **(docker)** Add dockerfile for coreboot:4.19 - (2dbdbb5) - Patrick Rudolph
- **(taskfile)** shuffle tests - (6d9e914) - AtomicFS
- add VBT and EC blob support - (fb89593) - AtomicFS
- add new release of coreboot 4.20.1 - (f3372e3) - AtomicFS
- extend build to remaining images - (06883b6) - AtomicFS
- add edk2-stable202211 Dockerfile - (1ad870a) - Patrick Rudolph
- add edk2-stable202208 Dockerfile - (b535645) - Patrick Rudolph
- initial commit - (d0ed9d2) - Marvin Drees
#### Miscellaneous Chores
- **(action)** remove obsolete todo - (e1ffbd3) - AtomicFS
- **(action)** go mod tidy - (6425ef8) - AtomicFS
- **(action)** go mod tidy - (14eb34d) - AtomicFS
- **(action)** go mod tidy - (6410108) - AtomicFS
- **(action)** cleanup input names - (0e5636b) - AtomicFS
- **(action)** go mod tidy - (b233d50) - AtomicFS
- **(coreboot)** change ubuntu release from 22.04 to jammy - (5a9bc71) - AtomicFS
- **(go)** cleanup action go module (WIP) - (dae43fb) - AtomicFS
- **(megalinter)** cosmetic fixes - (52617bd) - AtomicFS
- **(megalinter)** cosmetic fixes - (1a12f64) - AtomicFS
- update dagger and other dependencies - (e033863) - AtomicFS
- fix typos and cosmetic problems - (8e3b447) - AtomicFS
- typos and cosmetic changes - (083b846) - AtomicFS
- move coreboot defconfig into separate directory - (727aaf7) - AtomicFS
- update dagger - (1ad3dfe) - Marvin Drees
#### Refactoring
- **(action)** switch to structs with embedded types - (4b1dada) - AtomicFS
- **(action)** switch from github action to JSON config - (87cbec9) - AtomicFS
- **(action)** cleanup linux and coreboot builds - (7994ae1) - AtomicFS
- **(action)** minor cleanup in filesystem - (b6d9f06) - AtomicFS
- **(action)** cleanup container tests - (8c8280d) - AtomicFS
- **(action)** filesystem functions - (90aa121) - AtomicFS
- **(action)** complete refactor of coreboot and recepies (WIP) - (85a3aef) - AtomicFS
- **(action.yml)** inputs - (f04dd56) - AtomicFS
- **(action/container)** lot of improvements - (9cf5625) - AtomicFS
- **(dagger)** fixes according to isort and ruff - (9750c13) - AtomicFS
- **(dagger)** fixes according to python-black code formatter - (e7fcb13) - AtomicFS
- **(dagger)** refactor handling results into separate class - (70a6308) - AtomicFS
- **(dagger)** finish rewrite into dagger with python (WIP) - (c011f1a) - AtomicFS
- **(docker)** refactor all edk2 dockerfiles - (7d3b5b4) - AtomicFS
- rework CLI and GitHub integration - (d70013b) - AtomicFS
- docker-compose file parser - (e94de83) - AtomicFS
- parametric Dockerfile for edk2 - (f574da7) - AtomicFS
- parametric Dockerfile for coreboot - (6163598) - AtomicFS
- move testing into script - (21b3724) - AtomicFS
- merge building and testing into single workflow - (8bd5ba2) - AtomicFS
- migrate kconfig to pkg - (bf192a7) - Patrick Rudolph
- rename go package to main - (f36fbcb) - Marvin Drees
#### Style
- **(tests/docker_compose)** more descriptive variable names - (a38d356) - AtomicFS
- fix according to review - (f3da6e9) - AtomicFS
- fixes according to megalinter - (f3236e0) - AtomicFS
- unify use of 'dockerfile' and 'docker-compose' - (8a2fb30) - AtomicFS
- rename workflow - (122ce88) - AtomicFS
- cosmetic - (aeaa824) - AtomicFS
- cosmetic - (7373f2d) - AtomicFS
- fix megalinter complaining - (04c8c34) - AtomicFS
- fix megalinter complaining - (902b405) - AtomicFS
- fix indentation to editorconfig specs - (7fdd2fe) - Vojtech Vesely
- add editorconfig - (3ca5baf) - Vojtech Vesely
- remove trailing whitespaces - (d016bbf) - Marvin Drees
- replace rome with megalinter/prettier - (d76c095) - Marvin Drees
#### Tests
- **(action)** run tests in parallel - (23945e9) - AtomicFS
- **(dagger)** use taskfile, fix pytest - (92af6b1) - AtomicFS
- **(docker-compose)** 100% coverage - (298e3c2) - AtomicFS
- **(go)** go mod tidy - (a95f28c) - AtomicFS
- **(go)** add coverage package - (f5d22cb) - AtomicFS
- **(pytest)** allow passing additional arguments to run_pytest.sh - (86f1dab) - AtomicFS
- **(pytest)** fix test with broken dockerfile - (04f3c97) - AtomicFS
- add pytest job - (736a21a) - AtomicFS
- add docker testing workflow - (43741c3) - Vojtech Vesely


