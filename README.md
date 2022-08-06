# vzip
📁 A simple, fast and lightweight file and folder compression utility

## 🏃 Usage
`vzip <file name | folder name> [--level=int] [--method=string]`

## ❓ Options
`--level=int`: What compression level should be used (0-9). The higher the number, the slower the operation, but the smaller the file.

> **Note** 0 disables compression

`--method=string`: What compression method should be used.

Accepted values:
  - `none`
  - `gzip`
  - `zlib`
