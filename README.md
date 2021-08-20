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

# Upload local git repo to Terraform Cloud and set the Terraform Working Directory to the path, relative to git root
tfd upload-config --path /path/to/terraform/dir --workspace myworkspace

# Start a run
tfd run start --workspace myworkspace

# Start a run and auto approve it once it is ready 
tfd run start --workspace myworkspace --auto-approve

# Start a destroy run
tfd run destroy --workspace myworkspace

# Apply the current run that is waiting for approval
tfd run apply --workspace myworkspace

# Stop queued runs and the current run
tfd run stop --workspace myworkspace
```

## Configure

Each of the command-line options can be read from env, config file (default: ~/.tfd.yaml) or as an option. Env vars will
be prefixed with "TFD_"

## Install

```sh
go get github.com/logandavies181/tfd@v1.0.0
```
