## About

This project is a simple tool to check links of a website.

## How to use

- Git clone
- Go to the directory
- Execute the command "make build" to build the tool. PS: you don't need to build the tool if you want to run it using go run main.go
- Run the command below:
```bash

./checker-website-links -api-key <api-key> -link <link> -limit <limit> -disable-cache <disable-cache> -timeout <timeout> -output <output> -max-time-ms-accepted <max-time-ms-accepted>

```

## Options required

- api-key: API key of the website to check links.
- link: Link of the website to check links.

## Options

- api-key: API key of Firecrawl(https://www.firecrawl.dev/).
- link: Link of the website to check links. For example: https://www.firecrawl.dev/ or https://www.abacatepay.com/ or https://www.google.com/
- limit: Limit of links to check. Default: 100
- disable-cache: Disable cache of the website to check links. Default: false . If true, the website will not use cache to check links, because will add "?v=unix_time_milliseconds" to the link to avoid cache.
- timeout: Timeout of the website to check links. Default: 5
- output: Output file of the website to check links. Default: output.json . The file will be saved in the same directory of the tool and contains the links checked.
- max-time-ms-accepted: Max time accepted of the website to check links. Default: 5000

## Examples:

- Check Github actions workflow on folder .github/workflows to see how to use.
