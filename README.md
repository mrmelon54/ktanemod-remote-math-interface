# Remote Math Interface

This is an interface to connect [Remote Math](https://github.com/mrmelon54/ktanemod-remote-math) to [Remote Math Server](https://github.com/mrmelon54/ktanemod-remote-math-server) using secure websockets.

Unfortunately as Unity doesn't support secure websockets in 2017.4.22f1 which is used for KTaNE modding. This interface is a hack which enables secure websockets.

The Go source code is compiled into DLL/SO/DYLIB binaries and is loaded by the mod in game.

## Build

Just use GitHub actions its easier

By default GitHub actions is configured to build the following binaries:

  - Windows: amd64, 386, arm64, arm
  - Linux: amd64, 386, arm64, arm
  - Mac: amd64, arm64
