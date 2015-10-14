# [Stackdriver](http://www.stackdriver.com/) API Library

**DEPRECATED**: This repo is no longer in use and is not being maintained. The authoritative version of this code base can be found at [https://github.com/AMeng/stackdriver](https://github.com/AMeng/stackdriver).


A client library for accessing the Stackdriver API.

Currently implemented:
* Custom Metrics
* Code Deploy Events
* Annotation Events


## Usage

```go
import (
    "github.com/bellycard/stackdriver"
    "time"
)

// Create new Stackdriver API client.
client := stackdriver.NewStackdriverClient("apikey")

// Create new Stackdriver API gateway message.
apiMessages := stackdriver.NewGatewayMessage()

now := time.Now().Unix()

// Populate gateway message with metrics.
apiMessages.CustomMetric("my-metric1", "i-axd939f", now, 50)
apiMessages.CustomMetric("my-metric2", "i-afdsf9f", now, 6.5)
apiMessages.CustomMetric("my-metric3", "i-a3d923f", now, 25)

// Send gateway message to Stackdriver API.
client.Send(apiMessages)
```


## Author(s) & Credit

**Christian Vozar**

+ [http://twitter.com/christianvozar](http://twitter.com/christianvozar)
+ [http://github.com/christianvozar](http://github.com/christianvozar)

## Copyright and License

Copyright 2014 Belly, Inc. under [the Apache 2.0 license](LICENSE.md).
