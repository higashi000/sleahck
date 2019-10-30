# sleahck

## Description
This repository is take Slack channels and channels history for [sleahck.vim](https://github.com/higashi000/sleahck.vim).<br>

## Library
[gin-goinc/gin](https://github.com/gin-gonic/gin)<br>
[nlopes/slack](https://github.com/nlopes/slack)<br>

## Install
```
go get -u github.com/higashi000/sleahck
```

After, please edit your bashrc.<br>
```
export PATH=$PATH:$GOPATH/bin/
export SLACK_TOKEN="your slack token"
sleahck &
```

## How to use
### Get Channel List
Please access `http://localhost:8080/channelList`

### Get Channel History
Please access `http://localhost:8080/GetHistory/:channel_name`
