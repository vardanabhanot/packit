# 🧰 Packit

**Packit** is a lightweight CLI tool for creating compressed archives — currently supporting **ZIP** files.
It was originally built for a specific use case, so the code is intentionally simple and verbose.
Support for additional formats like **TAR** and others is planned in future releases.

> It generates a .packit file in that directory to handle ignores.

---

## 🚀 Features

- Create ZIP archives quickly from the command line
- Ignore specific files or directories
- Lightweight — no external dependencies
- Simple, minimal, and easy to extend

---

## ⚙️ Usage
```bash
Usage: packit <command> [flags] [arguments]

Available Commands:

 packit build
        -f      (Default: zip)  the format of archive you want to create zip/tar
        -o      (Default: )     name of the output file relative to current directory without any extension

 packit ignore
        -l      (Default: false)        Lists all the files which are in the ignore list

        Adding files to ignore
        packit ignore file1 file2 file3
```


## 📦 Installation

### Using Go
```bash
go install github.com/vardanabhanot/packit
```
