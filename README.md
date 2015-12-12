# codic

This project is command line tool for [codic](https://codic.jp/my/api_status).

## Usage

### 1. Install

To install, use `go get`:

```bash
$ go get -d github.com/ladicle/codic
```

### 2. Setup AccessToken

Can not open config file.

Login to codic.
https://codic.jp/login

Get the AccessToken in the API page.
And, save it:

```bash
$ echo 'YOUR_ACCESS_TOKEN' > ~/.codic
```

## Contribution

1. Fork ([https://github.com/ladicle/codic/fork](https://github.com/ladicle/codic/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[ladicle](https://github.com/ladicle)
