# pops

## Overview
A progressive rewrite of [platform_ops](https://github.com/MYOB-Technology/platform_ops) and some other useful ops utilities in Go.

## Why creating another repo?
- No runtime dependencies. Much easier to distribute to any user without setting up ruby environment and building gems.
- Much much faster. This is important for quick tasks. For instance, `knife data bag show` took 6s while `pops` can do it in under 0.1s.
- Single command to do everything, just like `git`.
- Self-documented with `--help`.
- Open source, including all the dependencies. No need to use your ssh key to pull the source or use the distribution.
