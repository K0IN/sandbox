
# Sandbox

Please note this is a work in progress, there is nothing to see here.

## Install

> curl -sSL https://raw.githubusercontent.com/K0IN/sandbox/main/install.sh | sh

which will install the `sandbox` command to `/usr/local/bin`.

## Usage

> sandbox --help

## Examples

Start a sandbox of your current system:

> sandbox try "<your command>"

Start a container (pull from docker hub):

> sandbox container -i library/python:latest
