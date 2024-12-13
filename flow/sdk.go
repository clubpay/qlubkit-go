package flow

import (
	"context"
	"time"

	"github.com/clubpay/go-service-business/pkg/core/settings"
	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.temporal.io/api/namespace/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Config struct {
	fx.In

	Settings settings.Settings
	Logger   *log.Logger
}

type SDK struct {
	app    *fx.App
	nsCli  client.NamespaceClient
	cli    client.Client
	replay worker.WorkflowReplayer
	w      worker.Worker
	l      *log.Logger

	taskQ     string
	namespace string
}

func newSDK(
	lc fx.Lifecycle, cfg Config, namespace, taskQ string,
) (*SDK, error) {
	sdk := &SDK{
		taskQ:     taskQ,
		namespace: namespace,
		replay:    worker.NewWorkflowReplayer(),
		l:         cfg.Logger,
	}

	err := sdk.invoke(lc, cfg.Settings)
	if err != nil {
		return nil, err
	}

	return sdk, nil
}

func (sdk *SDK) invoke(lc fx.Lifecycle, set settings.Settings) error {
	var err error
	sdk.nsCli, err = client.NewNamespaceClient(client.Options{
		HostPort: settings.GetTemporalHostPort(set),
		Logger:   log.NopLogger.Sugared(),
	})
	if err != nil {
		return err
	}

	if _, err = sdk.nsCli.Describe(context.Background(), sdk.namespace); err != nil {
		err = sdk.nsCli.Register(
			context.Background(),
			&workflowservice.RegisterNamespaceRequest{
				Namespace:                        sdk.namespace,
				WorkflowExecutionRetentionPeriod: &durationpb.Duration{Seconds: 72 * 3600},
			},
		)
		if err != nil {
			sdk.l.Info(
				"got error on create namespace",
				log.String("hostport", settings.GetTemporalHostPort(set)),
				log.Error(err),
			)
		}
	}

	sdk.cli, err = client.NewLazyClient(
		client.Options{
			HostPort:  settings.GetTemporalHostPort(set),
			Namespace: sdk.namespace,
			Logger:    log.NopLogger.Sugared(),
		},
	)
	if err != nil {
		return err
	}

	sdk.w = worker.New(
		sdk.cli,
		sdk.taskQ,
		worker.Options{
			DisableRegistrationAliasing: true,
		},
	)

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				return sdk.w.Start()
			},
			OnStop: func(_ context.Context) error {
				sdk.w.Stop()

				return nil
			},
		},
	)

	return nil
}

func (sdk *SDK) TaskQueue() string {
	return sdk.taskQ
}

type UpdateNamespaceRequest struct {
	Description                      *string
	WorkflowExecutionRetentionPeriod *time.Duration
}

func (sdk *SDK) UpdateWorkflowRetentionPeriod(ctx context.Context, d time.Duration) error {
	res, err := sdk.nsCli.Describe(ctx, sdk.namespace)
	if err != nil {
		return err
	}

	if res.Config == nil {
		res.Config = &namespace.NamespaceConfig{}
	}
	res.Config.WorkflowExecutionRetentionTtl = durationpb.New(d)

	return sdk.nsCli.Update(
		ctx,
		&workflowservice.UpdateNamespaceRequest{
			Namespace: sdk.namespace,
			Config:    res.Config,
		},
	)
}
