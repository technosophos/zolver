# Zolver: Your Personal URL Shortener

Zolver is a simple tool for building local short links. You can use it
to...

- Create custom shortened URLs (`my/news` goes to `http://techmeme.com`).
- Redirect a short domain to a long one (`gh/technosophos` goes to
  `https://github.com/technosophos`)
- Rewrite a short local URL into a more complex one (`ticket/1` goes to
  `http://mybugtracker.example.com/issues/?item=1&view=full`)

## Installing

```
$ go install github.com/technosophos/zolver
```

## Using

Right now, using Zolver takes three steps:

### 1. Write your Zolver.yaml file

```
# Simple redirect
gh:
  to: https://github.com
# Shortened URLs
my:
  # Shorten or just search.
  short:
    # my/ts goes here
    ts: http://technosophos.com
    # my/news goes here
    news: http://techmeme.com
# rewrite a URL based on the contents of an incomming URL
q:
  tpl: https://ddg.gg?q={{.Path}}
```

## 2. Edit your hosts file with your simple domains:

```
127.0.0.1 localhost gh my q
```

## 3. Start the server
```
sudo ./zolver
```

You need to run Zolver as a root user so that you can take over port 80.
