# RVC-Models-Downloader

[English](README.md) | [简体中文](README_sc.md) | 한국어

yaml의 batch 파일을 쉽게 다운로드 할 수 있는 도구입니다. (Hugging Face 🤗의 RVC 모델 등).

![tui demo](https://github.com/fumiama/RVC-Models-Downloader/assets/41315874/db577dfb-8a6d-4909-b071-9d36cc77afc6)

## 빠른 시작

### 준비

1. [릴리스](https://github.com/fumiama/RVC-Models-Downloader/releases) 페이지에서 프로그램을 다운로드를 받아주세요.
2. 해당 프로그램을 RVC의 루트 디렉토리(또는 파일을 다운로드하고 싶은 위치)에 넣어주세요.
3. 이 도구를 어디에서나 사용할 수 있도록 `PATH`에 추가할 수도 있습니다. 패키지 매니저를 통해 이 프로그램을 설치했다면 이미 `PATH`에 등록되어 있을 수 있습니다.

### 다운로드

#### RVC의 모든 자산

```bash
rvcmd assets/rvc
```

#### ChatTTS의 모든 자산

```bash
rvcmd -w 1 assets/chtts
```

### 사용자 정의 다운로드

#### 예시 1. hubert & rmvpe 다운로드

1. 다음 내용을 포함한 `cust.yaml`을 작성하고 저장합니다.
    ```yaml
    BaseURL: https://huggingface.co/fumiama/RVC-Pretrained-Models/resolve/main
    Targets:
      - Refer: hubert
      - Refer: rmvpe
    ```
2. 같은 폴더에서 `rvcmd`를 실행합니다.
    ```bash
    rvcmd -c cust
    ```

#### 예시 2. 🤗의 다른 저장소 다운로드

> [Stable Diffusion v1-5](https://huggingface.co/runwayml/stable-diffusion-v1-5)를 예시로 사용합니다.

1. 다음 내용을 포함한 `cust.yaml`을 작성하고 저장합니다.
    ```yaml
    BaseURL: https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main
    Targets:
      - Folder: sd1.5 # 다운로드할 폴더
        Copy: # 다운로드할 파일
          - v1-5-pruned-emaonly.ckpt
          - v1-5-pruned-emaonly.safetensors
      - Folder: sd1.5/vae # 다운로드할 폴더
        Copy: # 다운로드할 파일
          - vae/diffusion_pytorch_model.bin
    ```

#### 예시 3. GitHub에서 릴리스 다운로드

> [yousa-ling-diffsinger-v1.3](https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/tag/v1.3)를 예시로 사용합니다.

1. 다음 내용을 포함한 `cust.yaml`을 작성하고 저장합니다.
    ```yaml
    BaseURL: https://github.com/yousa-ling-official-production/yousa-ling-diffsinger-v1/releases/download/v1.3
    Targets:
      - Folder: . # 다운로드할 폴더
        Copy: # 다운로드할 파일
          - yousaV1.3.zip
    ```
2. 같은 폴더에서 `rvcmd`를 실행합니다.
    ```bash
    rvcmd -c cust
    ```

## 전체 사용법

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

## 데모 비디오

https://github.com/fumiama/RVC-Models-Downloader/assets/41315874/da2b5827-8b1a-45f8-a9c0-03a5618ad5f8
