# grpcp

## Description

This is a simple tool to copy files between local and remote hosts using gRPC stream.

## Usage

```
Usage: grpcp [<src> [<dest>]]

Arguments:
  [<src>]
  [<dest>]

Flags:
  -h, --help         Show context-sensitive help.
  -p, --port=8022    port number
  -l, --listen=""    listen address
  -s, --server       run as server
  -q, --quiet        quiet mode for client
```

Start the server on the remote host:
```console
$ grpcp --server
```

### Examples

Copy a file from the remote host to the local host:
```console
$ grpcp remote_host:/path/to/file /path/to/destination
```

Copy a file from the local host to the remote host:
```console
$ grpcp /path/to/file remote_host:/path/to/destination
```

## TODO

- Use TLS encryption for the connection

## LICENSE

MIT
