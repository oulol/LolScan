
# LolScan - a pure Go port scanner
![Logo](./images/logo.png)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
![GitHub forks](https://img.shields.io/github/forks/oulol/LolScan)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/oulol/LolScan)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/oulol/LolScan)

**LolScan** is a fast and multi-threaded port scanner written entirely in Go. It has a built-in bruteforcer to automatically test for weak credentials.

## Warning
> [!WARNING]
This tool is provided for educational and research purposes only. 
The author assumes no liability for any misuse, damage, or illegal activity caused by this software. 
You are solely responsible for complying with all applicable local, national, and international laws. 
By using this tool, you agree to these terms. Use at your own risk.

## Installation

### Download Pre-Compiled Binaries
To install LolScan quickly, just download pre-compiled binary for your OS and architecture in Releases.

### Building From Source
To build it manually you will need **Go 1.26.4 or higher** and **Git**.

#### Linux / macOS / Git Bash
```sh
# Clone the repository
git clone https://github.com/oulol/LolScan.git
cd LolScan

# Build the executable
go build -ldflags="-s -w -X main.Version=$(git describe --tags)" -trimpath -o LolScan
```

#### Windows (PowerShell)
```powershell
# Clone the repository
git clone https://github.com/oulol/LolScan.git
cd LolScan


# Build the executable
go build -ldflags="-s -w -X main.Version=\$(git describe --tags)" -trimpath -o LolScan.exe
```

You will end up with an executable file named "LolScan". 
## Usages

To get a list of all available arguments, run:
```sh
./LolScan -help
```

Example scanning without bruteforce for port 80:
```sh
# Adding a subnet to the ips.txt file (can be done any other way)
echo 192.168.0.1/24 > ips.txt

./LolScan -ips ips.txt -ports 80 -threads 120 -nobrute
```
## Roadmap

- Add more services support
- Make scanning faster
- Add snapshot grabbing from supported services
- Add SYN scanning that works on Windows and Linux

## Contribution

Feel free to contribute to this project!
For more details see file [CONTRIBUTING](CONTRIBUTING.md) 
## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
