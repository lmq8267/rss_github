# rss_github

[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Flmq8267%2Frss_github&count_bg=%2395C10D&title_bg=%23555555&icon=github.svg&icon_color=%238DC409&title=%E8%AE%BF%E9%97%AE%E6%95%B0&edge_flat=false)](https://hits.seeyoufarm.com)

#### 解析github仓库的atom
```shell
  -u string  指定 GitHub 仓库链接（必填项）
  -all       同时解析 releases.atom 和 commits.atom
  -r	       解析版本记录 releases.atom
  -c	       解析提交记录 commits.atom
  -n int     指定检查的数量 (default 1)
  -o string  指定输出文件的目录路径（文件夹）
```
```shell
#效果
lmq8267@ubuntu:go build 
lmq8267@ubuntu:./rss -u https://github.com/lmq8267/rss_github -c
=== 提交记录信息 ===
提交说明：Create go.sum
提交时间：2024-10-23 11:16:22
详细链接：https://github.com/lmq8267/rss_github/commit/2dd7a984efbeca0b7927a174553c0db071e824e9
```

###### 编译：
```shell
#本机编译
CGO_ENABLED=0 go build -ldflags="-s -w"
```
```shell
#交叉编译aarch64
CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w"
```
```shell
#交叉编译mispel
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w"
```

