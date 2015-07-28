# Zolver: Your Personal URL Shortener

Your personal productivity is important. Don't waste time typing URLs.

Zolver is a simple tool for building local short links. You can use it
to...

- Create custom shortened URLs (`my/news` goes to `http://techmeme.com`).
- Redirect a short domain to a long one (`gh/technosophos` goes to
  `https://github.com/technosophos`)
- Rewrite a short local URL into a more complex one (`ticket/1` goes to
  `http://mybugtracker.example.com/issues/?item=1&view=full`)

Run it locally, use it in any browser. Or run it in the cloud and share
between computers.

## Installing

Assuming you have a valid Go environment:

```
$ go get -u github.com/technosophos/zolver
```

### Developers

If you're a developer, you may prefer to check out the Git repo, and
then use Glide to set up your environment.

```
$ cd zolver
$ glide install
```

## Using

Right now, using Zolver takes three steps:

### 1. Write your Zolver.yaml file

The main configuration file for Zolver is a simple YAML file.

Generally, each entry starts with the name of the domain you will use,
followed by an indented list of directives.

Starting with the most simple example, say I want to have `gh` redirect
me to `github.com`, so that I can type `gh/technosophos` and have it
take me to `https://github.com/technosophos`. The YAML for this is:

```
gh:
  to: https://github.com
```

The supported directives are:

* `to`: Redirect to a URL. If you add a path (`gh/technosophos/zolver`),
  the path will be moved to the new domain
  (`github.com/technosophos/zolver`).
* `short`: This provides a basic URL shortening map like you'd get with
  Bit.ly. `short`'s map takes the local path on the left, and the remote
  URL on the right.
* `tpl`: This allows you to build up templates where values from the URL
  you entered are put in very specific places in the result. It
  supports all of Go's `text/template` package plus the additional
  functions from https://github.com/Masterminds/sprig.

Here's a complate example that declares three domains: gh, my, and q.

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

Get your DNS resolver to play nicely with your new service by adding
an entry to your `/etc/hosts` file.

```
127.0.0.1 localhost gh my q
```

Don't worry, Zolver will tell you what the record should look like.

## 3. Start the server
```
sudo ./zolver
```

You need to run Zolver as a root user so that you can take over port 80.

## License

This is licensed under an MIT-style license.
