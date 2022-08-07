# vzip
📁 A simple, fast and lightweight file and folder compression utility

## 🏃 Usage
```sh-session
$ vzip <file name | folder name> [--dest=string] [--level=int] [--method=string]
```

## ❓ Options
`--dest=string`: Where the compressed file will be saved. Defaults to `./<file name>.zip`

`--level=int`: What compression level should be used (0-9). The higher the number, the slower the operation, but the smaller the file.

> **Note** 0 disables compression

`--method=string`: What compression method should be used.

Accepted values:
  - `none`
  - `gzip`
  - `zlib`

## 💻 Building from source
Clone the repository:

```sh-session
$ git clone https://github.com/davipatricio/vzip.git
```

Build the project:

```sh-session
$ cd vzip
$ go build vzip.go
```

If you wish a smaller binary size, compile the program using `go build -ldflags="-w -s" -gcflags=all="-l -B" vzip.go`.

## 📝 License
This project is licensed under the [MIT](LICENSE) license.
