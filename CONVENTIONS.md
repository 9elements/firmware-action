# Conventions


## File formatting

To maintain consistent coding styles, indentation, line length and so on we are using [.editorconfig](https://editorconfig.org/).


## Coding style

For golang code, please use `gofumpt` (a strict formatter for Go language, stricter than gofmt) before contributing. On top of that please fix all reported issues by linters such as `revive`, `go vet`, `staticcheck` and `golangci-lint` (up-to-date list of used linters is in `Taskfile.yml` in `lint` task.


## Guidelines for writing good functional code

When writing code, you MUST follow these principles:
- Code should be easy to read and understand.
- Keep the code as simple as possible. Avoid unnecessary complexity.
- Use meaningful names for variables, functions, etc. Names should reveal intent.
- Functions should be small and do one thing well. They should not exceed a few lines.
- Function names should describe the action being performed.
- Prefer fewer arguments in functions. Ideally, aim for no more than two or three.
- Only use comments when necessary, as they can become outdated. Instead, strive to make the code self-explanatory.
- When comments are used, they should add useful information that is not readily apparent from the code itself.
- Properly handle errors and exceptions to ensure the software's robustness.
- Consider security implications of the code. Implement security best practices to protect against vulnerabilities and attacks.


## Commit message guidelines

Please follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.


### Developer Sign-Off

For purposes of tracking code-origination, we follow a simple sign-off process. If you can attest to the [Developer Certificate of Origin](https://developercertificate.org/) then you append in each git commit text a line such as:

```
Signed-off-by: Your Name <username@youremail.com>
```


### Use of AI generated content

Use of purely AI generated or partially AI assisted content is permitted, however please mark it as AI generated in the commit message body by adding:
```
AI-Generated: true
AI-Model: <model used>
```
For example:
```
feat(...): add new feature

- adding new fun thing

AI-Generated: true
AI-Model: ChatGPT o3-mini
Signed-off-by: ...
```

This must follow [git trailer convention](https://git-scm.com/docs/git-interpret-trailers).
