package persistent_volume_claim_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/openshift/lvm-operator/controllers"
	persistentvolumeclaim "github.com/openshift/lvm-operator/controllers/persistent-volume-claim"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPersistentVolumeClaimReconciler_Reconcile(t *testing.T) {
	defaultNamespace := "openshift-storage"

	tests := []struct {
		name                     string
		req                      controllerruntime.Request
		objs                     []client.Object
		wantErr                  bool
		expectNoStorageAvailable bool
	}{
		{
			name: "testing set deletionTimestamp",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-deletionTimestamp",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-deletionTimestamp",
						DeletionTimestamp: &metav1.Time{Time: time.Now()}, Finalizers: []string{"random-finalizer"}},
				},
			},
		},
		{
			name: "testing empty storageClassName",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-emptyStorageClassName",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-emptyStorageClassName"},
				},
			},
		},
		{
			name: "testing non-applicable storageClassName",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-nonApplicableStorageClassName",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-nonApplicableStorageClassName"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String("blabla"),
					},
				},
			},
		},
		{
			name: "testing non-pending PVC is skipped",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-non-pending-PVC",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-non-pending-PVC"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String(controllers.StorageClassPrefix + "bla"),
					},
					Status: v1.PersistentVolumeClaimStatus{
						Phase: v1.ClaimBound,
					},
				},
			},
		},
		{
			name: "testing pending PVC is proccessed",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-pending-PVC",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-pending-PVC"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String(controllers.StorageClassPrefix + "bla"),
					},
					Status: v1.PersistentVolumeClaimStatus{
						Phase: v1.ClaimPending,
					},
				},
			},
			expectNoStorageAvailable: true,
		},
		{
			name: "testing PVC requesting more storage than capacity in the node",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-pending-PVC",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-pending-PVC"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String(controllers.StorageClassPrefix + "bla"),
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceStorage: *resource.NewQuantity(100, resource.DecimalSI),
							},
						},
					},
					Status: v1.PersistentVolumeClaimStatus{
						Phase: v1.ClaimPending,
					},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{persistentvolumeclaim.CapacityAnnotation + "bla": "10"}},
				},
			},
			expectNoStorageAvailable: true,
		},
		{
			name: "testing PVC requesting less storage than capacity in the node",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-pending-PVC",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-pending-PVC"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String(controllers.StorageClassPrefix + "bla"),
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceStorage: *resource.NewQuantity(10, resource.DecimalSI),
							},
						},
					},
					Status: v1.PersistentVolumeClaimStatus{
						Phase: v1.ClaimPending,
					},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{persistentvolumeclaim.CapacityAnnotation + "bla": "100"}},
				},
			},
			expectNoStorageAvailable: false,
		},
		{
			name: "testing PVC requesting less storage than capacity in one node, having another node without annotation",
			req: controllerruntime.Request{NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      "test-pending-PVC",
			}},
			objs: []client.Object{
				&v1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace, Name: "test-pending-PVC"},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: pointer.String(controllers.StorageClassPrefix + "bla"),
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceStorage: *resource.NewQuantity(10, resource.DecimalSI),
							},
						},
					},
					Status: v1.PersistentVolumeClaimStatus{
						Phase: v1.ClaimPending,
					},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{Name: "annotated", Annotations: map[string]string{persistentvolumeclaim.CapacityAnnotation + "bla": "100"}},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{Name: "idle"},
				},
			},
			expectNoStorageAvailable: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := record.NewFakeRecorder(1)
			r := &persistentvolumeclaim.PersistentVolumeClaimReconciler{
				Client:   fake.NewClientBuilder().WithObjects(tt.objs...).Build(),
				Recorder: recorder,
			}
			got, err := r.Reconcile(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, controllerruntime.Result{}) {
				t.Errorf("Reconcile() got non default Result")
				return
			}

			select {
			case event := <-recorder.Events:
				if !strings.Contains(event, "NotEnoughCapacity") {
					t.Errorf("event was captured but it did not contain the reason NotEnoughCapacity")
					return
				}
				if !tt.expectNoStorageAvailable {
					t.Errorf("event was captured but was not expected")
					return
				}
			case <-time.After(100 * time.Millisecond):
				if tt.expectNoStorageAvailable {
					t.Errorf("wanted event that no storage is available but none was sent")
					return
				}
			}
		})
	}
}
