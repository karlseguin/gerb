# Gerb
Gerb is a work in progress. Support for `for` and template inheritance are
missing.

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
to parse the same content (based on the content's hash). The cache will
automatically evict old items.

Once you have a template, you use the `Render` method to output the template
to the specified `io.Writer` using the specified data. Data must be a
`map[string]interface{}`.

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
precedence is left to right and parenthesis cannot be used.

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
In addition to custom builtins, it is possible to alias package functions. This
makes functions such as `strings.ToUpper` available to use.

By default, many functions from the `strings` package in addition to
`fmt.Sprintf` are available.

Since Go doesn't make it possible to automatically reflect functions exposed from
a package, registration must manually be done on a per-method basis. Like
builtins, this process is not thread safe:

    func init() {
      core.RegisterAliases("strings",
        "ToUpper", strings.ToUpper,
        "ToLower", strings.ToLower,
        ...
      )
      core.RegisterAliases("fmt",
        "Sprintf", fmt.Sprintf
      )
    }

You can see a list of what's currently aliased by looking at
<https://github.com/karlseguin/gerb/blob/master/core/aliases.go>

## Multiple Return Values
If you call a function which returns multiple values, only the first value is
considered/returned. The exception to this rule is with any assignment.

## Assignments
It's possible to create/assign to new or existing variables within a template:

    <% name := "leto" %>

Assignment supports multiple values, either via an explicit list of values or
a function which returns multiple values:

    <% name, power := "goku", 9000 %>
    <% name, power := sayan.Stats() %>

In fact, you can (but probably shouldn't) mix the two:

    <% name, power := sayan.Name(), 9000 %>

This is also true for assignments within an `if` or `else if` tag:

    <% if n, ok := strconv.Atoi(value); ok; { %>
      The number is <%= n %>
    <% } %>

## ++, --, += and -=
There's limited support for these four operators. As a general rule, they should
only ever be used on simple values (support was basically added to support the
i++ in a `for` loop).

Here's a couple examples of what **is not supported**:

    <% user.PowerLevel++ %>
    <% ranks[4]++ %>

Put differently, these 4 operators should only ever be used as such:

    <% counter++ %>


## Errors and Logs
`Render` should never fail. By default, `Render` will log errors to stdout. This
behavior can be changed. To disable all logging:

    gerb.Configure().Logger(nil)

Alternatively, to use your own logger:

    gerb.Configure().Logger(new(MyLogger))

Your logger must implement `Error(v ...interface{})`.
