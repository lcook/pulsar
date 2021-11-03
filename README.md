## Pulseline Discord bot

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
as well as the configuration file. To run, simply:

```console
# ./pulseline
```

If you want to use a custom configuration file separate
of the global one (residing under `/usr/local/etc/pulseline`)
then pass the `-c` flag, followed with a desired absolute path.

## License

[BSD 2-Clause](LICENSE)
