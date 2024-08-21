# grpcp

## Description

This is a simple tool to copy files between local and remote hosts using gRPC stream.

## Usage

```
Usage: grpcp [<src> [<dest>]] [flags]

Arguments:
  [<src>]
  [<dest>]

Flags:
  -h, --help                Show context-sensitive help.
  -h, --host="localhost"    host name
  -p, --port=8022           port number
  -q, --quiet               quiet mode
  -d, --debug               enable debug log
      --[no-]tls            enable TLS (default: true)
  -s, --server              run as server
      --cert=STRING         certificate file for server
      --key=STRING          private key file for server
      --verify              TLS verification for client
      --kill                send shutdown command to server
      --ping                send ping message to server
```

Start the server on the remote host:
```console
$ grpcp --server
```

Copy a file from the remote host to the local host:
```console
$ grpcp remote_host:/path/to/file /path/to/destination
```

Copy a file from the local host to the remote host:
```console
$ grpcp /path/to/file remote_host:/path/to/destination
```

grpcp does not support copying directories, local to local, or remote to remote.

### TLS Configuration

grpcp enables TLS with self-signed certificate by default. If you want to use your own certificate, you can specify the certificate and private key files:
```console
$ grpcp --server --cert server.crt --key server.key
```

grpcp client does not verify the server certificate by default. If you want to verify the server certificate, you can specify the `--verify` flag:
```console
$ grpcp --verify remote_host:/path/to/file /path/to/destination
```

## LICENSE

MIT
