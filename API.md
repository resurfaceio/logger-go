# API

## Contents

<ul>
<li><a href="#creating_loggers">Creating Loggers</a></li>
<li><a href="#logging_http">Logging HTTP Calls</a></li>
<li><a href="#setting_default_rules">Setting Default Rules</a></li>
<li><a href="#setting_default_url">Setting Default URL</a></li>
<li><a href="#enabling_and_disabling_loggers">Enabling and Disabling Loggers</a></li>
</ul>

<a name="creating_loggers"/>

## Creating Loggers

To get started, first you'll need to create a `HttpLogger` instance. Here there are options to specify a URL (for where JSON
messages will be sent) and/or a specific set of <a href="https://resurface.io/rules.html">logging rules</a> (for what privacy
protections to apply). Default values will be used for either of these if specific values are not provided.

```golang
const resurfaceio = require('resurfaceio-logger');
const HttpLogger = resurfaceio.HttpLogger;

// with default url and rules
logger := NewHttpLogger(Options{});

// with specific url and default rules
opt := Options{
    Url: "https://...",
}
logger := NewHttpLogger(opt);

// with specific url and rules
opt := Options{
    Url: "https://...",
    Rules: "include_debug\n",
}
logger := NewHttpLogger(opt);

// with specific url and rules from local file
opt := Options{
    Url: "https://...",
    Rules: "file://./rules.txt\n",
}
logger := NewHttpLogger(opt);
```

<a name="logging_http"/>

## Logging HTTP Calls

Now that you have a logger instance, let's do some logging. Here you can pass standard request/response objects, as well
as response body and request body content when these are available.

```golang
const HttpMessage = resurfaceio.HttpMessage;

// with standard objects
SendHttpMessage(logger, response, request, start_time)
```

Request and Response bodies are automatically logged.

If standard request and response objects aren't available in your case, create mock implementations to pass instead.

```golang
// define request to log
request := &http.Request{}
request.Header.Set("Content-Type", "application/json")
request.Method = "POST"
request.Body = ioutil.NopCloser(strings.NewReader("body_content"))
request.URL, err = url.Parse("http://resurface.io")
if err != nil {/*handle error*/}


// define response to log
response = http.Response{}
response.Header.Set("Content-Type", "text/html: charset=utf-8")
response.StatusCode = 200

// log objects defined above
HttpMessage.send(logger, request, response);
```

<a name="setting_default_rules"/>

## Setting Default Rules

If no <a href="https://resurface.io/rules.html">rules</a> are provided when creating a logger, the default value of
`include strict` will be applied. A different default value can be specified as shown below.

```golang
HttpRules.SetDefaultRules("include_debug")
```

When specifying multiple default rules, put a new line character between each rule.

```golang
HttpRules.SetDefaultRules("include_debug\nallow_http_url")
```

<a name="setting_default_url"/>

## Setting Default URL

If your application creates more than one logger, or requires different URLs for different environments (development vs
testing vs production), then set the `USAGE_LOGGERS_URL` environment variable as shown below. This value will be applied if no
other URL is specified when creating a logger.

```bash
# from command line
export USAGE_LOGGERS_URL="https://..."

# for Heroku cli
heroku config:set USAGE_LOGGERS_URL=https://...
```

<a name="enabling_and_disabling_loggers"/>

## Enabling and Disabling Loggers

Individual loggers can be controlled through their `enable` and `disable` methods. When disabled, loggers will
not send any logging data, and the result returned by the `log` method will always be true (success).

All loggers for an application can be enabled or disabled at once with the `UsageLoggers` class. This even controls
loggers that have not yet been created by the application.

```golang
UsageLoggers.Disable();    // disable all loggers
UsageLoggers.Enable();     // enable all loggers
```

All loggers can be permanently disabled with the `USAGE_LOGGERS_DISABLE` environment variable. When set to true,
loggers will never become enabled, even if `UsageLoggers.enable()` is called by the application. This is primarily
done by automated tests to disable all logging even if other control logic exists.

```bash
# from command line
export USAGE_LOGGERS_DISABLE="true"

# for Heroku app
heroku config:set USAGE_LOGGERS_DISABLE=true
```

---

<small>&copy; 2016-2024 <a href="https://resurface.io">Graylog, Inc.</a></small>
