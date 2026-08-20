### Project overview

Discord bot currently in use by the [FreeBSD Discord server](https://wiki.freebsd.org/Discord/DiscordServer).
The objective is to provide a bridge between the different FreeBSD
services (Bugzilla, Phabricator, Git, etc) and the Discord platform
itself while also introducing preventative measures to combat
spam/malicious users.

So far implemented is a handful of [commands](internal/bot/handler/command),
[event handlers](internal/bot/handler/event) and a [webhook](internal/pulse/hook/git)
that forwards commits of the FreeBSD GitHub repositories to Discord.

| Command | Description |
| ------- | ----------- |
| !help | Displays available commands |
| !status | Shows bot status |
| !role <role> | Allows users the ability to assign themselves roles |
| !bug <id> | Sends a message embed detailing a problem report from Bugzilla. Additionally, messages matching the FreeBSD Bugzilla URL will trigger this event |
| !review <id> | Sends a message embed detailing a Differntial revision from Phabricator. Additionally, messages matching the FreeBSD Phabricator URL will trigger this event |
| !user <id> | Sends a message embed detailing a user |

Key events on Discord including message updates, deletions, member
removals and bans are logged in a public channel to ensure transparency
within our community. Recently, we've seen users attempting to promote
malicious advertisements or spam channels. To combat this, we have
implemented an "antispam" measure to help identify and reduce these
issues as they arise.

While it's not possible to create heuristics that cover every type
of behavior, we make a basic attempt to identify the most significant
offenders and take appropriate action. In this repository, you will
find `bot.antispam.rules` in the default [YAML file](config.example.yaml)
that outlines common patterns of spam along with associated timeout values.
The goal is to expand this file over time to address more advanced cases
effectively.

### Building and deployment

Before proceeding to build anything ensure a valid configuration
file exists in the root of the project. Example can be found
[here](config.example.yaml).

`go` and `gmake` must be installed to build the project. Optionally,
`golangci-lint` for linting the code.

<details open>
<summary>Container image (recommended)</summary>

Optionally, you can build OCI images and deploy through `podman`
or `docker`.

```console
# make container
```

Once successfully built run the images as follows, passing the
`config.yaml` configuration file as a volume mount, replacing
`$HASH` with the according git sha:

```console
# podman run --rm -v ./config.yaml:/app/config.yaml:ro localhost/pulsar-bot:$HASH
# podman run --rm -v ./config.yaml:/app/config.yaml:ro localhost/pulsar-relay:$HASH
```

Container images are automatically [published to GitHub](https://github.com/lcook?tab=packages&repo_name=pulsar)
on each commit passing the build pipeline. Like above, run the
following:

```console
# podman run --rm -v ./config.yaml:/app/config.yaml:ro ghcr.io/lcook/pulsar/bot:$HASH
# podman run --rm -v ./config.yaml:/app/config.yaml:ro ghcr.io/lcook/pulsar/relay:$HASH
```
</details>

<details>
<summary>Manually building</summary>
Run:

```console
# make build
```

This will build both the `bot` and `relay` Go binaries.

</details>

### License

[BSD 2-Clause](LICENSE)
