# RVC-Models-Downloader
Quickly download RVC models in 🤗 Hugging Face.

## Quick Start
### Preparation
Put this program into the root directory of RVC. You can download it at [Release](https://github.com/RVC-Project/RVC-Models-Downloader/releases) page.
### Download
#### All Assets
```bash
rvcmd assets/all
```
#### Latest General Pack (Windows Only)
```bash
rvcmd packs/general/latest
```
#### ffmpeg Tools (Windows Only)
```bash
rvcmd tools/ffmpeg
```
### Customized Download
#### Ex.1. Download ffmpeg Tools & Latest Intel Pack
1. Write and save the following `cust.yaml`.
    ```yaml
    BaseURL: https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main
    Targets:
      - Refer: tools/ffmpeg
      - Refer: packs/intel/latest
    ```
2. Run `rvcmd` in the same folder.
    ```bash
    rvcmd -c cust
    ```
#### Ex.2. Download other Repositories in 🤗
> Use [Stable Diffusion v1-5](https://huggingface.co/runwayml/stable-diffusion-v1-5) as the example.
1. Write and save the following `cust.yaml`.
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
#### Ex.3. Download Releases in GitHub
> Use [yousa-ling-diffsinger-v1.3](https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/tag/v1.3) as the example.
1. Write and save the following `cust.yaml`.
    ```yaml
    BaseURL: https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/download/v1.3
    Targets:
      - Folder: . # the folder you want to download into
        Copy: # files to download
          - yousaV1.3.zip
    ```
2. Run `rvcmd` in the same folder.
    ```bash
    rvcmd -c cust
    ```
## Full Usage
```bash
Usage: rvcmd [-notrs] [-dns dns.yaml] 'target/to/download'
  -c    use custom yaml instruction
  -dns string
        custom dns.yaml
  -f    force download even file exists
  -notrs
        use standard TLS client
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
