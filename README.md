# awssh
SSH to an EC2 instance in AWS without the IP, just the instance:{custom} tag

## Example
```bash
# SSH to server with your public key already registered 
$ awssh {instance tag}

# SSH to server using aws key (public key not registered)
$ awssh {instance tag} --key
```

## Installation
```bash
# get awssh
$ go get -v github.com/andhikagama/awssh

# create config dir
$ mkdir $HOME/awssh
```

create config.json inside config directory and replace values with your aws credential (see [`config.example.json`](https://github.com/andhikagama/awssh/blob/master/config.example.json))

## FAQ
### awssh: command not found
Make sure your `$GOPATH/bin` directory is already in `$PATH`


```bash
$ export PATH=$PATH:$(go env GOPATH)/bin
```

### Platforms
Only works on macOS

### iTerm2
By default `awssh` will launch macOS default Terminal app. In order to set iTerm2 as default terminal, you must first install [`duti`](http://duti.org/documentation.html).

```bash
$ brew install duti
```

Then download the [`iterm.duti`](https://github.com/andhikagama/awssh/blob/master/iterm.duti) file, and run it with `duti`

```bash
$ duti /path/to/your/iterm.duti
```