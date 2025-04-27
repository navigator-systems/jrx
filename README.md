# JRX

![JRX Logo](images/logo.png)

**`jrx`** is a modern, Cargo-inspired command-line tool written in Golang, designed to streamline the development lifecycle of Go projects. From scaffolding to building and analyzing dependencies, `jrx` offers a developer-friendly interface with powerful functionality

## Features

- 🔧 `jrx new <project>` — Create a new project scaffold with a structured layout:
    ```bash
    myProject/
        ├── bin/
        ├── doc/
        ├── lib/
        ├── main.go
        ├── Makefile
        ├── jrx.toml
        └── Dockerfile
    ```
- 📦 `jrx mod <project>` — Equivalent to `go mod init`.

- 🛠️ `jrx build <project>` — Builds the project. If the module is not initialized, it initializes it first.
   - With flags as `arch` and `os` to build multiarchitecture binaries.

- 🔍 `jrx info <project>` — Displays:
   - Size of the binary file(s) in the `bin` folder.
   - Dependencies listed in `go.sum`.
   - Known vulnerabilities using the [OSV.dev](https://osv.dev) CVE database.

- 🛠️ `jrx ci --template <jenkins, github> <project>`  — Creates either a Jenkinsfile or a simple gituhub workflow yaml file for building the application. 
---
## jrx file
With every new project `jrx` will create a config file called `jrx.toml`, here you can write relevant information about the project. 
Standard jrx.toml file looks like this: 

```
name = "MyProject"
version = "0.0.1"
authors = [
    "My Name",
    
]

[builds.laptop] 
   arch = "386"
   os = "linux"

[builds.raspPi] 
   arch = "arm64"
   os = "linux"

[builds.release]
   arch = "amd64"
   os = "linux"
   flags = "-ldflags= -s -w"
    
```
The keys `name`, `version`, `authors` are descriptions of the project, users can update these values according to the needs of the project. Users can now 
define multiple custom build targets directly in jrx.toml, specifying architecture, operating system, and build flags under `[builds.<target>]` sections. 

This makes it easy to cross-compile for different environments like Linux, Windows, Raspberry Pi, or even older platforms — all from a single command `jrx build MyProject`.

The output binaries will be created in the bin directory with the following name sctructure `projectName-OS-ARCH`. If several build targets use the same OS and ARCH these are going to get overwritten. 


## 🛡️ Vulnerability Scanning

`jrx info --osv` uses the [OSV.dev API](https://osv.dev) to check for known vulnerabilities in Go dependencies. It automatically scans each dependency version listed in `go.sum` and outputs any relevant CVEs.

---

## 📁 Usage

```bash
NAME:
   jrx - Just a simple go wrapper CLI

USAGE:
   jrx [global options] command [command options]

COMMANDS:
   info, i   Get information from the project
   new, n    Create a new project
   build, b  Build and compile a project
   mod, m    Start a simple go.mod file
   ci        Add a CI template (Jenkins, GitHub Actions or Gitlab Template) to the project
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
