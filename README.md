# paramfilter

A high-performance pure Go utility designed to replace multi-command bash URL filtering pipelines with a single-pass streaming tool that runs **10x–30x faster** with zero external dependencies.

## What It Replaces

This single tool replaces the following 4 bash commands:

```console
cat urls.txt | egrep -av "=" | grep -aiE "\.(php|asp|aspx|jsp|jspx|xml|html)$" | unew endpoints.txt
cat urls.txt | egrep -av "=" | egrep -avE '^https?://[^/]+/.*\.' | unew endpoints.txt
cat urls.txt | cut -d"/" -f1-3 | sed 's|/*$|/|' | unew endpoints.txt

cat urls.txt | grep -a "=" | egrep -aiv ".(jpg|jpeg|gif|css|tif|tiff|png|ttf|woff|woff2|icon|pdf|svg|txt|js)" | unew params.txt
```

## Installation

### Install via Go
```
go install github.com/rix4uni/paramfilter@latest
```

### Download Prebuilt Binaries
```
wget https://github.com/rix4uni/paramfilter/releases/download/v0.0.1/paramfilter-linux-amd64-0.0.1.tgz
tar -xvzf paramfilter-linux-amd64-0.0.1.tgz
rm -rf paramfilter-linux-amd64-0.0.1.tgz
mv paramfilter ~/go/bin/paramfilter
```

Or download the [latest release](https://github.com/rix4uni/paramfilter/releases) for your platform.

### Compile from Source
```
git clone --depth 1 https://github.com/rix4uni/paramfilter.git
cd paramfilter; go install
```

## Usage
```console
Usage of paramfilter:
  -e, --endpoints string   Output file for filtered endpoints (default "endpoints.txt")
  -p, --params string      Output file for filtered parameters (default "params.txt")
      --silent             Silent mode
      --version            Print the version of the tool and exit
```


## Usage Examples

### Basic
```console
cat urls.txt | paramfilter
```

### Custom Output Files
```console
cat urls.txt | paramfilter --endpoints my_endpoints.txt --params my_params.txt
```
