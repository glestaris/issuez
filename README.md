# Issuez

Importing Tickets to Jira from Markdown.

Inspired by [Onsi's Prolific](https://github.com/onsi/prolific).

## Installation

If you have Go installed, you can install by running:

```
$> go get -u github.com/glestaris/issuez
```

Otherwise, you may use one of the static binaries found in releases.

## Capabilities

`issuez` is a CLI that can be configured to connect to a JIRA instance and
import issues from a markdown file. The markdown file looks as follows:

```
[Task] Task title

Task **description** in markdown.

Epic: EPIC-123
Labels: label-a, label-b

---

[Bug] Found a bug

Markdown description supports lists:

1. A
1. B
1. C

Epic: EPIC-123

---

This is a story

_As a user, ..._

Code blocks work too:

\`\`\`python
x = 12
\`\`\`

E: EPIC-123
L: my-label
```

Run `issuez` using the following flags:

```
$> ./issuez --api https://foo.atlassian.com/jira/ --username foo@gmail.com --token abc123 import --project-key PROJ ./issues.md
```

The flags are used as follows:

- `--api` or `-a`: The API endpoint URL.
- `--username` or `-u`: Your Jira username.
- `--token` or `-t`: Your Jira API token. This can be generating by looking at
  Account Settings > Security > API Tokens.
- `--project-key` or `-p`: The project you want the issues to be imported in.

## Contributing

### Building the tool

Simply run:

```
$> make issuez 
```

It will produce `./issuez` in your local directory.

### Running the tests

Use `./hack/test.sh` or `make` to run tests:

```
$> make test

== UNIT TESTS ==================
...

== INTEGRATION TESTS ===========
...

== ACCEPTANCE TESTS ============
...
```

`richgo` is required in order to run the tests. Install it by running:
`go get -u github.com/kyoh86/richgo`.
