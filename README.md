# GOGH

> Share files via GitHub API

> All files are encrypted with AES256

## Installation

1. Download the release via the [release page](https://github.com/Devil666face/gogh/releases/latest).
2. Unzip the file: `tar -xf gogh_linux_amd64.tar.gz`
3. Add it to your path: `sudo cp gogh /usr/local/bin` or use it directly: `./gogh`

## Usage

#### Settings

1. Set your GitHub session cookie token

- Open your browser
- Log in to GitHub and open the GitHub page
- Press F12
- Go to Store -> Cookies -> user_session
- Copy the value, which looks like `bzzsDpsejMS8jiKH-eezHwUlQWERTYvPix0qGPxg8YyL-b-7`

![](https://github.com/Devil666face/gogh/blob/main/.github/cookies.png)

2. Set the token in Gogh

- `settings`
- `set token`
- Enter your token: `bzzsDpsejMS8jiKH-eezHwUlQWERTYvPix0qGPxg8YyL-b-7`

```
❯ gogh
🪣 >> settings
🪣 >> ⚙️ settings >> set token
```

3. Check the token

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
- Choose a file

```
❯ gogh
🪣 >> store
🪣 >> 🗄️ store >> upload
👌 uploaded test.txt
```

##### Download file

- `store`
- `download`
- Choose a file

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

> You can share Gogh files with all links and AES keys

##### Export

- `share`
- `export`
- Choose an uploaded remote file
- The exported file will be saved as `file_name.gogh`
- Send `file_name.gogh` to your friend

```
❯ gogh
🪣 >> share
🪣 >> 🔁 share >> export
💾 file test.txt saved
```

##### Import

- `share`
- `import`
- Choose a local file in the current directory with the `.gogh` extension
- You can now download the file via `store`

```
❯ gogh
🪣 >> share
🪣 >> 🔁 share >> import
📥 file test.txt imported
```

#### Configuration

1. Gogh saves the database in `~/.config/gogh/gogh.enc`.
2. Use the [config file](https://github.com/Devil666face/gogh/blob/main/config.yaml) to override this to `~/.config/gogh/config.yaml`. For example, the base file could be saved as `~/files/second-brain.md/gogh.enc`.

```yaml
database: ~/files/second-brain.md/gogh.enc
```
