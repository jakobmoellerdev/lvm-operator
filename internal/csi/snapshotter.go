package csi

import (
	"context"
	"net/http"
	"time"

	"github.com/kubernetes-csi/csi-lib-utils/connection"
	clientset "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	informers "github.com/kubernetes-csi/external-snapshotter/client/v6/informers/externalversions"
	controller "github.com/kubernetes-csi/external-snapshotter/v6/pkg/sidecar-controller"
	"github.com/kubernetes-csi/external-snapshotter/v6/pkg/snapshotter"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	ResyncPeriodOfSnapshotInformer = 15 * time.Minute
)

type SnapshotterOptions struct {
	DriverName          string
	CSIEndpoint         string
	CSIOperationTimeout time.Duration
}

type Snapshotter struct {
	config  *rest.Config
	client  *http.Client
	options SnapshotterOptions
}

func (s *Snapshotter) NeedLeaderElection() bool {
	return true
}

var _ manager.Runnable = &Snapshotter{}
var _ manager.LeaderElectionRunnable = &Snapshotter{}

func NewSnapshotter(mgr manager.Manager, options SnapshotterOptions) *Snapshotter {
	return &Snapshotter{
		config:  mgr.GetConfig(),
		client:  mgr.GetHTTPClient(),
		options: options,
	}
}

func (s *Snapshotter) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	onLostConnection := func() bool {
		logger.Info("lost connection to csi driver, attempting to reconnect due to in tree provisioning...")
		return true
	}
	grpcClient, err := connection.Connect(s.options.CSIEndpoint, nil,
		connection.OnConnectionLoss(onLostConnection),
		connection.WithTimeout(s.options.CSIOperationTimeout))
	defer grpcClient.Close() //nolint:errcheck,staticcheck
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfigAndClient(s.config, s.client)
	if err != nil {
		return err
	}

	snapClient, err := clientset.NewForConfigAndClient(s.config, s.client)
	if err != nil {
		return err
	}
	factory := informers.NewSharedInformerFactory(snapClient, ResyncPeriodOfSnapshotInformer)
	snapShotter := snapshotter.NewSnapshotter(grpcClient)

	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(time.Second, 5*time.Minute)
	ctrl := controller.NewCSISnapshotSideCarController(
		snapClient,
		kubeClient,
		s.options.DriverName,
		factory.Snapshot().V1().VolumeSnapshotContents(),
		factory.Snapshot().V1().VolumeSnapshotClasses(),
		snapShotter,
		nil,
		s.options.CSIOperationTimeout,
		ResyncPeriodOfSnapshotInformer,
		"snapshot",
		-1,
		"groupsnapshot",
		-1,
		false,
		rateLimiter,
		false,
		factory.Groupsnapshot().V1alpha1().VolumeGroupSnapshotContents(),
		factory.Groupsnapshot().V1alpha1().VolumeGroupSnapshotClasses(),
		rateLimiter,
	)

	// run...
	factory.Start(ctx.Done())
	ctrl.Run(1, ctx.Done())

	return nil
}
