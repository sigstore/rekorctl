# rekor-ctl

`rekorctl` is a CLI tool to interact with a rekor server. It provides a set of commands that are
useful to anyone wanting to interact with a rekor instance, beyond just making and verifying entries.

> :warning: **If you are a developer** and just want to make use of rekor within your project, please use [rekor-cli](https://github.com/sigstore/rekor/)

```
%  rekorctl --help
  get         Gets an entry using the artefact
  getleaf     Get an entry via a Leaf Index
  sigs-by-artifact Look up all signature entries for an artifact and return them.
  sigs-by-pub Look up signatures by public key
  help        Help about any command
  update      Rekor update command
```

## Security

Should you discover any security issues, please refer to sigstores [security
process](https://github.com/sigstore/community/blob/main/SECURITY.md)
