# Confluence Pages Creator GitHub Action

This GitHub Action scans your repository for OpenAPI definition files and Markdown documentation to automatically create or update Confluence pages.

![GitHub Action CI](https://github.com/danaketh/ga-confluence-gen/workflows/CI/badge.svg)

## Features

- Scans your repository for OpenAPI `.yaml` or `.json` files.
- Scans your repository for Markdown documentation files.
- Automatically creates or updates Confluence pages with the found documentation.

## Requirements

- A running instance of Confluence that is accessible from GitHub Actions.
- A Confluence account with permissions to create and update pages.

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
        confluence-url: ${{ vars.CONFLUENCE_URL }}
        confluence-username: ${{ vars.CONFLUENCE_USERNAME }}
        confluence-token: ${{ secrets.CONFLUENCE_TOKEN }}
```

## Configuration

### Environment Variables

- `CONFLUENCE_URL`: The URL of your Confluence instance.
- `CONFLUENCE_USERNAME`: The username of the Confluence account.
- `CONFLUENCE_TOKEN`: API token for the Confluence account.

### Action Inputs

| Input                 | Description             | Required | Default |
|-----------------------|-------------------------|----------|---------|
| `confluence-url`      | Confluence Instance URL | `true`   |         |
| `confluence-username` | Confluence Username     | `true`   |         |
| `confluence-token`    | Confluence API Token    | `true`   |         |

## How it Works

1. The action first scans the repository for any OpenAPI `.yaml` or `.json` files.
2. It then scans for Markdown files.
3. The action then interacts with the Confluence API to create or update pages based on the found files.

## Development

Built with Go.

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

[MIT](LICENSE.md)
