package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/AlehaWP/YaPracticum.git/internal/handlers"
	"github.com/AlehaWP/YaPracticum.git/internal/middlewares"
	"github.com/AlehaWP/YaPracticum.git/internal/models"
	"github.com/go-chi/chi"
)

type Server struct {
	http.Server
}

//Start server with router.
func (s *Server) Start(ctx context.Context, repo models.Repository, opt models.Options) {
	r := chi.NewRouter()
	handlers.NewHandlers(repo, opt)
	middlewares.NewCookie(repo)
	r.Use(middlewares.SetCookieUser, middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)

	r.HandleFunc("/debug/pprof/*", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	// r.Use(middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)

	r.Get("/api/user/urls", handlers.HandlerUserPostURLs)
	r.Get("/ping", handlers.HandlerCheckDBConnect)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middlewares.URLCtx)
		r.Get("/", handlers.HandlerURLGet)
	})
	r.Post("/", handlers.HandlerURLPost)
	r.Get("/api/user/urls", handlers.HandlerUserPostURLs)
	r.Post("/api/shorten", handlers.HandlerAPIURLPost)
	r.Post("/api/shorten/batch", handlers.HandlerAPIURLsPost)
	r.Delete("/api/user/urls", handlers.HandlerDeleteUserUrls)

	s.Addr = opt.ServAddr()
	s.Handler = r
	go s.ListenAndServe()

	<-ctx.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	s.Shutdown(ctx)
}
