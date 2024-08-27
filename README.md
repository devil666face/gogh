# GOGH

> Just share file via github api
> All files crypted via AES256

## Installation

1. Download release via [release page](https://github.com/Devil666face/gogh/releases/latest)
2. Unzip `tar -xf gogh_linux_amd64.tar.gz`
3. Add in your path `sudo cp gogh /usr/local/bin` or use via `./gogh`

## Usage

#### Settings

1. Set your github session cookie token

- Open your browser
- Login on github and open github page
- F12
- Store -> Cookies -> user_session
- Copy value like `bzzsDpsejMS8jiKH-eezHwUlQWERTYvPix0qGPxg8YyL-b-7`

2. Set token on gogh

- `settings`
- `set token`
- enter your token `bzzsDpsejMS8jiKH-eezHwUlQWERTYvPix0qGPxg8YyL-b-7`

```
❯ gogh
🪣 >> settings
🪣 >> ⚙️ settings >> set token
```

3. Check what token sets

- `show`

```
🪣 >> ⚙️ settings >> show
 Key      Value
 compress +
 token    K9j-9lgbr2yjuodm5mQWEQWEQx3pmD_EkJa5kq-qXY52y
```

#### Store

##### Upload file

- `store`
- `upload`
- choose file

```
❯ gogh
🪣 >> store
🪣 >> 🗄️ store >> upload
👌 uploaded test.txt
```

##### Download file

- `store`
- `download`
- choose file

```
❯ gogh
🪣 >> store
🪣 >> 🗄️ store >> download
✅ downloaded test.txt
```

##### Show files

- `store`
- `show`

```
❯ gogh
🪣 >> store
🪣 >> 🗄️ store >> show
 # Name                        Created      Size     Compress
 1 test.txt                    Aug 27 16:04 4 B      +
```

#### Share

> You can share gogh files with all links and AES keys

##### Export

- `share`
- `export`
- choose uploaded remote file
- exported file save on `file_name.gogh`
- send `file_name.gogh` to you friend

```
❯ gogh
🪣 >> share
🪣 >> 🔁 share >> export
💾 file test.txt saved
```

##### Import

- `share`
- `import`
- choose local file in current dir with `.gogh`
- now you can download file via store

```
❯ gogh
🪣 >> share
🪣 >> 🔁 share >> import
📥 file test.txt imported
```

#### Configuration

1. Gogh save database in `~/.config/gogh/gogh.enc`
2. Use [config](https://github.com/Devil666face/gogh/blob/main/config.yaml) for override it `~/.config/gogh/config.yaml`. In this example base would be save in `~/files/second-brain.md/gogh.enc`

```yaml
database: ~/files/second-brain.md/gogh.enc
```
