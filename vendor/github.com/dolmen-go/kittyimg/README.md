# kittyimg

`kittyimg` is a Go library that allows to display images in terminal emulators implementing [kitty's *terminal graphics protocol*](https://sw.kovidgoyal.net/kitty/graphics-protocol.html).

[![Go Reference](https://pkg.go.dev/badge/github.com/dolmen-go/kittyimg.svg)](https://pkg.go.dev/github.com/dolmen-go/kittyimg)
[![CI](https://github.com/dolmen-go/kittyimg/actions/workflows/test.yml/badge.svg)](https://github.com/dolmen-go/kittyimg/actions)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=dolmen-go_kittyimg&metric=coverage)](https://sonarcloud.io/summary/overall?id=dolmen-go_kittyimg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolmen-go/kittyimg)](https://goreportcard.com/report/github.com/dolmen-go/kittyimg)


<!--
[![Codecov](https://img.shields.io/codecov/c/github/dolmen-go/kittyimg/main.svg)](https://codecov.io/gh/dolmen-go/kittyimg/branch/main)
[![Codecov](https://codecov.io/gh/dolmen-go/kittyimg/graph/badge.svg?token=YpWZ0OBMUa)](https://app.codecov.io/gh/dolmen-go/kittyimg/branch/main)
-->

## ✨ Features

A [basic API](https://pkg.go.dev/github.com/dolmen-go/kittyimg)
([`Fprint`](https://pkg.go.dev/github.com/dolmen-go/kittyimg#Fprint),
[`Fprintln`](https://pkg.go.dev/github.com/dolmen-go/kittyimg#Fprintln),
[`Transcode`](https://pkg.go.dev/github.com/dolmen-go/kittyimg#Transcode)) allows to display an image (loaded with stdlib's [image](https://pkg.go.dev/image) package) at the cursor position.

```console
go get github.com/dolmen-go/kittyimg@latest
```

A command-line tool ([`icat`](https://pkg.go.dev/github.com/dolmen-go/kittyimg/cmd/icat)) is provided.

```console
go install github.com/dolmen-go/kittyimg/cmd/icat@latest
```

`icat <image>` works the same as [Kitty's command](https://sw.kovidgoyal.net/kitty/kittens/icat/) `kitten icat --transfer-mode=stream --align=left <image>`.

## 🏗️ Status

Production ready.

## 🔄 See also

The [Go Playground](https://go.dev/play) has [support for displaying images](https://play.golang.org/p/LXmxkAV0z_M) with its own protocol: `IMAGE:` prefix followed by base64 image file data.

Display tools for images on terminals:
* [tycat](https://git.enlightenment.org/apps/terminology.git/tree/src/bin/tycat.c):
Similar tool to `icat`, but for the Enlightenment Terminology app (which uses a different terminal protocol).
* [timg](https://github.com/hzeller/timg)
* [viu](https://github.com/atanunq/viu)
* [chafa](https://hpjansson.org/chafa/)

## 🛡️ License

Copyright 2021-2026 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
