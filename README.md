下载器小玩具
==================

[![Build Status](https://github.com/zion-chuu/fdlr/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/zion-chuu/fdlr/actions)

利用 MapReduce 的思路, 将大块数据拆分成小块数据处理, 然后最后再合并, 输出结果, 所以对于很小的下载数据来说, 效果并不好 
下载功能利用 HTTP 头的 Range 功能

## 功能:
* 支持 HTTP 与 HTTPS 应用层协议(其实是没特殊处理 HTTPS, 所以这两种效果都一样)
* 可以设置具体将要下载的数据进行分段数量, 每段一个 Goroutine
* 分段之后的数据, 每段并行处理
* 如果中途比如说, ctrl-D 了, 之后可以继续前面的下载(前提是之前的下载位置没变)

## 构建要求:
Go v1.6+

## 栗子:
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

## 灵感来自:
- [https://github.com/ytdl-org/youtube-dl](https://github.com/ytdl-org/youtube-dl)
- [https://github.com/cavaliercoder/grab](https://github.com/cavaliercoder/grab)
