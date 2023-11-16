package csi

import (
	"context"
	"net/http"
	"time"

	"github.com/kubernetes-csi/external-resizer/pkg/controller"
	"github.com/kubernetes-csi/external-resizer/pkg/csi"
	"github.com/kubernetes-csi/external-resizer/pkg/resizer"
	"k8s.io/client-go/informers"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	resyncPeriod           = 10 * time.Minute
	retryIntervalStart     = time.Second
	retryIntervalMax       = 5 * time.Minute
	handleVolumeInUseError = true
	workers                = 1
)

type ResizerOptions struct {
	SharedInformerFactory informers.SharedInformerFactory
	DriverName            string
	CSIEndpoint           string
	CSIOperationTimeout   time.Duration // 10*time.Second
}

type Resizer struct {
	config  *rest.Config
	client  *http.Client
	options ResizerOptions
}

func (r *Resizer) NeedLeaderElection() bool {
	return true
}

var _ manager.Runnable = &Resizer{}
var _ manager.LeaderElectionRunnable = &Resizer{}

func NewResizer(mgr manager.Manager, options ResizerOptions) *Resizer {
	return &Resizer{
		config:  mgr.GetConfig(),
		client:  mgr.GetHTTPClient(),
		options: options,
	}
}

func (r *Resizer) Start(ctx context.Context) error {
	csiClient, err := csi.New(r.options.CSIEndpoint, r.options.CSIOperationTimeout, nil)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfigAndClient(r.config, r.client)
	if err != nil {
		return err
	}

	csiResizer, err := resizer.NewResizerFromClient(
		csiClient,
		r.options.CSIOperationTimeout,
		clientset,
		nil,
		r.options.DriverName)
	if err != nil {
		return err
	}

	resizerController := controller.NewResizeController(r.options.DriverName, csiResizer, clientset, resyncPeriod, r.options.SharedInformerFactory,
		workqueue.NewItemExponentialFailureRateLimiter(retryIntervalStart, retryIntervalMax),
		handleVolumeInUseError)

	resizerController.Run(workers, ctx)

	ctrl.Log.Info("resizer finished shutdown")

	return nil
}
