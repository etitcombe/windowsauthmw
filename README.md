# windowsauthmw

A Go middleware which adds the Windows domain and user name to the request context.
It assumes the web application is running under IIS via the [HttpPlatformHandler](https://docs.microsoft.com/en-us/iis/extensions/httpplatformhandler/httpplatformhandler-configuration-reference).

It is based on [windowsauthtoken](https://github.com/mfcollins3/windowsauthtoken)
by [Michael Collins](https://github.com/mfcollins3).

There is currently only one method: AddDomainUser which takes and returns an
[http.Handler](https://golang.org/pkg/net/http/#Handler).

Wrap your mux with it at the beginning or close to the beginning of your middleware
chain. You'll want this invoked early so that the information is added to the
Context before you need to use it.

    srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: windowsauthmw.AddDomainUser(mux)}

The middlware adds the domain and user name to the context in this format:
Domain\Username

It doesn't offer a way to get the value out of the context. That is the responsibility
of the application code. Use the provided DomainUserKey to get the value from the context:

    ctx := r.Context()
    domainUser := ctx.Value(windowsauthmw.DomainUserKey).(string)
