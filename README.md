# Cloudbit CLI

This command line interface serves as an additional frontend for the [Cloudbit](https://my.cloudbit.ch) platform.

![Tests](https://github.com/cloudbit-ch/cli/workflows/test/badge.svg)

## Installation

If you have GoLang installed, you can download and install the CLI with 

```shell script
go install github.com/cloudbit-ch/cli/v2/cmd/cloudbit@latest
```

Otherwise, you will need to download the executable for your system from the
release section in the github repository.

## Usage

After downloading you first need to authenticate the cli with using application
token. You can get a new token by creating one in the [Cloudbit](https://my.cloudbit.ch/#/organization/applications) 
portal.

Once you have a token, you need to set it up in the cli. You can do this by
creating a `.cloudbit/config.json` file in your home directory with the following
content:

```json
{
  "token": "YOUR_TOKEN_HERE"
}
```

Alternatively, you can pass the token as an argument to the cli with the
`--token` flag or by setting the `CLOUDBIT_TOKEN` environment variable.

Once you have successfully logged in into your account, you can start
manipulating things in your organization. As a first step it would be a good
idea to upload your personal ssh key onto our platform. You will need this for
every linux virtual machine you deploy. 
```shell script
cloudbit compute key-pair create \
    --name my-key-pair \
    --public-key ~/.ssh/id_rsa.pub
```

Just to test things out, you can try creating an ubuntu virtual machine using
the previously uploaded key pair:
```shell script
cloudbit compute server create \
    --name my-server \
    --location bit1 \
    --image ubuntu-20.04 \
    --product b1.1x1 \
    --key-pair my-key-pair
```

Further usage manuals can be found in the application itself using the `-h` or
`--help` flags.