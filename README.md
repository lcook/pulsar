## Pulsar Discord bot

Currently used in the FreeBSD Discord server to
forward GitHub webhook events to a chosen
channel, embedding commit events from various
repositories in a central manner. Also included
are a number of utility commands as described in
[COMMANDS.md](COMMANDS.md).

## Tentative goals

Broader reach to Bugzilla and possibly Phabricator
events, or any additional services that serves
us of useful information.

## Build and deployment

`go` and `bmake` must be installed to build the
project. Optionally, `golangci-lint` for linting.

To get started, ensure a valid configuration file
exists in the root of the project. Example can be
found [here](config.toml.example).

```console
# make install
```

This will both build and install the resulting Go binary,
as well as the configuration file. An RC service script
comes included so that pulsar can be daemonized.
To enable the service, run:

```console
# sysrc pulsar_enable=YES
# service pulsar start
```

If you want to use a custom configuration file separate
of the global one (residing under `/usr/local/etc/pulsar`)
then pass the `-c` flag, followed with a desired absolute path.

Alternatively, specify the configuration that the RC service
uses:

```console
# sysrc pulsar_config=/path/to/config.toml
```

### Container images

Optionally, you can build container images through
the use of `podman`. This is as simple as running:

```console
# make container
```

Once sucessfully built, run the image as follows,
passing the `config.yaml` configuration file as a
volume mount, replacing `$HASH` with the according git
sha:

```console
# podman run localhost/pulsar:$HASH -v ./config.toml:/app/config.toml /app/pulsar
```

Container images are automatically [published to GitHub](https://github.com/lcook/pulsar/pkgs/container/pulsar)
on each successful commit passing the build pipeline. Like
above, run the following:

```console
# podman run ghcr.io/lcook/pulsar:$HASH -v ./config.toml:/app/config.toml /app/pulsar
```

## License

[BSD 2-Clause](LICENSE)
