package windowsauthmw

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/sys/windows"
)

// ContextKeyType is the type used to create context keys.
type ContextKeyType string

// DomainUserKey is the key to use to get the domain user from the conext.
const DomainUserKey ContextKeyType = "domain-user"

// AddDomainUser adds the Windows domain and user name to the request context.
func AddDomainUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("X-Iis-Windowsauthtoken")
		if tokenHeader != "" {
			handle, err := strconv.ParseUint(tokenHeader, 16, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			token := windows.Token(handle)
			defer token.Close()

			user, err := token.GetTokenUser()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			username, domainName, _, err := user.User.Sid.LookupAccount("")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			domainUser := fmt.Sprintf("%s\\%s", domainName, username)

			ctx := context.WithValue(r.Context(), DomainUserKey, domainUser)
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	})
}
