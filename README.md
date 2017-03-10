# Concourse Monitor

Concourse Monitor is a simple tool that will sit in your OS system tray and
display the status of your concourse builds. When a build goes from non green
to green or green to non green then you will also receive a system
notification.

**Concourse Monitor currently only supports monitoring public pipelines.**

**Concourse Monitor is new and may have problems, feel free to create a GitHub
issue and I will address the problem as soon as I can.**

# Usage

Assuming you have Go installed:

```
go get github.com/bradylove/concourse-monitor/cmd/concourse-monitor
```

Then run `concourse-monitor`.

```
$ concourse-monitor --help
Usage of ./concourse-monitor:
  -concourse-url string
        url for concourse instance
  -refresh-interval int
        interval for pulling status from concourse (default 15)
  -team-name string
        team name to monitor (default "main")

$ concourse-monitor --concourse-url=http://concourse-domain.tld
```

# MIT License

Copyright (c) 2017 Brady Love <love.brady@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
