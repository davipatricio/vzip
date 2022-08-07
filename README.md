# vzip
📁 A simple, fast and lightweight file and folder compression utility

## 🏃 Usage
```sh-session
$ vzip <file name | folder name> [--dest=string] [--level=int] [--method=string]
```

## ❓ Options
`--dest=int`: Where the compressed file will be saved. Defaults to `./<file name>.zip`

`--level=int`: What compression level should be used (0-9). The higher the number, the slower the operation, but the smaller the file.

> **Note** 0 disables compression

`--method=string`: What compression method should be used.

Accepted values:
  - `none`
  - `gzip`
  - `zlib`
