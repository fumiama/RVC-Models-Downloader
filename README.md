# RVC-Models-Downloader
Quickly download RVC models in 🤗 Hugging Face.

## Quick Start
### Preparation
Put this program into the root directory of RVC.
### Download All Assets
```bash
rvcmd assets/all
```
### Download Latest General Pack (Windows Only)
```bash
rvcmd packs/general/latest
```
### Download ffmpeg Tools (Windows Only)
```bash
rvcmd tools/ffmpeg
```

## Full Usage
```bash
Usage: rvcmd [-notrs] [-dns dns.yaml] 'target/to/download'
  -dns string
        custom dns.yaml
  -notrs
        use standard TLS client
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
