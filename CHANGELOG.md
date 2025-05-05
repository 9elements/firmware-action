# Changelog

All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

- - -
## v0.17.2 - 2025-05-05
#### Bug Fixes
- **(action)** handle nested output directories in artifact caching - (7e5b2e4) - AtomicFS
- **(cmd)** store artifact path in txt file - (a75dd06) - AtomicFS
#### Build system
- **(deps)** bump golangci/golangci-lint-action from 7 to 8 - (03578de) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (00af04c) - dependabot[bot]
- **(deps)** bump github.com/sethvargo/go-githubactions - (5e32903) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (f5e3af8) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (209376f) - dependabot[bot]
#### Continuous Integration
- **(lint)** enable GitHub reporter for megalinter - (ffdc34f) - AtomicFS
- limit maximum parallel jobs - (aa11e5d) - AtomicFS
#### Documentation
- update notes on adding new container - (7ea554f) - AtomicFS
- fix forgotten old naming in migration guide - (10991a2) - AtomicFS
#### Miscellaneous Chores
- **(docker)** add new containers - (aa41509) - AtomicFS
- **(linter)** formatting changes according to linter - (c2ee01b) - AtomicFS
#### Tests
- add example to test nested output dir - (c6ff4b8) - AtomicFS

- - -

## v0.17.1 - 2025-04-15
#### Bug Fixes
- **(cmd)** convert version command to flag with improved help formatting - (206ce27) - AtomicFS
#### Build system
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (bfd36cb) - dependabot[bot]
- **(deps)** bump github.com/go-git/go-git/v5 in /cmd/firmware-action - (1a91224) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (682ef42) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (df30d3b) - dependabot[bot]
#### Miscellaneous Chores
- **(docker)** bump golang to v1.24 for u-root - (5d9471e) - AtomicFS
#### Refactoring
- **(cmd)** use T.Context for in tests - (fa92def) - AtomicFS
- **(cmd)** use T.Chdir for cleaner tests - (b9a68cc) - AtomicFS

- - -

## v0.17.0 - 2025-04-02
#### Bug Fixes
- **(cmd)** update detection of discontinued containers - (0a9cef4) - AtomicFS
#### Build system
- **(deps)** bump github.com/alecthomas/kong in /cmd/firmware-action - (f173548) - dependabot[bot]
- **(deps)** update pytest-cov requirement in /.dagger-ci/daggerci - (2bfe9db) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (0425abf) - dependabot[bot]
- **(deps)** bump github.com/go-playground/validator/v10 - (b65c794) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (3e73e88) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (7aa1fce) - dependabot[bot]
#### Documentation
- update paths in documentation - (ea6c140) - AtomicFS
#### Features
- **(goreleaser)** add Taskfile - (2139058) - AtomicFS
- **(goreleaser)** add nFPM to goreleaser - (816c294) - AtomicFS
- **(goreleaser)** add config to automatically load linux kernel modules - (21f78af) - AtomicFS
- **(goreleaser)** make defaults explicit - (0888d83) - AtomicFS
#### Miscellaneous Chores
- **(cmd)** improve logging in GitHub - (9978796) - AtomicFS
- changes according to review - (2c346eb) - AtomicFS

- - -

## v0.16.0 - 2025-03-25
#### Bug Fixes
- **(action)** trailing slash - (90984b9) - AtomicFS
- **(action)** pass debug into firmware-action executable - (8e1bf7c) - AtomicFS
- **(action)** upload artifact only for specified target - (a86f0cb) - AtomicFS
- **(action)** remove forward slashes from artifact name - (879ea76) - AtomicFS
- **(action)** typo - (a586ba6) - AtomicFS
- **(cmd)** typo - (c70e61e) - AtomicFS
- **(example)** undefined environment variable - (20ba823) - AtomicFS
- **(example)** remove windows an macos test - (6ca45b0) - AtomicFS
- release workflow - (68615ca) - AtomicFS
#### Build system
- **(deps)** bump golangci/golangci-lint-action from 6 to 7 - (0ae9038) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (f9c7fff) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (7afd080) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (3c5fcd9) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (64974ab) - dependabot[bot]
#### Continuous Integration
- speed up go-test - (e184893) - AtomicFS
- optimize when docker containers are build - (0378964) - AtomicFS
- disable scheduled container-cleanup - (04bada2) - AtomicFS
#### Features
- **(action)** make compile option work on the outside too - (54e799c) - AtomicFS
- **(action)** add optional input for version - (56c885a) - AtomicFS
- **(cmd)** move artifact and cache preparation into golang - (ed0c552) - AtomicFS
#### Miscellaneous Chores
- **(action)** simplify unpacking - (3f546be) - AtomicFS
- **(action)** cleanup after unpacking - (260b8af) - AtomicFS
- **(mdbook)** improve cache - (c8a056a) - AtomicFS
- add comment - (8751734) - AtomicFS
#### Refactoring
- **(cmd)** according to new golangci-lint v2 - (daeefef) - AtomicFS
#### Tests
- **(cmd)** improve logging - (107533c) - AtomicFS
- **(lint)** update golangci config to v2 - (f75f420) - AtomicFS
- reduce number of example variants - (a984ccd) - AtomicFS

- - -

## v0.15.0 - 2025-03-18
#### Bug Fixes
- **(action)** multiple configuration files in cache and artifact uploads - (26da39e) - AtomicFS
- **(ci)** update container cleanup workflow - (35d5d95) - AtomicFS
- **(ci)** coreboot 24.02 container was discontinued - (451ef19) - AtomicFS
- **(cmd)** generate-config u-boot was empty - (7009a28) - AtomicFS
- **(cmd)** prevent nil pointer dereference in DetectChanges method - (4879282) - AtomicFS
- **(cmd)** golang test for linux module - (d53c1bd) - AtomicFS
- **(example)** use shallow fetch - (75aa7fb) - AtomicFS
- **(example)** update when golang is compiled in examples - (701b57f) - AtomicFS
- **(example)** fail on release - (d0d7042) - AtomicFS
- **(python)** downloading docker-compose executable - (5c616d0) - AtomicFS
- release workflow - (986bec65) - AtomicFS
- update Taskfile - (771f70b) - AtomicFS
#### Build system
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (1011e11) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (314b8c2) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /cmd/firmware-action - (7af6129) - dependabot[bot]
- **(deps)** bump github.com/jedib0t/go-pretty/v6 - (ddd1295) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (12a238e) - dependabot[bot]
- **(deps)** bump github.com/go-git/go-git/v5 in /cmd/firmware-action - (f29f69b) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (04b95e8) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (26426ce) - dependabot[bot]
- **(deps)** bump github.com/google/go-cmp in /cmd/firmware-action - (fed74da) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (870af31) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (9ef4ae6) - dependabot[bot]
- **(deps)** update dagger-io requirement in /.dagger-ci/daggerci - (0112e79) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (5b518b8) - dependabot[bot]
#### Continuous Integration
- speed comparison between pruning and not pruning - (d14d821) - AtomicFS
- add prune example test - (092b6ed) - AtomicFS
#### Documentation
- update comment in action.yml - (12d1827) - AtomicFS
- add notes on migrating to new version - (12109cf) - AtomicFS
- update GitHub CI usage documentation - (4f26c24) - AtomicFS
- add notes on YAML multi-line strings - (4fd9898) - AtomicFS
- add artifact and caching documentation to GitHub CI usage - (a8bb3f4) - AtomicFS
- move examples into separate section - (c6baee1) - AtomicFS
- add note on discontinued container in docker-compose - (b488d50) - AtomicFS
- cosmetic changes - (dbcf1d2) - AtomicFS
- add documentation for git commit hash change detection - (96ef385) - AtomicFS
- cosmetic - (8c3c8ac) - AtomicFS
- add notes about change detection - (b5a15a3) - AtomicFS
- fix typo - (dfbcf29) - AtomicFS
- update shell-completions - (5a59061) - AtomicFS
- add citation-file - (1edea9f) - AtomicFS
- add CONVENTIONS.md and upadte CONTRIBUTING.md - (cf233a7) - AtomicFS
- support multiple configuration files - (e737afc) - AtomicFS
- add tips and tricks page - (f11698d) - AtomicFS
- add bug report template - (b2cafeb) - AtomicFS
- fix broken link - (fe305bb) - AtomicFS
#### Features
- **(action)** add input options to control action behavior - (eec0710) - AtomicFS
- **(action)** automatically upload artifacts - (a6fda90) - AtomicFS
- **(action)** automatically cache - (c103f63) - AtomicFS
- **(cmd)** add option to prune Dagge Engine - (fb86de0) - AtomicFS
- **(cmd)** add docker cleanup function - (bc381ef) - AtomicFS
- **(cmd)** large refactor for change detection - (9ebfb19) - AtomicFS
- **(cmd)** add ErrNotGitRepository into runGit - (7eee655) - AtomicFS
- **(cmd)** make GitDescribe more universal - (c962e1c) - AtomicFS
- **(cmd)** detect changed in configuration - (ad7d104) - AtomicFS
- **(cmd)** support multiple configuration files - (2c35324) - AtomicFS
- **(docker)** add script to switching to python2 into edk2 containers - (b1bc223) - AtomicFS
- **(example)** enable debug in CI - (9d1c501) - AtomicFS
- add validate-config command to CLI for config validation - (26b12f9) - AtomicFS
#### Miscellaneous Chores
- **(action)** improve CI logging by grouping output log - (33eb36b) - AtomicFS
- **(cmd)** better warning - (740dfdc) - AtomicFS
- **(cmd)** move GitHub CI detection into separate function - (e0bce25) - AtomicFS
- **(cmd)** cosmetic changes - (77e5e12) - AtomicFS
- **(cmd)** add debug message into AnyFileNewerThan - (330c67b) - AtomicFS
- **(cmd)** use Filenamify function for time-stamps too - (132ae26) - AtomicFS
- **(cmd)** cosmetic - (c778e1b) - AtomicFS
- **(cmd)** add Filenamify function into filesystem - (ec52f42) - AtomicFS
- **(cmd)** replace hardcoded version strings with ldflag - (aa105fc) - AtomicFS
- **(cmd)** formatter cleanup - (79aabdd) - AtomicFS
- **(cmd)** cleanup in config test - (f3b2ec6) - AtomicFS
- **(cmd)** add old linux containers into discontinued list - (13e81b4) - AtomicFS
- **(cmd)** sort discontinued containers alphabetically - (b8b2fcd) - AtomicFS
- **(docker)** add NodeJS to uboot and uroot containers - (562fbf7) - AtomicFS
- **(docker)** add new linux containers and remove old - (f910c3b) - AtomicFS
- **(example)** update example - (ab0964d) - AtomicFS
- **(example)** cleanup - (08ca77b) - AtomicFS
- **(linter)** cspell - (ca4955f) - AtomicFS
- **(linter)** cspell - (fc5f193) - AtomicFS
- **(linter)** cspell - (48f4d1e) - AtomicFS
- changes according to review - (e94ebb8) - AtomicFS
- cosmetic - (68de45b) - AtomicFS
- go mod tidy - (f6a76e4) - AtomicFS
- update .gitignore - (b1eec2e) - AtomicFS
- update goreleaser configuration - (33017f7) - AtomicFS
- fix typo - (9d0e2ba) - AtomicFS
- go mod tidy - (2b240bb) - AtomicFS
- remove temporary symlinks for defconfigs - (020e041) - AtomicFS
- add temporary symlinks for defconfigs - (d3d44ef) - AtomicFS
- update linux defconfigs - (4e85464) - AtomicFS
#### Refactoring
- **(action)** update action.yml multi-line string styles - (9418ccd) - AtomicFS
- **(cmd)** tweak AnyFileNewerThan function - (7943906) - AtomicFS
- **(cmd)** use reflection in Merge method - (4889bbd) - AtomicFS
- **(cmd)** use reflection in AllModules method - (69f5eda) - AtomicFS
#### Tests
- **(cmd)** add test for the new up-to-date detection - (601d332) - AtomicFS
- **(cmd)** use reflection in AllModules method - (c8d7636) - AtomicFS
- **(example)** update Linux matrix - (d3daf24) - AtomicFS
- **(lint)** add goreleaser check - (cf0a9d5) - AtomicFS
- update Taskfile to simulate GitHub - (844ec07) - AtomicFS
- add universal module into examples - (877de4e) - AtomicFS

- - -

## v0.14.1 - 2025-02-17
#### Bug Fixes
- **(cmd)** empty coreboot blob handling - (5d58179) - AtomicFS
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (0a47874) - dependabot[bot]
- **(deps)** bump github.com/go-playground/validator/v10 - (96a430f) - dependabot[bot]

- - -

## v0.14.0 - 2025-02-15
#### Bug Fixes
- **(cmd)** coreboot blobs directory vs file handling - (00be63f) - AtomicFS
- **(docker)** compilation of coreboot utils - (c37115b) - AtomicFS
- **(docker)** typo in edk2-stable202411 - (ccc219a) - AtomicFS
- **(docker)** edk2 missing branch - (17519e7) - AtomicFS
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (a94a109) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /cmd/firmware-action - (4d82bde) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (c09c7e4) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (7de8daa) - dependabot[bot]
#### Continuous Integration
- **(lint)** disable markdown-table-formatter - (45dd6d7) - AtomicFS
#### Documentation
- add notes on migrating to new coreboot config - (2c10f80) - AtomicFS
- add link to DockerHub containers - (871fb89) - AtomicFS
- add u-boot into list of supported modules - (25e291c) - AtomicFS
- update list of containers - (2572282) - AtomicFS
#### Features
- **(cmd)** refactor how coreboot blobs are handled - (d4437a0) - AtomicFS
#### Miscellaneous Chores
- **(docker)** add git wget and curl to all containers - (76c7a21) - AtomicFS
- **(docker)** add new edk2-stable202411 container - (6cca466) - AtomicFS
- **(linter)** cosmetic fixes - (56e3ee8) - AtomicFS

- - -

## v0.13.0 - 2025-02-11
#### Bug Fixes
- **(ci)** go-test - (8bd0a26) - AtomicFS
- **(cmd)** path validation in JSON - (7eaa619) - AtomicFS
- **(cmd)** incorrect module path - (f1be43f) - AtomicFS
- **(test)** example Taskfile for Linux - (c10087d) - AtomicFS
- trivy false positive AVD-DS-0001 - (80c9b1d) - AtomicFS
#### Build system
- **(deps)** bump github.com/jedib0t/go-pretty/v6 - (51bac62) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /cmd/firmware-action - (2e1b921) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (7f0cc81) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (3ab7832) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /cmd/firmware-action - (5b8c9e6) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (f966bba) - dependabot[bot]
#### Continuous Integration
- **(docker)** cleanup and update env vars - (6ba5218) - AtomicFS
- **(docker)** increase git commit sha slug to 12 characters - (b146d3d) - AtomicFS
- **(docker)** move container cleanup into separate workflow - (1ee5926) - AtomicFS
- **(labeler)** improve labeling for modules - (6126d79) - AtomicFS
- better caching for coreboot - (ad0ed57) - AtomicFS
- add u-boot example - (1e03ffe) - AtomicFS
- remove cache cleanup workflows - (82ef4b3) - AtomicFS
- cleanup coreboot toolchain PRs - (1acfd3a) - AtomicFS
#### Documentation
- add u-boot module config - (8a57478) - AtomicFS
- notes about discontinued containers - (1ecc76d) - AtomicFS
- update references to firmware-action-example - (f35bd79) - AtomicFS
- add notes about building contianers on the fly - (7861a25) - AtomicFS
- add link to toolchains repo - (5228032) - AtomicFS
- add link to firmware-action-example repo - (f1f0b73) - AtomicFS
#### Features
- **(ci)** automatically re-run failed container builds - (b7a842d) - AtomicFS
- **(cmd)** add support for u-boot - (a2b42d7) - AtomicFS
- **(cmd)** warning about using discontinued containers - (ed551ea) - AtomicFS
- **(cmd)** expose building contianers on the fly - (2af58ac) - AtomicFS
- **(docker)** add clang into uboot container - (6d351b4) - AtomicFS
- **(docker)** add support for DockerHub - (9421298) - AtomicFS
- **(docker)** add uboot container - (cf5668f) - AtomicFS
#### Miscellaneous Chores
- **(ci)** cleanup - (86e9cc6) - AtomicFS
- **(cmd)** fixup logging in tests - (5b573b8) - AtomicFS
- **(docker)** add new coreboot containers - (d56e2fa) - AtomicFS
- **(docker)** pass Taskfile CLI_ARGS to dagger-ci - (98831af) - AtomicFS
- **(docker)** update get_env_var_value - (cd3ced9) - AtomicFS
- **(docker)** remove ssh server from containers - (2ed69b8) - AtomicFS
- **(linter)** cspell - (24744df) - AtomicFS
- **(linter)** cspell - (55841a2) - AtomicFS
- **(linter)** cspell - (5dd031a) - AtomicFS
#### Tests
- **(example)** add u-boot into Taskfile for local run - (1600333) - AtomicFS

- - -

## v0.12.0 - 2025-01-27
#### Bug Fixes
- **(ci)** example runs on windows and macos - (b171694) - AtomicFS
#### Build system
- **(deps)** bump dagger.io/dagger in /cmd/firmware-action - (c840941) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (6aa77aa) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (82322c2) - dependabot[bot]
- **(deps)** bump github.com/go-git/go-git/v5 in /cmd/firmware-action - (c31eb83) - dependabot[bot]
#### Documentation
- update notes on interactive debugging - (a58aa67) - AtomicFS
#### Features
- **(cmd)** add DefconfigPath to list of sources - (2a109e2) - AtomicFS
#### Miscellaneous Chores
- **(cmd)** remove interactivity via SSH - (bf44f0d) - AtomicFS
- go mod tidy and drop toolchain pin - (76aa2b3) - Marvin Drees

- - -

## v0.11.1 - 2025-01-22
#### Bug Fixes
- release workflow - (5e845cd) - AtomicFS

- - -

## v0.11.0 - 2025-01-22
#### Bug Fixes
- **(action)** add coreboot blobs into GetSources - (0f0a9a6) - AtomicFS
- **(ci)** stop scheduled docker container builds - (8617665) - AtomicFS
- **(cmd)** passing coreboot version into container - (38c7b20) - AtomicFS
- **(docker)** downgrade NodeJS for udk2017 - (23c8d31) - AtomicFS
- **(docker)** install newer NodeJS into edk2 containers - (52d5cd8) - AtomicFS
- create subgraph when multiple roots present - (942bc1b) - Marvin Drees
- undefined variable in Taskfile - (bbc46be) - AtomicFS
#### Build system
- **(deps)** bump github.com/go-playground/validator/v10 - (c90f456) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto from 0.31.0 to 0.32.0 in /action - (5ac5e4a) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /action - (943d480) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (3d6789e) - dependabot[bot]
- **(deps)** bump github.com/go-git/go-git/v5 in /action - (cb6d824) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (f60fd1d) - dependabot[bot]
- **(deps)** bump github.com/go-git/go-git/v5 in /action - (e1ec222) - dependabot[bot]
#### Continuous Integration
- tweak release-prepare - (5d8cef1) - AtomicFS
#### Documentation
- fix changelog headline - (4fa835d) - AtomicFS
- fix dates in changelog - (01261dc) - AtomicFS
- unify firmare-action name in README - (c948e47) - AtomicFS
- update firmware-action installation instruction - (c760068) - AtomicFS
- update remaining '/action' reference - (1d6a465) - Marvin Drees
- add golang installation instructions - (ec9bbd1) - AtomicFS
#### Features
- **(cmd)** pass env vars into coreboot container - (16085ab) - AtomicFS
- **(cmd)** add functions to help with passing env vars into container - (1cefcf7) - AtomicFS
- **(cmd)** add functions for handling git describe - (2aebc88) - AtomicFS
#### Miscellaneous Chores
- **(action)** move golang code - (94b500f) - AtomicFS
- **(cmd)** unify assert parameters - (586199f) - AtomicFS
- **(docker)** bump coreboot docker containers from jammy to noble - (4ec7b9f) - AtomicFS
- **(linter)** cspell - (990ce12) - AtomicFS
- sort .gitignore alphabetically - (fcd2bf6) - AtomicFS
- change go mod path and imports - (be9a59a) - Marvin Drees
- cleanup after moving golang code - (07a5942) - AtomicFS
#### Refactoring
- **(cmd)** use DIY cache in /tmp for coreboot tests - (abdab2c) - AtomicFS
- **(cmd/test)** according to review - (6c14e86) - AtomicFS
- **(cmd/test)** use /tmp for cache in all tests - (4ef189c) - AtomicFS
- **(test)** test embedded coreboot version in compiled binary - (88ba083) - AtomicFS
#### Tests
- **(cmd)** passing of git describe into container - (1ccd2d4) - AtomicFS
- **(cmd)** orphaned node break DAG - (402447e) - AtomicFS

- - -


## v0.10.2 - 2024-12-19
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (4d59952) - dependabot[bot]
- **(deps)** bump github.com/jedib0t/go-pretty/v6 in /action - (9454cd6) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (c150bd1) - dependabot[bot]
#### Miscellaneous Chores
- **(dependabot)** split updates for golang - (27d42d8) - AtomicFS
- **(deps)** bump x/net due to an upstream CVE - (7823b44) - Marvin Drees
#### Tests
- **(action)** add timeout to SSH server start - (d27861e) - AtomicFS
- **(action)** update taskfile test - (8bd720b) - AtomicFS

- - -

## v0.10.1 - 2024-12-14
#### Build system
- **(deps)** bump golang.org/x/crypto from 0.30.0 to 0.31.0 in /action - (8e739ed) - dependabot[bot]
- **(deps)** update dagger-io requirement in /.dagger-ci/daggerci - (fd6eae2) - dependabot[bot]
#### Miscellaneous Chores
- **(action)** bump version to v0.10.1 - (5509437) - AtomicFS

- - -

## v0.10.0 - 2024-12-11
#### Bug Fixes
- **(action)** skipping .git - (f57cc93) - AtomicFS
- **(action)** check for existing non-empty output directory - (f1f9418) - AtomicFS
- **(action)** check for existing output directory - (04337ab) - AtomicFS
- **(action)** missing universal module - (e244f7c) - AtomicFS
- **(action)** generate-config missing uroot - (079883f) - AtomicFS
- **(docs)** fix missing content from submodules - (b4deb3a) - AtomicFS
#### Build system
- **(deps)** bump github.com/alecthomas/kong - (c6a22e4) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (2f219b4) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (958fbab) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto in /action in the golang group - (94f65f6) - dependabot[bot]
#### Continuous Integration
- **(lint)** markdown link check as warning - (1b4c19e) - AtomicFS
#### Documentation
- **(action)** fix missing docs for uroot - (1d216ff) - AtomicFS
- **(action)** improve configuration overview - (be9538d) - AtomicFS
#### Features
- **(action)** add source changes detection - (3df87c3) - AtomicFS
- **(action)** add functions to detect changes in files - (c88853b) - AtomicFS
- **(action)** add universal module - (3318c85) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.10.0 - (d4beebe) - AtomicFS
- **(action)** cosmetic fixes - (ad462d3) - AtomicFS
- **(action)** prettify summary table - (f5d8ae3) - AtomicFS
- **(action)** speed up filesystem.DirTree - (81e0b94) - AtomicFS
- **(action)** fix typo - (ad14dfc) - AtomicFS
- **(action)** run formatter - (e212627) - AtomicFS
- **(action)** cosmetic fixes - (ed72d9c) - AtomicFS
#### Tests
- **(action)** add universal go-test - (34658bd) - AtomicFS

- - -

## v0.9.0 - 2024-12-03
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (a91d6a9) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong - (8a4052d) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong - (319583c) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (c73ac3c) - dependabot[bot]
#### Continuous Integration
- **(megalinter)** re-enable trivy - (17b231c) - AtomicFS
#### Documentation
- **(action)** print digest sha for use container - (089db89) - AtomicFS
- **(action)** add link to source-code into --help - (49b2cb2) - AtomicFS
- update git submodule update suggestion - (cc644e2) - AtomicFS
- update dagger troubleshooting link - (8df641b) - AtomicFS
- add common troubleshooting steps - (589dff0) - AtomicFS
- add notes on offline usage - (9e596f6) - AtomicFS
#### Features
- **(action)** check for any undefined environment variable in config - (af1bc10) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.9.0 - (7a94195) - AtomicFS
- **(action)** on dagger error check of kernel module - (cd01c93) - AtomicFS
- **(action)** run formatter - (a002d9c) - AtomicFS

- - -

## v0.8.1 - 2024-11-28
#### Bug Fixes
- **(ci)** next version calculation - (ce6e90f) - AtomicFS
- **(docker)** udk2017 needs python2 - (c3c3f22) - AtomicFS
#### Documentation
- **(action)** improve error message on missing coreboot blob - (6dae629) - AtomicFS
- **(contributing)** update - (e72fb0f) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.8.1 - (fc387d4) - AtomicFS

- - -

## v0.8.0 - 2024-11-27
#### Bug Fixes
- **(ci)** cache cleanup - (1a4af0c) - AtomicFS
- **(ci)** run example tests on change in golang code - (b60076b) - AtomicFS
- **(ci)** automerge on pull_request_target instead of pull_request - (b946bcf) - AtomicFS
- **(ci)** allow automerge to fail - (46c5296) - AtomicFS
- **(ci)** replace only first semver in Taskfile - (bd1063d) - AtomicFS
- **(docker)** forgotten debug variable - (c5e7f0e) - AtomicFS
- **(linter)** cspell, editorconfig and shellcheck - (db36238) - AtomicFS
- **(taskfile)** tweaks - (5619f15) - AtomicFS
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (5cad9d3) - dependabot[bot]
- **(deps)** bump github.com/stretchr/testify - (3491a1b) - dependabot[bot]
#### Documentation
- **(ci)** update comment - (190b283) - AtomicFS
- add notes on containers - (8b6cb99) - AtomicFS
#### Features
- **(action)** add shell completion - (8fe18fe) - AtomicFS
- **(action/coreboot)** add support for 10gbe blob - (66619f0) - AtomicFS
- **(docker)** download toolchains from separate repository - (109fb8e) - AtomicFS
- **(docker)** store toolchains in separate repository - (2032b7e) - AtomicFS
- **(docker)** do not pre-compile utilities - (b363acc) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.8.0 - (495b125) - AtomicFS

- - -

## v0.7.0 - 2024-11-22
#### Bug Fixes
- **(action)** infinite symlink issue - (3d9a259) - AtomicFS
- **(action)** linux make defconfig file conflict - (50e1d3d) - AtomicFS
- **(action/linux)** missing dir - (1d824a9) - AtomicFS
- **(action/linux)** fix cross-compilation detection - (5e9f8f3) - AtomicFS
- **(ci)** compile firmware-action in examples - (4e88c4b) - AtomicFS
- **(ci)** tweak dynamically generated matrix - (761da59) - AtomicFS
- **(ci)** update dagger use in python - (c9f52e4) - AtomicFS
- **(ci)** reminder bot - (6074ee4) - AtomicFS
- **(ci/automerge)** use PAT token - (dd4f88b) - AtomicFS
- **(ci/cleanup)** do not fail when cache not found - (2d9c722) - AtomicFS
- **(dagger)** run python black formatter - (fd3d412) - AtomicFS
- **(dagger)** install build-essential package - (540b4fd) - AtomicFS
- **(docker)** apply patches in tests too - (5e01648) - AtomicFS
- **(docker)** edk2-stable202408 missing submodule - (88151a1) - AtomicFS
- **(docker)** add cross toolchain for x86 into linux container - (b023bf9) - AtomicFS
- **(docker)** linux container missing package on arm64 - (39083ee) - AtomicFS
- **(docker)** add omitted python into compose.yaml - (bbdcfa1) - AtomicFS
- **(docker)** typo in compose.yaml - (9203d56) - AtomicFS
- **(docker)** bump all linux docker containers to noble release - (8573152) - AtomicFS
- **(docker)** update GCC vs GCC5 also in examples - (22450d2) - AtomicFS
- **(docker)** enable again tests of containers - (b91a043) - AtomicFS
- **(docker)** download latest version of docker-compose - (8a85b80) - AtomicFS
- **(docker)** slim down edk2 container - (5649f08) - AtomicFS
- **(docker/edk2)** shallow submodules for edk2-stable202008 - (a46ffb7) - AtomicFS
- **(docker/edk2)** shallow submodules - (54bfeda) - AtomicFS
- **(docker/edk2)** possibly broken docker arguments - (607fa71) - AtomicFS
- **(docs)** inter-linking was broken - (e5289aa) - AtomicFS
- **(edk2)** toolchain GCC5 was deprecated in edk2-stable202305 - (f900fb2) - AtomicFS
- **(examples)** fixed typo - (a6c6001) - AtomicFS
- **(examples)** partial revert of 32583f79 - (77759c3) - AtomicFS
- **(examples)** conflicting artifact names - (84f4916) - AtomicFS
- **(examples)** make arch into env variable - (e95fabc) - AtomicFS
- **(examples)** artifact names - (41cb3ae) - AtomicFS
- **(examples)** path to artifacts - (fdf1df3) - AtomicFS
- **(linter)** revive fixes in golang - (9b40826) - AtomicFS
- **(linter)** pylint fixes in python - (bbcbe25) - AtomicFS
- **(megalinter)** fix spelling - (5d64dd4) - AtomicFS
- **(megalinter)** fix spelling - (57ea19c) - AtomicFS
- **(taskfile)** cleanup - (8636a55) - AtomicFS
- **(tests)** forgotten env variable - (81a433e) - AtomicFS
#### Build system
- **(deps)** bump docs/src/firmware-action-example - (dca2df1) - dependabot[bot]
- **(deps)** bump github.com/go-playground/validator/v10 - (64fe7bc) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (d5c0c5c) - dependabot[bot]
- **(deps)** update pytest-flake8 requirement in /.dagger-ci/daggerci - (f4362ea) - dependabot[bot]
- **(deps)** update dagger-io requirement in /.dagger-ci/daggerci - (bdb564a) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto in /action in the golang group - (82dac48) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong - (9dad109) - dependabot[bot]
- **(deps)** bump docs/src/firmware-action-example - (8191bcf) - dependabot[bot]
- **(deps)** bump the golang group across 1 directory with 2 updates - (289e338) - dependabot[bot]
- **(deps)** update pytest-benchmark requirement in /.dagger-ci/daggerci - (9bb6963) - dependabot[bot]
- **(deps)** update prettytable requirement in /.dagger-ci/daggerci - (fac4f9f) - dependabot[bot]
- **(deps)** update pytest-cov requirement in /.dagger-ci/daggerci - (f024ffc) - dependabot[bot]
- **(deps)** update pytest-benchmark requirement in /.dagger-ci/daggerci - (02e489e) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (8be1bb1) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (9c5fa6f) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (51cfe2a) - dependabot[bot]
- **(deps)** bump golang.org/x/crypto in /action in the golang group - (745b5c7) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (f095979) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (6db7264) - dependabot[bot]
- **(deps)** update anyio requirement in /.dagger-ci/daggerci - (aa84ffa) - dependabot[bot]
- **(deps)** update dagger-io requirement - (b01f5bf) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (1990387) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong - (7533e02) - dependabot[bot]
- **(deps)** bump the golang group in /action with 2 updates - (2e09d8e) - dependabot[bot]
- **(deps)** bump github.com/alecthomas/kong in /action - (64039cf) - dependabot[bot]
- **(deps)** bump the golang group across 1 directory with 3 updates - (b1f41cb) - dependabot[bot]
- **(deps)** bump peter-evans/create-pull-request from 6 to 7 - (afd39d9) - dependabot[bot]
- **(deps)** bump dagger.io/dagger in /action in the golang group - (1818d4a) - dependabot[bot]
- **(deps)** bump github.com/sethvargo/go-githubactions - (09d8897) - dependabot[bot]
#### Continuous Integration
- **(automerge)** do not run on draft PR - (d7463b6) - AtomicFS
- **(cache)** improve caching in example workflow - (8de8ef8) - AtomicFS
- **(cache)** add cleanup jobs - (4ebd4df) - AtomicFS
- **(docker)** dynamically generate the matrix - (0059567) - AtomicFS
- **(docker)** build monthly instead of weekly - (5461bec) - AtomicFS
- **(docker)** skip testing - (674a5c9) - AtomicFS
- **(docker)** add timeouts - (720b7ac) - AtomicFS
- **(megalinter)** temporarily disable trivy because - (f024bc6) - AtomicFS
- **(reminder)** run daily instead of hourly - (606b105) - AtomicFS
- a bit more automation to pull requests - (5629978) - AtomicFS
- fix cache cleanup - (4d0c1fb) - AtomicFS
#### Documentation
- apply suggestions from code review - (9216098) - AtomicFS
- update documentaion - (fb985f1) - AtomicFS
- add SECURITY.md - (ab781d4) - AtomicFS
- cosmetic changes according to megalinter - (e2be0bd) - AtomicFS
- add notes about creating new containers - (c22deb4) - AtomicFS
#### Features
- **(action)** build edk2 basetools on the fly - (f92e178) - AtomicFS
- **(ci)** fetch coreboot toolchains - (92907ad) - AtomicFS
- **(ci)** add 2nd part to reminder bot - (5743166) - AtomicFS
- **(ci)** use arduino/setup-task to get taskfile - (78c19f0) - AtomicFS
- **(dagger)** make containers multi-arch - (f2abe55) - AtomicFS
- **(dependabot)** use also for git submodules - (ea075ac) - AtomicFS
- **(docker)** apply patches to edk2 containers - (4486886) - AtomicFS
- **(docker)** add edk2 containers between 2023-02 and 2024-05 - (22b5d11) - AtomicFS
- **(docker)** bump GCC_VERSION for edk2-202408.01 - (1750ca3) - AtomicFS
- **(docker)** use pre-compiled coreboot toolchain in dockerfile - (9c3128f) - AtomicFS
- **(docker)** compile coreboot tool-chains separately - (c608076) - AtomicFS
- **(docker)** make coreboot container support multi-arch builds - (c91219c) - AtomicFS
- **(docker)** make edk2 container support multi-arch builds - (9e361bb) - AtomicFS
- **(docker)** make linux container support multi-arch builds - (4fe8fd2) - AtomicFS
- **(docker)** add zstd into linux container - (c48e1f4) - AtomicFS
- **(docker)** expand fleet of Linux containers - (2d9178b) - AtomicFS
- **(docker)** add missing packages - (c4133f5) - AtomicFS
- **(docker)** add new edk2 202408 container - (2d8fa0e) - AtomicFS
- **(taskfile)** include docker building taskfile in the main one - (e48b996) - AtomicFS
- **(taskfile)** add tasks to build containers - (d627efd) - AtomicFS
- **(taskfile)** add python virtual environment setup task - (b103b16) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.7.0 - (6427558) - AtomicFS
- **(action)** add debug for arch inputs - (22afde5) - AtomicFS
- **(action)** unify architecture strings - (32583f7) - AtomicFS
- **(action)** remove obsolete code - (6c8fc70) - AtomicFS
- **(action)** run go mod tidy - (738f44f) - AtomicFS
- **(ci)** correction to example caching - (b51b3a2) - AtomicFS
- **(ci)** remove leftover junk in go-test workflow - (aa4fa83) - AtomicFS
- **(dependabot)** remove docker - (3c33923) - AtomicFS
- **(dependabot)** remove reviewers and assignees - (a31adf9) - AtomicFS
- **(docker)** depreciating coreboot 24.02 in favour of 24.02.01 - (273cadf) - AtomicFS
- **(docker)** update gitignore - (5f13a2f) - AtomicFS
- **(docker)** remove unnecessary dependencies - (52366f7) - AtomicFS
- **(docker)** add omitted arguments into compose file - (5062b12) - AtomicFS
- **(docker)** add omitted arguments into compose file - (bbc194c) - AtomicFS
- **(docker)** cleanup python code - (e107818) - AtomicFS
- **(python)** switch from hardcoded version to latest stable - (e71a9eb) - AtomicFS
#### Revert
- **(ci)** remove edk2 matrix - (a849e26) - AtomicFS
#### Style
- cosmetic change in edk2 test script - (8ebee2a) - AtomicFS
#### Tests
- **(docker)** expand examples to include new Linux containers - (da84266) - AtomicFS
- **(docker)** test containers also on change of tests - (54cdfc6) - AtomicFS
- **(examples)** add taskfile to run Linux example locally - (630e8aa) - AtomicFS
- **(examples)** add cleanup step - (ea35ff3) - AtomicFS
- **(examples)** expand to also run on arm64 machine - (9319dfa) - AtomicFS

- - -

## v0.6.1 - 2024-08-27
#### Bug Fixes
- **(action/linux)** defconfig filename - (4b9e9d4) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to v0.6.1 - (e9fdd48) - AtomicFS

- - -

## v0.6.0 - 2024-08-26
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

## v0.5.0 - 2024-08-01
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

## v0.4.0 - 2024-07-30
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

## v0.3.2 - 2024-07-11
#### Features
- **(action)** allow multi-module workspaces for u-root - (f54803d) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.3.2 - (4cac34c) - AtomicFS

- - -

## v0.3.1 - 2024-07-11
#### Bug Fixes
- **(again)** build docker containers on release - (2d33a7e) - AtomicFS
#### Miscellaneous Chores
- **(action)** bump version to 0.3.1 - (ca343ff) - AtomicFS

- - -

## v0.3.0 - 2024-07-11
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

## v0.2.1 - 2024-04-12
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

## v0.2.0 - 2024-03-11
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

## v0.1.2 - 2024-03-04
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

## v0.1.1 - 2024-03-04
#### Bug Fixes
- **(action)** naming mistake in JSON config - (58ebf67) - AtomicFS

- - -

## v0.1.0 - 2024-01-31
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

## v0.0.1 - 2024-01-22
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


