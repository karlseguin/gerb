# Gerb
Gerb is a work in progress. Critically, only output tags are currently supported.

## Usage

    template, err := gerb.ParseString("....", true)
    if err != nil {
      panic(err)
    }
    data := map[string]interface{}{
      "name": ....
    }
    template.Render(os.Stdout, data)

There are three available methods for creating a template:

1. Parse(data []byte, useCache bool)
2. ParseString(data string, useCache bool)
3. ParseFile(path string, useCache bool)

Unless `useCache` is set to `false`, an internal cache is used to avoid having
to reparse the same content (based on the content's hash). The cache is currently
dumb and will hold onto references forever.

Once you have a template, you use the `Render` method to output the template
to the specified `io.Writer` using the specified data. Data must be a `map[string]interface{}`.

It's safe to call `Render` from multiple threads.

## Output Tags
Gerb supports two types of output tags: escaped and non-escaped. The only difference
is that escaped tags will have < and > characters HTML-escaped.

    <%= "<script>alert('this will be escaped')</script>" %>
    <%! "<script>alert('this won't be escaped')</script>" %>

## Variables
Gerb attempts to behave as close to Go as possible. The biggest difference is that
method calls and field access is case-insensitive. Gerb supports:

* ints
* float64
* strings
* byte
* fields
* methods
* arrays
* basic operations

For example, the following works:

    <%= user.Analysis(count)[start+7:] %>

+, -, /, * and % are the only support operations. Currently (and sadly) order of
precendence is left to right and parenthesis cannot be used.

Gerb might occasionally be a little less strict with type conversions, but not
by much.

## Builtins
Go's builtins aren't natively available. However, custom builtin functions can
be registered. The custom buitlin `len` comes pre-registered and behaves much
like the real `len` builtin. You can register your own builtin:

    import "github.com/karlseguin/gerb/core"
    func init() {
      RegisterBuiltin("add", func(a, b int) int {
        return a + b
      })
    }

The builtin isn't threadsafe. It's expected that you'll register your builtins
at startup and then leave it alone.

## Aliases and Package Functions
In adddition to custom builtins, it is possible to alias package functions. This
makes functions such as `strings.ToUpper` available to use.

By default, many functions from the `strings` package in addition to
`fmt.Sprintf` are available.

Since Go doesn't make it possible to automatically reflect functions exposed from
a package, registration must manually be done on a per-method basis. Like
builtins, this process is not thread safe:

    func init() {
      core.Registerliases("strings",
        "ToUpper", strings.ToUpper,
        "ToLower", strings.ToLower,
        ...
      )
      core.Registerliases("fmt",
        "Sprintf", fmt.Sprintf
      )
    }

You can see a list of what's currently aliased by looking at
<https://github.com/karlseguin/gerb/blob/master/core/aliases.go>

## Multiple Return Values
If you call a function which returns multiple values, only the first value is
considered/returned.

This is likely to change, but the scope might be limited. Specifically, a
`var` tag might be added which would properly handle multiple return
values:

    `<% var x, y, _ = DoSomethihing() %>`

It's up in the air.

