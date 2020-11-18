package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	promclient "github.com/grafana/loki/pkg/promtail/client"
	"github.com/prometheus/common/model"
)

var (
	defaultLokiURL = "http://127.0.0.1:3100/loki/api/v1/push"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogfmtLogger(os.Stdout)

	// TODO(kavi): remove hardcoded config

	// Sets your Google Cloud Platform project ID.
	projectID := "grafanalabs-dev"

	// Creates a pubsub client.
	cloudclient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		level.Error(logger).Log("event", "failed to create pubsub client", "cause", err)
		os.Exit(1)
	}

	sub := cloudclient.SubscriptionInProject("loki-subscription", projectID)

	lu := os.Getenv("LOKI_URL")
	if lu == "" {
		lu = defaultLokiURL
	}

	u, err := url.ParseRequestURI(lu)
	if err != nil {
		level.Error(logger).Log("event", "failed to parse LOKI_URL", "cause", err)
		os.Exit(1)
	}

	// creates promtail client
	pclient, err := promclient.New(promclient.Config{
		URL:     flagext.URLValue{URL: u},
		Timeout: 2 * time.Second,
	}, logger)

	// TODO(kavi): Add goroutines lifecycle with run.Group
	// with signal handler for graceful shutdown.
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			// TODO(kavi): send it o pclient.Handle
			fmt.Println(string(m.Data))
			fmt.Println()
			pclient.Handle(model.LabelSet{"source": "cloudtail"}, time.Now(), string(m.Data))
			m.Ack()
		})

		if err != nil {
			level.Error(logger).Log("event", "failed to receive from subscription", "cause", err)
		}

		cancel()
	}()

	<-ctx.Done()
}
