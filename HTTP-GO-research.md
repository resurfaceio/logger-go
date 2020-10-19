## [a guide to writing logging middleware in GO](https://blog.questionable.services/article/guide-logging-middleware-go/)

## [GO log package](https://golang.org/pkg/log/)
- [SetOutput](https://golang.org/src/log/log.go?s=8750:8777#L266) provides a way to route the out put of the logger (potentially will help with packaging the data of the transaction being recorded)
- [New](https://golang.org/src/log/log.go?s=3133:3189#L55) can also be used to set output destination when creating a new logger. The first argument specifies the destination

## [kit/log](https://github.com/go-kit/kit/tree/master/log)
- [go-kit json_logger](https://github.com/go-kit/kit/blob/master/log/json_logger.go)
- 

## [honeybdger.io](https://www.honeybadger.io/blog/golang-logging/)
- Popular Logging frameworks
    - [glog](https://github.com/golang/glog)
    - [logrus](https://github.com/Sirupsen/logrus)
- [Logging Framework Options](https://github.com/avelino/awesome-go#logging)

## [Golang logging http responses (in addition to requests)](https://stackoverflow.com/questions/38443889/golang-logging-http-responses-in-addition-to-requests)

## [Attaching Logger to Router](https://golangcode.com/attach-a-logger-to-your-router/)

## [How to collect, standardize, and centralize Golang logs](https://www.datadoghq.com/blog/go-logging/)
- [Syslog Package](https://golang.org/pkg/log/syslog/) can be used to centralize logs