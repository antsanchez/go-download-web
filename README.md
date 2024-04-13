# Go Download Web

Go Download Web is a command-line application developed with Go that allows you to download an entire online website, including CSS, JavaScript, Images, and other assets.

## Project Status
The application is functional and has been improved recently. However, there are still some tasks pending:

- Add headless browser support to download JS-generated content.
- Implement download resuming.
- Enable parallel downloading of attachments.

## Installation
To install, follow these simple steps:

```bash
$ git clone https://github.com/antsanchez/go-download-web
$ cd go-download-web
$ go build
```

## Usage
The application can be used with various flags to customize the download process:

```bash
$ ./go-download-web -u <URL>
```
- `-u` or `--url`: The URL of the website to download content from. This is a required field.

```bash
$ ./go-download-web -u <URL> -new <NEW_URL>
```
- `-new` or `--new-url`: The new URL to use for the downloaded content. This is an optional field.

```bash
$ ./go-download-web -u <URL> -r <INCLUDED_URLS>
```
- `-r` or `--included-urls`: The URL prefixes/root paths that should be included in the download. This is an optional field.

```bash
$ ./go-download-web -u <URL> -s <SIMULTANEOUS_CONNECTIONS>
```
- `-s` or `--simultaneous`: The number of concurrent connections. The default value is 3, and the minimum is 1.

```bash
$ ./go-download-web -u <URL> -q
```
- `-q` or `--use-queries`: A flag to ignore query strings in URLs. This is an optional field.

```bash
$ ./go-download-web -u <URL> -path <DOWNLOAD_PATH>
```
- `-path` or `--download-path`: The local path to save the downloaded files. The default value is `./website`.

For help, use the `-h` or `--help` flag:

```bash
$ ./go-download-web -h
```

## Development

### Generate Mocks

This project uses [uber-go/mock](https://github.com/uber-go/mock) to generate mocks for testing. To generate mocks, run the following commands:

For the console:
```bash
mockgen -destination=pkg/console/mock_console.go -package=console github.com/antsanchez/go-download-web/pkg/scraper Console
```

For the HTTP client:
```bash
mockgen -destination=pkg/get/mock_get.go -package=get github.com/antsanchez/go-download-web/pkg/scraper HttpGet
```


These commands generate mocks for the `Console` and `HttpGet` interfaces in the `scraper` package. The generated mocks are saved in the `pkg/console` and `pkg/get` packages, respectively.

### Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

## Author
[Antonio Sánchez](https://asanchez.dev)

## License
[MIT](https://choosealicense.com/licenses/mit/)
