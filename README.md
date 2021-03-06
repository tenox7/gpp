Go Ping Plot
============

![screenshot](gpp.png)

* Multi target, ping multiple hosts at once
* Line color changes with latency, red on timeout
* Multi OS, works on Linux, FreeBSD, macOS and Windows
* Full screen mode for devices like Raspberry PI
* Portable, self contained binaries
* Written in GoLang, uses SDL2

Usage
-----
```
gpp [flags] host [host [...]]

flags:
    -t  - threshold for "bad ping" - line color will be more reddish
    -fs - full screen mode (eg. for Raspberry PI)
    -fg - border and text color in hex rgb
    -bg - background color in hex rgb
```

Downloads
---------
Under [Releases](https://github.com/tenox7/gpp/releases)

Copyright
---------
Copyright 2019 Google LLC

Licenses
--------
* Go Ping Plot is Apache 2.0
* Uses SDL2 which is BSD-3-Clause
* Uses go-ping which is MIT
* Includes Google Noto Font which is SIL Open Font License, Version 1.1.
* Uses rakyll/statik which is Apache 2.0
