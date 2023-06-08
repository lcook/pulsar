/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package relay

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitMux(resp any, hooks []Hook, config, port string) (*http.Server, error) {
	srv, err := registerMux(resp, hooks, config, port)
	if err != nil {
		return srv, err
	}
	/*
	 * Spawn a channel listening for CTRL-C key-presses
	 * (interrupts) and SIGTERM signals.  If so, gracefully
	 * shutdown the server.
	 */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go shutdown(srv, sig)

	return srv, nil
}

func registerMux(resp any, hooks []Hook, config, port string) (*http.Server, error) {
	mux := http.NewServeMux()
	/*
	 * Register the `Response` handler function with it's corresponding
	 * endpoint in each of the hooks provided.
	 *
	 * There is a few middlewares a hook can use:
	 *
	 * `OptionCheckMethod`: Verify the incoming payload is a `POST` method.
	 * `OptionCheckType`: Verify the incoming payload is of type `application/json`.
	 */
	for _, hook := range hooks {
		err := hook.LoadConfig(config)
		if err != nil {
			return nil, err
		}

		mux.HandleFunc(hook.Endpoint(), func(httpfn http.HandlerFunc) http.HandlerFunc {
			return func(writer http.ResponseWriter, req *http.Request) {
				if (hook.Options()&OptionCheckMethod != 0) &&
					(req.Method != http.MethodPost) {
					writer.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				if (hook.Options()&OptionCheckType != 0) &&
					(req.Header.Get("Content-Type") != "application/json") {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
				httpfn(writer, req)
			}
		}(hook.Response(resp)))
	}

	return &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}, nil
}

func shutdown(server *http.Server, sig <-chan os.Signal) {
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Could not gracefully shutdown...", err)
	}
}
