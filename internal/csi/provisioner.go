package csi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/external-provisioner/pkg/capacity"
	"github.com/kubernetes-csi/external-provisioner/pkg/capacity/topology"
	provisionctrl "github.com/kubernetes-csi/external-provisioner/pkg/controller"
	"github.com/kubernetes-csi/external-provisioner/pkg/owner"
	snapclientset "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	csitrans "k8s.io/csi-translation-lib"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlRuntimeMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v9/controller"
	libmetrics "sigs.k8s.io/sig-storage-lib-external-provisioner/v9/controller/metrics"
)

const (
	ResyncPeriodOfCsiNodeInformer = 1 * time.Hour
)

type ProvisionerOptions struct {
	DriverName          string
	CSIEndpoint         string
	CSIOperationTimeout time.Duration
}

type Provisioner struct {
	config  *rest.Config
	client  *http.Client
	options ProvisionerOptions
}

func (p *Provisioner) NeedLeaderElection() bool {
	return true
}

func init() {
	ctrlRuntimeMetrics.Registry.MustRegister([]prometheus.Collector{
		libmetrics.PersistentVolumeClaimProvisionTotal,
		libmetrics.PersistentVolumeClaimProvisionFailedTotal,
		libmetrics.PersistentVolumeClaimProvisionDurationSeconds,
		libmetrics.PersistentVolumeDeleteTotal,
		libmetrics.PersistentVolumeDeleteFailedTotal,
		libmetrics.PersistentVolumeDeleteDurationSeconds,
	}...)
}

var _ manager.Runnable = &Provisioner{}
var _ manager.LeaderElectionRunnable = &Provisioner{}

func NewProvisioner(mgr manager.Manager, options ProvisionerOptions) *Provisioner {
	return &Provisioner{
		config:  mgr.GetConfig(),
		client:  mgr.GetHTTPClient(),
		options: options,
	}
}

func (p *Provisioner) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)

	onLostConnection := func() bool {
		logger.Info("lost connection to csi driver, attempting to reconnect due to in tree provisioning...")
		return true
	}
	grpcClient, err := connection.Connect(p.options.CSIEndpoint, nil,
		connection.OnConnectionLoss(onLostConnection),
		connection.WithTimeout(p.options.CSIOperationTimeout))
	defer grpcClient.Close() //nolint:errcheck,staticcheck
	if err != nil {
		return err
	}
	pluginCapabilities, controllerCapabilities, err := provisionctrl.GetDriverCapabilities(grpcClient, p.options.CSIOperationTimeout)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfigAndClient(p.config, p.client)
	if err != nil {
		return err
	}
	snapClient, err := snapclientset.NewForConfigAndClient(p.config, p.client)
	if err != nil {
		return err
	}

	// Generate a unique ID for this provisioner
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	identity := strconv.FormatInt(timeStamp, 10) + "-" + strconv.Itoa(rand.Intn(10000)) + "-" + p.options.DriverName

	translator := csitrans.New()
	factory := informers.NewSharedInformerFactory(clientset, ResyncPeriodOfCsiNodeInformer)

	scLister := factory.Storage().V1().StorageClasses().Lister()
	csiNodeLister := factory.Storage().V1().CSINodes().Lister()
	nodeLister := factory.Core().V1().Nodes().Lister()
	claimLister := factory.Core().V1().PersistentVolumeClaims().Lister()
	vaLister := factory.Storage().V1().VolumeAttachments().Lister()
	provisioner := provisionctrl.NewCSIProvisioner(
		clientset,
		p.options.CSIOperationTimeout,
		identity,
		"pvc",
		-1,
		grpcClient,
		snapClient,
		p.options.DriverName,
		pluginCapabilities,
		controllerCapabilities,
		"",
		false,
		true,
		translator,
		scLister,
		csiNodeLister,
		nodeLister,
		claimLister,
		vaLister,
		nil,
		false,
		"",
		nil,
		false,
		false,
	)

	var capacityController *capacity.Controller
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return fmt.Errorf("need NAMESPACE env variable for CSIStorageCapacity objects")
	}
	podName := os.Getenv("NAME")
	if podName == "" {
		return fmt.Errorf("need NAME env variable to determine CSIStorageCapacity owner")
	}
	controllerRef, err := owner.Lookup(p.config, namespace, podName,
		schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}, 2)
	if err != nil {
		return fmt.Errorf("look up owner(s) of pod failed: %w", err)
	}
	logger.Info(fmt.Sprintf("using %s/%s %s as owner of CSIStorageCapacity objects", controllerRef.APIVersion, controllerRef.Kind, controllerRef.Name))
	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(time.Second, 5*time.Minute)

	topologyInformer := topology.NewNodeTopology(
		p.options.DriverName,
		clientset,
		factory.Core().V1().Nodes(),
		factory.Storage().V1().CSINodes(),
		workqueue.NewRateLimitingQueueWithConfig(rateLimiter, workqueue.RateLimitingQueueConfig{
			Name: "csitopology",
		}),
	)
	go topologyInformer.RunWorker(ctx)

	managedByID := "external-provisioner"

	factoryForNamespace := informers.NewSharedInformerFactoryWithOptions(clientset,
		ResyncPeriodOfCsiNodeInformer,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = labels.Set{
				capacity.DriverNameLabel: p.options.DriverName,
				capacity.ManagedByLabel:  managedByID,
			}.AsSelector().String()
		}),
	)
	// We use the V1 CSIStorageCapacity API if available.
	clientFactory := capacity.NewV1ClientFactory(clientset)
	cInformer := factoryForNamespace.Storage().V1().CSIStorageCapacities()

	capacityController = capacity.NewCentralCapacityController(
		csi.NewControllerClient(grpcClient),
		p.options.DriverName,
		clientFactory,
		// Metrics for the queue is available in the default registry.
		workqueue.NewRateLimitingQueueWithConfig(rateLimiter, workqueue.RateLimitingQueueConfig{
			Name: "csistoragecapacity",
		}),
		controllerRef,
		managedByID,
		namespace,
		topologyInformer,
		factory.Storage().V1().StorageClasses(),
		cInformer,
		time.Minute,
		false,
		p.options.CSIOperationTimeout,
	)
	// Wrap Provision and Delete to detect when it is time to refresh capacity.
	provisioner = capacity.NewProvisionWrapper(provisioner, capacityController)

	claimRateLimiter := workqueue.NewItemExponentialFailureRateLimiter(time.Second, 5*time.Minute)
	claimQueue := workqueue.NewRateLimitingQueueWithConfig(claimRateLimiter, workqueue.RateLimitingQueueConfig{
		Name: "claims",
	})
	claimInformer := factory.Core().V1().PersistentVolumeClaims().Informer()

	provisionController := controller.NewProvisionController(
		clientset,
		p.options.DriverName,
		provisioner,
		controller.LeaderElection(false), // Always disable leader election in provisioner lib. Leader election should be done here in the CSI provisioner level instead.
		controller.FailedProvisionThreshold(0),
		controller.FailedDeleteThreshold(0),
		controller.RateLimiter(rateLimiter),
		controller.Threadiness(1),
		controller.CreateProvisionedPVLimiter(workqueue.DefaultControllerRateLimiter()),
		controller.ClaimsInformer(claimInformer),
		controller.NodesLister(nodeLister),
		controller.AddFinalizer(true),
	)

	csiClaimController := provisionctrl.NewCloningProtectionController(
		clientset,
		claimLister,
		claimInformer,
		claimQueue,
		controllerCapabilities,
	)

	factory.Start(ctx.Done())
	// Starting is enough, the capacity controller will
	// wait for sync.
	factoryForNamespace.Start(ctx.Done())
	cacheSyncResult := factory.WaitForCacheSync(ctx.Done())
	for _, v := range cacheSyncResult {
		if !v {
			return fmt.Errorf("failed to sync Informers")
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		capacityController.Run(ctx, 1)
		logger.Info("capacity controller finished shutdown")
	}()
	go func() {
		defer wg.Done()
		csiClaimController.Run(ctx, 1)
		logger.Info("csi claim controller finished shutdown")
	}()
	go func() {
		defer wg.Done()
		provisionController.Run(ctx)
		logger.Info("provisioner controller finished shutdown")
	}()

	wg.Wait()
	logger.Info("provisioner finished shutdown")

	return nil
}
