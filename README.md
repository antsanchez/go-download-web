# Go Download Web
A simple command-line application to download an entire online website, including CSS, JSS, Images, and other assets. 
Coded with Go.

## Project status
There are still some to-do's, and some refactoring is needed, but the app is already functional. 

To do: add headless browser support to download JS-generated content.
To do: resume downloads
To do: download attachments in parallel

## Installation
There is nothing special to do here, just download the code, build it as you would do with any other Go app, and you are set to go.

```bash
$ git clone https://github.com/antsanchez/go-download-web
$ cd go-download-web
$ go build
```

## Usage
```bash
# Default mode:
$ ./go-download-web -u https://example.com

# Setting the number of concurrent connections to 10:
$ ./go-download-web -u https://example.com -s 10

# Scrape a website without ignoring query params on URLs
$ ./go-download-web -u https://example.com -q

# Change the domain name of the downloaded copy:
$ ./go-download-web -u https://example.com -new https://newname.com
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Author
[Antonio Sánchez](https://asanchez.dev)

## License
[MIT](https://choosealicense.com/licenses/mit/)
