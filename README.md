# deployster-slack

A very rough prototype for what a deployster plugin might look like using etcd to emit deploy events.


## Configuration

* `ETCDCTL_PEERS`: location of etcd peers
* `SLACK_TOKEN`: an API token for the slack bot to post messages as
* `SLACK_CHANNEL`: the ID of the channel that messages should be posted to


## Running

```
go build
./deployster-slack
```
