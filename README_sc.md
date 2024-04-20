# RVC模型下载器
[English](README.md) | 简体中文

一个能够批量下载`yaml`清单内文件的简单工具（例如 Hugging Face 🤗 中的 RVC 模型）。

![tui demo](https://github.com/RVC-Project/RVC-Models-Downloader/assets/41315874/faec35ea-b7af-4404-83f3-ecca73da9abc)

## 快速开始
### 准备工作
1. 在[发布](https://github.com/RVC-Project/RVC-Models-Downloader/releases)页面下载程序。
2. 将此程序放入RVC的根目录（或您想要下载文件的任何位置）。
3. 您也可以将它添加到`PATH`中以便在任何地方使用此工具。如果您已经通过包管理器安装了此程序，那么它可能已经位于`PATH`。
### 下载
#### RVC的所有资源文件
```bash
rvcmd assets/all
```
#### RVC的最新通用整合包（仅限Windows）
```bash
rvcmd packs/general/latest
```
#### ffmpeg工具（仅限Windows）
```bash
rvcmd tools/ffmpeg
```
### 自定义下载
#### 示例1. 下载ffmpeg工具和最新的Intel包
1. 编写并保存以下`cust.yaml`。
    ```yaml
    BaseURL: https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main
    Targets:
      - Refer: tools/ffmpeg
      - Refer: packs/intel/latest
    ```
2. 在同一文件夹中运行`rvcmd`。
    ```bash
    rvcmd -c cust
    ```
#### 示例2. 下载🤗中的其他仓库
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
#### 示例3. 下载GitHub中的发布版本
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
  -w uint
        connection waiting seconds (default 4)
  'target/to/download'
        like packs/general/latest
All available targets:
    assets:
        all    hubert    rmvpe    uvr5    v1    v2
    packs:
        amd:
            latest
            v2:
                20230813    20231006
        general:
            latest
            v1:
                20230331    20230416    20230428    20230508    20230513    20230516    20230717
            v2:
                20230528    20230618
        intel:
            latest
            v2:
                20230813    20231006
        nvidia:
            latest
            v2:
                20230813    20231006
    tools:
        ffmpeg
```
