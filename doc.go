/*
Package csrf (gorilla/csrf) provides Cross Site Request Forgery (CSRF)
prevention middleware for Go web applications & services.

It includes:

    * The `csrf.Protect` middleware/handler provides CSRF protection on routes
      attached to a router or a sub-router.
    * A `csrf.Token` function that provides the token to pass into your response,
      whether that be a HTML form or a JSON response body.
    * ... and a `csrf.TemplateField` helper that you can pass into your `html/template`
      templates to replace a `{{ .csrfField }}` template tag with a hidden input
      field.

gorilla/csrf is easy to use: add the middleware to individual handlers with
the below:

    CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
    http.HandlerFunc("/route", CSRF(YourHandler))

... and then collect the token with `csrf.Token(r)` before passing it to the
template, JSON body or HTTP header (you pick!). gorilla/csrf inspects the form body
(first) and HTTP headers (second) on subsequent POST/PUT/PATCH/DELETE/etc. requests
for the token.

Here's the common use-case: HTML forms you want to provide CSRF protection for,
in order to protect malicious POST requests being made:

    package main

    import (
        "net/http"

        "github.com/gorilla/csrf"
    )

    func main() {
        r := mux.NewRouter()
        r.HandleFunc("/signup", ShowSignupForm)
        // All POST requests without a valid token will return HTTP 403 Forbidden.
        r.HandleFunc("/signup/post", SubmitSignupForm)

        // Add the middleware to your router by wrapping it.
        http.ListenAndServe(":8000",
            csrf.Protect([]byte("32-byte-long-auth-key"))(r))
    }

    func ShowSignupForm(w http.ResponseWriter, r *http.Request) {
        // signup_form.tmpl just needs a {{ .csrfField }} template tag for
        // csrf.TemplateField to inject the CSRF token into. Easy!
        t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{
            csrf.TemplateTag: csrf.TemplateField(r),
        })
    }

    func SubmitSignupForm(w http.ResponseWriter, r *http.Request) {
        // We can trust that requests making it this far have satisfied
        // our CSRF protection requirements.
    }

You can also send the CSRF token in the response header. This approach is useful
if you're using a front-end JavaScript framework like Ember or Angular, or are
providing a JSON API:

    package main

    import (
        "github.com/gorilla/csrf"
        "github.com/gorilla/mux"
    )

    func main() {
        r := mux.NewRouter()

        api := r.PathPrefix("/api").Subrouter()
        api.HandleFunc("/user/:id", GetUser).Methods("GET")

        http.ListenAndServe(":8000",
            csrf.Protect([]byte("32-byte-long-auth-key"))(r))
    }

    func GetUser(w http.ResponseWriter, r *http.Request) {
        // Authenticate the request, get the id from the route params,
        // and fetch the user from the DB, etc.

        // Get the token and pass it in the CSRF header. Our JSON-speaking client
        // or JavaScript framework can now read the header and return the token in
        // in its own "X-CSRF-Token" request header on the subsequent POST.
        w.Header().Set("X-CSRF-Token", csrf.Token(r))
        b, err := json.Marshal(user)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        w.Write(b)
    }

In addition: getting CSRF protection right is important, so here's some background:

    * This library generates unique-per-request (masked) tokens as a mitigation
      against the [BREACH attack](http://breachattack.com/).
    * The 'base' (unmasked) token is stored in the session, which means that
      multiple browser tabs won't cause a user problems as their per-request token
      is compared with the base token.
    * Operates on a "whitelist only" approach where safe (non-mutating) HTTP methods
      (GET, HEAD, OPTIONS, TRACE) are the *only* methods where token validation is not
      enforced.
    * The design is based on the battle-tested
      [Django](https://docs.djangoproject.com/en/1.8/ref/csrf/) and [Ruby on
      Rails](http://api.rubyonrails.org/classes/ActionController/RequestForgeryProtection.html)
      approaches.
    * Cookies are authenticated and based on the [securecookie](https://github.com/gorilla/securecookie)
      library. They're also Secure (issued over HTTPS only) and are HttpOnly
      by default, because sane defaults are important.
    * Go's `crypto/rand` library is used to generate the 32 byte (256 bit) tokens
      and the one-time-pad used for masking them.

This library does not seek to be adventurous.

*/
package csrf
