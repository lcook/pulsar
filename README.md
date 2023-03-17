## Pulsar Discord bot

Currently used in the FreeBSD Discord server to
relay incoming GitHub web-hook events to a desired
channel, displaying information such as repository,
commit title and committer name in a nicely formatted
embedded message, allowing anyone to quickly view
recent changes in a centralized channel.

Among other user accessible commands, brings varying
aspects of the FreeBSD Project to Discord.

## Tentative goals

Broader reach to Bugzilla and possibly Phabricator
events, or any additional services that serves
us of useful information.

## Build and deployment

`go` and `bmake` must be installed to build the
project. Optionally, `golangci-lint` for linting.

An assumption is made that pulsar is built and ran
primarily on a FreeBSD host, so the provided RC
service script will not run properly (or at all) on
a different type of system (e.g., GNU/Linux), as
such, no guarantee is made for anything other than
FreeBSD. However, you are able to as least run the
application alongside a configuration file. If for
any reason you are using a non-default local prefix
where applications, configurations and friends get
installed, amend the `PREFIX` variable whilst building,
so for example: `make PREFIX=/opt build`.

To get started, ensure a valid configuration file
exists in the root of the project. Example can be
found [here](config.example.json).

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
# sysrc pulsar_config=/path/to/config.json
```

## License

[BSD 2-Clause](LICENSE)