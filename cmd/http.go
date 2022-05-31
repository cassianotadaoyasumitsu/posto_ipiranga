package cmd

import (
	"fmt"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"
	"git.wealth-park.com/cassiano/posto_ipiranga/internal/ufo"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/middleware"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/postgres"

	"github.com/go-kit/kit/endpoint"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type httpCmdConfig struct {
	*baseConfig `mapstructure:",squash"`
	HTTPAddr    string `mapstructure:"http_addr"`
	DSN         string `mapstructure:"dsn"`
}

type httpCommand struct {
	*cobra.Command
	config *httpCmdConfig
}

func NewHTTPCommand() *httpCommand {
	c := &httpCommand{}
	c.Command = &cobra.Command{
		Use:   "http",
		Short: "Start the HTTP server",
		Long: fmt.Sprintf(`%s
Start the server`, Art()),
		Run: c.run,
	}
	c.config = &httpCmdConfig{
		baseConfig: &baseConfig{},
	}
	c.Flags().String("http-addr", ":8081", "HTTP server listen address")
	c.Flags().String("dsn", "", "The DSN for the PostgreSQL database (format: postgres://postgres:postgres@host:port/db?sslmode=mode)")
	unmarshalConfig(c.Command, "POSTO_IPIRANGA", &c.config)
	return c
}

func (c *httpCommand) run(cobraCmd *cobra.Command, args []string) {
	l := initLogger(cobraCmd, c.config.baseConfig)

	ret := 0
	defer func() {
		l.Info("terminated", "code", ret)
		os.Exit(ret)
	}()

	conn, err := postgres.NewConnection(c.config.DSN, 0, "posto_ipiranga-postgres")
	if err != nil {
		l.Error(err, log.MessageFieldName, "failed to connect to the database")
		ret = 1
		return
	}

	ufoRepo := postgres.NewUfoRepository(conn.DB)
	mSvc := ufo.NewService(ufoRepo)
	mwChain := endpoint.Chain(
		middleware.EndpointTimeMiddleware(l),
		middleware.EndpointLoggingMiddleware(l),
	)
	ufoEps := &ufo.Endpoints{
		CreateUfo: mwChain(ufo.NewCreateUfoEndpoint(mSvc)),
		OfID:      mwChain(ufo.NewOfIDEndpoint(mSvc)),
		List:      mwChain(ufo.NewListEndpoint(mSvc)),
	}

	httpLogger := log.With(l, "server", "http")
	r := mux.NewRouter()
	ufo.RegisterRoutes(r, ufoEps, httpLogger)
	muxHandler := gorillahandlers.CombinedLoggingHandler(os.Stdout, r)

	errs := make(chan error, 1)
	go func() {
		httpLogger.Info("Listening for HTTP requests", "address", c.config.HTTPAddr)
		errs <- http.ListenAndServe(c.config.HTTPAddr, muxHandler)
	}()

	sig := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		for s := range c {
			sig <- fmt.Errorf("%s", s)
		}
	}()

	halt := make(chan bool, 1)
	go func() {
		halted := false
		for {
			select {
			case err := <-errs:
				if !halted {
					ret = 1
					l.Error(
						err,
						"reason", "error",
						log.MessageFieldName, "terminating",
					)
				}
			case s := <-sig:
				if !halted {
					l.Warn(
						"terminating",
						"reason", "signal",
						"signal", s,
					)
				}
			}
			if !halted {
				halted = true
				close(halt)
			}
		}
	}()
	// Wait for the "halt" channel to be closed upon error or signal received,
	// this operation is blocking
	<-halt
}
