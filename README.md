---
confluence: 458831
title: Confluence Pages Creator GitHub Action
toc: true
---

# Confluence Pages Creator GitHub Action

This GitHub Action scans your repository for OpenAPI definition files and Markdown documentation to automatically
update Confluence pages.

![GitHub Action CI](https://github.com/danaketh/ga-confluence-gen/workflows/CI/badge.svg)

## Features

- Scans your repository for OpenAPI `.yaml` or `.json` files.
- Scans your repository for Markdown documentation files.
- Automatically updates Confluence pages with the found documentation.

## Requirements

- A running instance of Confluence that is accessible from GitHub Actions.
- A Confluence account with permissions to update pages.

## Usage

Add the following steps at the end of your `.github/workflows/main.yml` or whatever your workflow file is:

```yaml
jobs:
  create-Confluence-pages:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout Code
      uses: actions/checkout@v2
      
    - name: Create Confluence Pages
      uses: your-username/confluence-pages-creator-action@v1
      with:
        confluence-domain: my-domain
        confluence-username: gh-user
        confluence-token: very-secret-token
        confluence-space: FOO
        source-paths: |-
            docs/
            api/
```

## Configuration

I highly recommend you to create a Confluence account specifically for this action and use it to generate an API token.
This way, you can restrict the permissions of the account to only update pages.

### Action Inputs

| Input                 | Description             | Required | Default |
|-----------------------|-------------------------|----------|---------|
| `confluence-domain`   | Confluence Instance URL | `true`   |         |
| `confluence-username` | Confluence Username     | `true`   |         |
| `confluence-token`    | Confluence API Token    | `true`   |         |
| `confluence-space`    | Confluence Space Key    | `true`   |         |
| `source-paths`        | Source Paths            | `false`  | `.`     |


## How it Works

1. The action first scans the repository for any OpenAPI `.yaml` or `.json` files.
2. It then scans for Markdown files.
3. The action then interacts with the Confluence API to update pages based on the found files.

### Markdown documentation

Your markdown files can be placed anywhere in the repository. The action will scan the entire repository for them.
To avoid conflicting with other GitHub Actions or whatever documentation files you may have in your repository,
only files which container front matter are considered.

Currently, the action only supports TOML and YAML front matter. The front matter must contain the `confluence` key
which contains the ID of the page in Confluence. I could not find a way to keep some kind of relationship between
the documentation file and the Confluence page (like storing hash in metadata or something like that), thus
it is necessary to first create the page in Confluence and then add the ID to the front matter.

While a bit annoying, this gives you the flexibility to create the page in Confluence however you want and then
just bind it to the documentation file.

#### Markdown support

Confluence uses their own markup language for formatting. To make things easier, this action converts the Markdown
to HTML and then modifies the HTML to match the Confluence markup. This is not perfect and some things may not
work as expected. If you find any issues, please open an issue or a pull request.

In most cases it's done by regex. You're welcome to open a pull request with a better solution.

#### Front Matter

| Key               | Value   | Required | Description                                                           |
|-------------------|---------|----------|-----------------------------------------------------------------------|
| `confluence`      | integer | Yes      | ID of the page in Confluence.                                         |
| `title`           | string  | No       | The title of the page. If missing, first H1 or the file name is used. |
| `toc`             | boolean | No       | Prepend Table of Contents component. Defaults to `false`              |
| `ag-warning-pre`  | boolean | No       | Prepend warning that the page is autogenerated. Defaults to `true`    |
| `ag-warning-post` | boolean | No       | Append warning that the page is autogenerated. Defaults to `false`    |

#### Example

```markdown
---
confluence: 666
---
# This is my first auto-updated page

Hello there!
```

### Build

```bash
go build -o main
```

### Test

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) first.

## License

[MIT](LICENSE)
