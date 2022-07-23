# tfd

Terraform Developer tool for interacting with Terraform Cloud and Terraform Enterprise

## Usage

```sh
# Set Terraform Cloud Token https://www.terraform.io/docs/cloud/users-teams-organizations/api-tokens.html
export TFD_TOKEN=xxx
export TFD_ORG=myorg
```

```sh
# List workspaces that you have access to in the current Org
tfd list-workspaces

# Upload local git repo to Terraform Cloud
tfd upload-config --path /path/to/project --workspace myworkspace

# Start a run
tfd run start --workspace myworkspace

# Start a run and auto approve it once it is ready
tfd run start --workspace myworkspace --auto-apply

# Start a destroy run
tfd run destroy --workspace myworkspace

# Apply the current run that is waiting for approval
tfd run apply --workspace myworkspace

# Stop queued runs and the current run
tfd run stop --workspace myworkspace
```

`tfd <subcommand> --help` for more info

## Configure

Each of the command-line options can be read from environment variables, config file (default: ~/.tfd.yaml) or as
command-line options. Environment variables are prefixed with "TFD_"

## Install

```sh
go get github.com/logandavies181/tfd@latest
```

Or check out [releases](https://github.com/logandavies181/tfd/releases)

## Contributing

Feel free to raise a PR or create an issue

## Developing

Ensure you have an up-to-date go compiler.

Before raising a Pull Request, ensure the unit tests pass and that the mocks are up to date by running
`generate_mocks.sh`

To create a release, run `goreleaser release`
