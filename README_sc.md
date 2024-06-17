# RVC 模型下载器

[English](README.md) | 简体中文 | [한국어](README_kr.md)

一个能够批量下载`yaml`清单内文件的简单工具（例如 Hugging Face 🤗 中的 RVC 模型）。

![tui demo](https://github.com/fumiama/RVC-Models-Downloader/assets/41315874/db577dfb-8a6d-4909-b071-9d36cc77afc6)

## 快速开始

### 准备工作

1. 在[发布](https://github.com/fumiama/RVC-Models-Downloader/releases)页面下载程序。
2. 将此程序放入 RVC 的根目录（或您想要下载文件的任何位置）。
3. 您也可以将它添加到`PATH`中以便在任何地方使用此工具。如果您已经通过包管理器安装了此程序，那么它可能已经位于`PATH`。

### 下载

#### RVC 的所有资源文件

```bash
rvcmd assets/all
```

#### ChatTTS 的所有资源文件
```bash
rvcmd -w 1 assets/chtts
```

### 自定义下载

#### 示例 1. 下载 hubert 和 rmvpe

1. 编写并保存以下`cust.yaml`。
    ```yaml
    BaseURL: https://huggingface.co/fumiama/RVC-Pretrained-Models/resolve/main
    Targets:
      - Refer: hubert
      - Refer: rmvpe
    ```
2. 在同一文件夹中运行`rvcmd`。
    ```bash
    rvcmd -c cust
    ```

#### 示例 2. 下载 🤗 中的其他仓库

> 以 [Stable Diffusion v1-5](https://huggingface.co/runwayml/stable-diffusion-v1-5) 为例。

1. 编写并保存以下`cust.yaml`。
    ```yaml
    BaseURL: https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main
    Targets:
      - Folder: sd1.5 # the folder you want to download into
        Copy: # files to download
          - v1-5-pruned-emaonly.ckpt
          - v1-5-pruned-emaonly.safetensors
      - Folder: sd1.5/vae # the folder you want to download into
        Copy: # files to download
          - vae/diffusion_pytorch_model.bin
    ```

#### 示例 3. 下载 GitHub 中的发布版本

> 以 [yousa-ling-diffsinger-v1.3](https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/tag/v1.3) 为例。

1. 编写并保存以下`cust.yaml`。
    ```yaml
    BaseURL: https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/download/v1.3
    Targets:
      - Folder: . # the folder you want to download into
        Copy: # files to download
          - yousaV1.3.zip
    ```
2. 在同一文件夹中运行`rvcmd`。
    ```bash
    rvcmd -c cust
    ```

## 完整用法

```bash
Usage: rvcmd [-notrs] [-dns dns.yaml] 'target/to/download'
  -c    use custom yaml instruction
  -dns string
        custom dns.yaml
  -f    force download even file exists
  -notrs
        use standard TLS client
  -notui
        use plain text instead of TUI
  -ua string
        customize user agent (default "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
  -w uint
        connection waiting seconds (default 4)
  'target/to/download'
        like packs/general/latest
All available targets:
    assets:
        chtts    hubert    rmvpe    rvc    uvr5    v1    v2
```

## 示例录屏

https://github.com/fumiama/RVC-Models-Downloader/assets/41315874/da2b5827-8b1a-45f8-a9c0-03a5618ad5f8
