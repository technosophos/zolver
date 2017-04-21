# Zolver: Get Personal with Your URLs
[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

[![Acid](https://img.shields.io/badge/Acid-Pass-brightgreen.svg)](http://localhost:7744/log/technosophos/zolver)


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

## Using

Right now, using Zolver takes three steps:

1. Write a `zolver.yaml` file.
2. Start the server.
3. Modify your `/etc/hosts` file.

### 1. Write your Zolver.yaml file

The main configuration file for Zolver is a simple YAML file.

Generally, each entry starts with the name of the domain you will use,
followed by an indented list of directives.

Starting with the most simple example, say I want to have `gh` redirect
me to `github.com`, so that I can type `gh/technosophos` and have it
take me to `https://github.com/technosophos`. The YAML for this is:

```yaml
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

```yaml
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

## 2. Start the server

Start the server as root, and tell it where your `zolver.yaml` is.

```
$ sudo $GOPATH/bin/zolver /path/to/zolver.yaml
```

You need to run Zolver as a root user so that you can take over port 80.

Every five minutes, the server will reload its configuration file, so
you can edit `zolver.yaml` and not have to reload the config.

## 3. Edit your hosts file with your simple domains:

Get your DNS resolver to play nicely with your new service by adding
an entry to your `/etc/hosts` file.

```
127.0.0.1 localhost gh my q
```

Don't worry, Zolver will tell you what the record should look like.

## Using Zolver

Now that Zolver is running, everything from Chrome, Firefox, and Safari
to Curl and WGet will be able to use your custom URLs.

## Advanced Usage

Now that you are familiar with the basics, learn some of the cool stuff.

### Using Redirects

A redirect maps a domain, and moves the path for you.

```yaml
gh:
  to: https://github.com
```

Now typing in `gh/Masterminds/sprig` will take you to
`https://github.com/Masterminds/sprig`. There is no need to do any
templating for this.

### Using the Shortener

The shortener allows you to pick one domain and provide a group of
mnemonic mappings under it:

```yaml
my:
  short:
    ts: http://technosophos.com
    news: http://techmeme.com
```

This maps `my/ts` to `http://technosophos.com` and `my/news` to
`http://techmeme.com`.

**This does not copy the path information from the source URL to the
destination.**

### Using Templates

When you use the `tpl` directive, you gain access to the template
engine. The template engine provides access to all of the pieces of the
URL that was passed in.

* `.Path`: The path portion of the URL passed in.
* `.RawQuery`: The raw query string
* `.Query`: The parsed query values
* `.Fragment`: Anything after a `#` in the URL
* `.Scheme`: The protocol scheme (usually HTTP).
* `.User`: The user portion of the URL (if you included one)
* `.Host`: The "host name", which is just the very first part of your
  URL.
* `.RequestURI`: `path?query` all together.
* `.String`: The entire URL.

```
scheme://[userinfo@]host/path[?query][#fragment]
```

There is also a function called `.Part` which takes one argument: An
integer (starting at 1) that indicates which part of the `.Path`.

This will reverse the order of `add/sprig/Masterminds` and covert it to
`https://github.com/Masterminds/sprig`:

```yaml
add:
  tpl: https://github.com/{{.Part 2}}/{{.Part 1}}
```

Along with those variables, there are a number of template functions
available to you, as well as some control structures. You may find
documentation for all of these here:

* Extended template functions from Sprig: https://github.com/Masterminds/sprig
* The core template engine: http://golang.org/pkg/text/template/

Here's an example that uses a number of template functions.

```yaml
add:
  tpl: https://github.com/{{"deis/" | repeat 2}}{{.Path | trimall "/" | lower}}
```

If you give it the URL `add/ISSUES` it will convert it to
`https://github.com/deis/deis/issues`.

The first template, `{{ "deis/" | repeat 2" }}` evaluates to
`deis/deis/`. And `{{ .Path | trimall "/" | lower }}` takes the path
(`/ISSUES`), trims off the leading and trailing slashes, and then
lowercases the rest.

### Hacking on Zolver

If you're a developer, you may prefer to check out the Git repo, and
then use Glide to set up your environment.

```
$ git clone https://github.com/technosophos/zolver.git
$ cd zolver
$ glide install
```

## License

This is licensed under an MIT-style license.
