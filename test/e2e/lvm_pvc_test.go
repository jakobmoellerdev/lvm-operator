/*
Copyright © 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/openshift/lvm-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	snapapi "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8sv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	PodImageForPVCTests   = "public.ecr.aws/docker/library/busybox:1.36"
	DevicePathForPVCTests = "/dev/xda"
	MountPathForPVCTests  = "/test1"
	VolumeNameForPVCTests = "vol1"
)

var (
	PodCommandForPVCTests = []string{"sh", "-c", "tail -f /dev/null"}
)

func pvcTest() {
	var cluster *v1alpha1.LVMCluster
	BeforeAll(func(ctx SpecContext) {
		cluster = GetDefaultTestLVMClusterTemplate()
		CreateResource(ctx, cluster)
		VerifyLVMSSetup(ctx, cluster)
		DeferCleanup(func(ctx SpecContext) {
			DeleteResource(ctx, cluster)
		})
	})

	for _, mode := range []k8sv1.PersistentVolumeMode{
		k8sv1.PersistentVolumeBlock,
		k8sv1.PersistentVolumeFilesystem,
	} {
		Context(string(mode), func() {
			Context("static", func() {
				pvcTestsForMode(mode)
			})
			Context("ephemeral", func() {
				ephemeralPVCTestsForMode(mode)
			})
		})
	}
}

func pvcTestsForMode(volumeMode k8sv1.PersistentVolumeMode) {
	var contentMode ContentMode
	switch volumeMode {
	case k8sv1.PersistentVolumeBlock:
		contentMode = ContentModeBlock
	case k8sv1.PersistentVolumeFilesystem:
		contentMode = ContentModeFile
	}

	pvc := generatePVC(volumeMode)
	pod := generatePodConsumingPVC(pvc)

	clonePVC := generatePVCCloneFromPVC(pvc)
	clonePod := generatePodConsumingPVC(clonePVC)

	snapshot := generateVolumeSnapshot(pvc, snapshotClass)
	snapshotPVC := generatePVCFromSnapshot(volumeMode, snapshot)
	snapshotPod := generatePodConsumingPVC(snapshotPVC)

	AfterAll(parallelDelete([]client.Object{
		snapshotPod, snapshotPVC, snapshot,
		clonePod, clonePVC,
		pod, pvc,
	}))

	expectedData := "TESTDATA"
	It("PVC and Pod", func(ctx SpecContext) {
		CreateResource(ctx, pvc)
		CreateResource(ctx, pod)
		validatePodIsRunning(ctx, client.ObjectKeyFromObject(pod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(pvc))

		Expect(contentTester.WriteDataInPod(ctx, pod, expectedData, contentMode)).To(Succeed())
		validatePodData(ctx, pod, expectedData, contentMode)
	})

	It("Cloning", func(ctx SpecContext) {
		CreateResource(ctx, clonePVC)
		CreateResource(ctx, clonePod)

		validatePodIsRunning(ctx, client.ObjectKeyFromObject(clonePod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(clonePVC))
		validatePodData(ctx, clonePod, expectedData, contentMode)
	})

	It("Snapshots", func(ctx SpecContext) {
		CreateVolumeSnapshotFromPVCOrSkipIfUnsupported(ctx, snapshot)
		validateSnapshotReadyToUse(ctx, client.ObjectKeyFromObject(snapshot))

		CreateResource(ctx, snapshotPVC)
		CreateResource(ctx, snapshotPod)

		validatePodIsRunning(ctx, client.ObjectKeyFromObject(snapshotPod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(snapshotPVC))
		validatePodData(ctx, snapshotPod, expectedData, contentMode)
	})
}

func ephemeralPVCTestsForMode(volumeMode k8sv1.PersistentVolumeMode) {
	var contentMode ContentMode
	switch volumeMode {
	case k8sv1.PersistentVolumeBlock:
		contentMode = ContentModeBlock
	case k8sv1.PersistentVolumeFilesystem:
		contentMode = ContentModeFile
	}

	pod := generatePodWithEphemeralVolume(volumeMode)

	// recreates locally what will be creatad as an ephemeral volume
	pvc := &k8sv1.PersistentVolumeClaim{}
	pvc.SetName(fmt.Sprintf("%s-%s", pod.GetName(), pod.Spec.Volumes[0].Name))
	pvc.SetNamespace(pod.GetNamespace())
	pvc.Spec.VolumeMode = &volumeMode

	clonePVC := generatePVCCloneFromPVC(pvc)
	clonePod := generatePodConsumingPVC(clonePVC)

	snapshot := generateVolumeSnapshot(pvc, snapshotClass)
	snapshotPVC := generatePVCFromSnapshot(volumeMode, snapshot)
	snapshotPod := generatePodConsumingPVC(snapshotPVC)

	AfterAll(parallelDelete([]client.Object{
		snapshotPod, snapshotPVC, snapshot,
		clonePod, clonePVC,
		pod,
	}))

	expectedData := "TESTDATA"
	It("Pod (with ephemeral Volume)", func(ctx SpecContext) {
		CreateResource(ctx, pod)
		validatePodIsRunning(ctx, client.ObjectKeyFromObject(pod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(pvc))

		Expect(contentTester.WriteDataInPod(ctx, pod, expectedData, contentMode)).To(Succeed())
		validatePodData(ctx, pod, expectedData, contentMode)
	})

	It("PVC Cloning", func(ctx SpecContext) {
		CreateResource(ctx, clonePVC)
		CreateResource(ctx, clonePod)

		validatePodIsRunning(ctx, client.ObjectKeyFromObject(clonePod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(clonePVC))
		validatePodData(ctx, clonePod, expectedData, contentMode)
	})

	It("Snapshots", func(ctx SpecContext) {
		CreateVolumeSnapshotFromPVCOrSkipIfUnsupported(ctx, snapshot)
		validateSnapshotReadyToUse(ctx, client.ObjectKeyFromObject(snapshot))

		CreateResource(ctx, snapshotPVC)
		CreateResource(ctx, snapshotPod)

		validatePodIsRunning(ctx, client.ObjectKeyFromObject(snapshotPod))
		validatePVCIsBound(ctx, client.ObjectKeyFromObject(snapshotPVC))
		validatePodData(ctx, snapshotPod, expectedData, contentMode)
	})
}

func generatePVCSpec(mode k8sv1.PersistentVolumeMode) k8sv1.PersistentVolumeClaimSpec {
	return k8sv1.PersistentVolumeClaimSpec{
		VolumeMode:  ptr.To(mode),
		AccessModes: []k8sv1.PersistentVolumeAccessMode{k8sv1.ReadWriteOnce},
		Resources: k8sv1.ResourceRequirements{Requests: map[k8sv1.ResourceName]resource.Quantity{
			k8sv1.ResourceStorage: *resource.NewScaledQuantity(1, resource.Giga),
		}},
		StorageClassName: ptr.To(storageClassName),
	}
}

func generatePVC(mode k8sv1.PersistentVolumeMode) *k8sv1.PersistentVolumeClaim {
	return &k8sv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(string(mode)),
			Namespace: testNamespace,
		},
		Spec: generatePVCSpec(mode),
	}
}

func generatePVCCloneFromPVC(pvc *k8sv1.PersistentVolumeClaim) *k8sv1.PersistentVolumeClaim {
	clone := generatePVC(*pvc.Spec.VolumeMode)
	gvk, _ := apiutil.GVKForObject(pvc, crClient.Scheme())
	clone.SetName(fmt.Sprintf("%s-clone", pvc.GetName()))
	clone.Spec.DataSource = &k8sv1.TypedLocalObjectReference{
		Kind: gvk.Kind,
		Name: pvc.GetName(),
	}
	return clone
}

func generatePVCFromSnapshot(mode k8sv1.PersistentVolumeMode, snapshot *snapapi.VolumeSnapshot) *k8sv1.PersistentVolumeClaim {
	pvc := &k8sv1.PersistentVolumeClaim{}
	pvc.SetName(snapshot.GetName())
	pvc.SetNamespace(snapshot.GetNamespace())
	pvc.Spec = generatePVCSpec(mode)
	gvk, _ := apiutil.GVKForObject(snapshot, crClient.Scheme())
	pvc.Spec.DataSource = &k8sv1.TypedLocalObjectReference{
		Kind:     gvk.Kind,
		APIGroup: ptr.To(gvk.Group),
		Name:     snapshot.GetName(),
	}
	return pvc
}

func generatePodConsumingPVC(pvc *k8sv1.PersistentVolumeClaim) *k8sv1.Pod {
	container := k8sv1.Container{
		Name:    "pause",
		Command: []string{"sh", "-c", "tail -f /dev/null"},
		Image:   PodImageForPVCTests,
		SecurityContext: &k8sv1.SecurityContext{
			RunAsNonRoot: ptr.To(true),
			SeccompProfile: &k8sv1.SeccompProfile{
				Type: k8sv1.SeccompProfileTypeRuntimeDefault,
			},
			Capabilities: &k8sv1.Capabilities{Drop: []k8sv1.Capability{"ALL"}},
		},
	}

	switch *pvc.Spec.VolumeMode {
	case k8sv1.PersistentVolumeBlock:
		container.VolumeDevices = []k8sv1.VolumeDevice{{Name: VolumeNameForPVCTests, DevicePath: DevicePathForPVCTests}}
	case k8sv1.PersistentVolumeFilesystem:
		container.VolumeMounts = []k8sv1.VolumeMount{{Name: VolumeNameForPVCTests, MountPath: MountPathForPVCTests}}
	}

	return &k8sv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       fmt.Sprintf("%s-consumer", pvc.GetName()),
			Namespace:                  testNamespace,
			DeletionGracePeriodSeconds: ptr.To(int64(1)),
		},
		Spec: k8sv1.PodSpec{
			Containers: []k8sv1.Container{container},
			Volumes: []k8sv1.Volume{{
				Name: VolumeNameForPVCTests,
				VolumeSource: k8sv1.VolumeSource{
					PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{ClaimName: pvc.GetName()},
				},
			}},
		},
	}
}

func generatePodWithEphemeralVolume(mode k8sv1.PersistentVolumeMode) *k8sv1.Pod {
	container := k8sv1.Container{
		Name:    "pause",
		Command: PodCommandForPVCTests,
		Image:   PodImageForPVCTests,
		SecurityContext: &k8sv1.SecurityContext{
			RunAsNonRoot: ptr.To(true),
			SeccompProfile: &k8sv1.SeccompProfile{
				Type: k8sv1.SeccompProfileTypeRuntimeDefault,
			},
			Capabilities: &k8sv1.Capabilities{Drop: []k8sv1.Capability{"ALL"}},
		},
	}

	switch mode {
	case k8sv1.PersistentVolumeBlock:
		container.VolumeDevices = []k8sv1.VolumeDevice{{Name: VolumeNameForPVCTests, DevicePath: DevicePathForPVCTests}}
	case k8sv1.PersistentVolumeFilesystem:
		container.VolumeMounts = []k8sv1.VolumeMount{{Name: VolumeNameForPVCTests, MountPath: MountPathForPVCTests}}
	}

	return &k8sv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       fmt.Sprintf("%s-ephemeral", strings.ToLower(string(mode))),
			Namespace:                  testNamespace,
			DeletionGracePeriodSeconds: ptr.To(int64(1)),
		},
		Spec: k8sv1.PodSpec{
			Containers: []k8sv1.Container{container},
			Volumes: []k8sv1.Volume{{
				Name: VolumeNameForPVCTests,
				VolumeSource: k8sv1.VolumeSource{
					Ephemeral: &k8sv1.EphemeralVolumeSource{VolumeClaimTemplate: &k8sv1.PersistentVolumeClaimTemplate{
						Spec: generatePVCSpec(mode)},
					},
				},
			}},
		},
	}
}

func generateVolumeSnapshot(pvc *k8sv1.PersistentVolumeClaim, snapshotClassName string) *snapapi.VolumeSnapshot {
	return &snapapi.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-snapshot", pvc.GetName()),
			Namespace: testNamespace,
		},
		Spec: snapapi.VolumeSnapshotSpec{
			Source:                  snapapi.VolumeSnapshotSource{PersistentVolumeClaimName: ptr.To(pvc.GetName())},
			VolumeSnapshotClassName: ptr.To(snapshotClassName),
		},
	}
}

func CreateVolumeSnapshotFromPVCOrSkipIfUnsupported(ctx context.Context, snapshot *snapapi.VolumeSnapshot) {
	GinkgoHelper()
	By(fmt.Sprintf("Creating VolumeSnapshot %q", snapshot.GetName()))
	err := crClient.Create(ctx, snapshot)
	if discovery.IsGroupDiscoveryFailedError(errors.Unwrap(err)) {
		Skip("Skipping Testing of VolumeSnapshot Operations due to lack of volume snapshot support")
	}
	Expect(err).ToNot(HaveOccurred(), "PVC should be created successfully")
}
