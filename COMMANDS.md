## General command overview

Short summaries of what each command does.

### bug #<######>

Displays information associated with a given Bugzilla problem
report using its ID (`bugs.freebsd.org/######`). We scan through
the whole contents of a message, and, matching a specific regex
(`bug\s#(?P<id>\d{1,6})`, e.g. `bug #249813`) will dispatch an
event handler to check whether the ID is valid and report back
the found result(s).

In particular, we show: status, product, component, proceeded
by a hyperlink to the bug itself shown as the report summary 
and lastly who created the problem report with creation date.

You may find the template residing [here](internal/bot/command/templates/report.tpl).

### role <name>

Users can self-assign certain roles as defined in [roles.json](internal/bot/command/data/roles.json)
simply by typing `!role <name>`, and type `!role` on it's
own to display what roles are available.
