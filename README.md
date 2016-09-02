## Overview
A single executable to do all the useful Ops commands.

## Why?
- No runtime dependencies. Much easier to distribute to any user without setting up ruby environment and building gems.
- Much much faster. This is important for quick tasks. For instance, `knife data bag show` took 6s while `pops` can do it in under 0.1s.
- Single command to do everything, just like `git`.
- Self-documented with `--help`.
- Open source, including all the dependencies. No need to use your ssh key to pull the source or use the distribution.

## Installation
Check [releases](https://github.com/whence/pops/releases) for downloads. Put the executable into your $PATH. That's it.

## Current available commands
- `dec` Decrypt Chef data bag
- `enc` Encrypt Chef data bag
- `db up` Start the database
- `db init` Init the database
- `db down` Destroy the database
- `rand iv` Generate an random initialization vector
- `rand secret` Generate an random secret (Can be used to encrypt/decrypt Chef data bag)

You can get the usage of all commands by `pops [command] -h`
