Download file in Golang
==================

[![Build Status](https://github.com/Imputes/fdlr/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/Imputes/fdlr/actions)

## Featues:
* support HTTP and HTTPS
* you can set the number of parallel to download
* download batches of files concurrently
* resume incomplete downloads

## Requires
Go v1.6+

## Example
```
$ fdlr download -c=3 https://download.jetbrains.com/go/goland-2020.2.2.dmg
Downloading IP is: 52.50.241.213 | 54.72.98.183
Start downloading with 3 connections 
Download target size: 398.9 MB
goland-2020.2.2.dmg - 0 1.73 MiB / 132.97 MiB    1.30% 04m41s                                          
goland-2020.2.2.dmg - 1 732.32 KiB / 132.97 MiB    0.54% 11m27s                                        
goland-2020.2.2.dmg - 2 1005.57 KiB / 132.97 MiB    0.74% 08m19s                                       
Interrupted, saving state ... 
Saving states data in /Users/xxx/.fdlr/goland-2020.2.2.dmg
```

```
$ fdlr resume https://download.jetbrains.com/go/goland-2020.2.2.dmg
Reading state from /Users/xxx/.fdlr/goland-2020.2.2.dmg/state.yaml
goland-2020.2.2.dmg - 0 510.94 KiB / 131.23 MiB    0.38% 19m58s                                        
goland-2020.2.2.dmg - 1 489.65 KiB / 131.98 MiB    0.36% 20m57s                                        
goland-2020.2.2.dmg - 2 1.20 MiB / 132.25 MiB    0.91% 08m19s
```

## Inspired
- [https://github.com/ytdl-org/youtube-dl](https://github.com/ytdl-org/youtube-dl)
- [https://github.com/cavaliercoder/grab](https://github.com/cavaliercoder/grab)
