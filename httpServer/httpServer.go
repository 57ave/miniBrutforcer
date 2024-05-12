package main

import (
	"errors"
	"fmt"
	"context"
	"net"
	"io"
	"io/ioutil"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	io.WriteString(w, "This is my website!\n")
}

func getAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")

	fmt.Printf("%s: got /admin request. first(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first)
	if hasFirst && first == "adminRoot" {
		io.WriteString(w, "You find admin root\n")
	} else {
		w.Header().Set("not-admin-query", "first")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /name request\n", ctx.Value(keyServerAddr))

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. , body: %s\n",
		ctx.Value(keyServerAddr),
		body)
	io.WriteString(w, "This is my website!\n")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/admin", getAdmin)
	mux.HandleFunc("/name", getName)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server one closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server one: %s\n", err)
	}
}
