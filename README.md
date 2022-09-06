## Pulsar Discord bot

Currently used in the FreeBSD Discord server to
relay incoming GitHub web-hook events to a desired
channel, showing information such as repository,
commit title and committer name.

## Tentative goals

Broader reach to Bugzilla and possibly Phabricator
events, or any additional services that serves
us of useful information.

## Build and deployment

Ensure a valid configuration file exists in the root
of the project. Example can be found [here](config.example.yaml).

```console
# make install
```

This will both build and install the resulting Go binary,
as well as the configuration file. An rc service script
comes included so that pulsar can be daemonized.
To enable the service, run:

```console
# sysrc pulsar_enable=YES
# service pulsar start
```

If you want to use a custom configuration file separate
of the global one (residing under `/usr/local/etc/pulsar`)
then pass the `-c` flag, followed with a desired absolute path.

Alternatively, specify the configuration that the rc service
uses:

```console
# sysrc pulsar_config=/path/to/config.yaml
```

## License

[BSD 2-Clause](LICENSE)
