package vgmanager

import (
	"fmt"
	"testing"

	"github.com/go-logr/logr/testr"
	lvmv1alpha1 "github.com/openshift/lvm-operator/api/v1alpha1"
	"github.com/openshift/lvm-operator/pkg/internal"
	mockExec "github.com/openshift/lvm-operator/pkg/internal/test"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
)

var mockLvsOutputThinPoolValid = `
{
      "report": [
          {
              "lv": [
                  {"lv_name":"thin-pool-1", "vg_name":"vg1", "lv_attr":"twi-a-tz--", "lv_size":"26.96g", "pool_lv":"", "origin":"", "data_percent":"0.00", "metadata_percent":"10.52", "move_pv":"", "mirror_log":"", "copy_percent":"", "convert_lv":""}
              ]
          }
      ]
  }
`

var mockLvsOutputThinPoolHighMetadataUse = `
{
      "report": [
          {
              "lv": [
                  {"lv_name":"thin-pool-1", "vg_name":"vg1", "lv_attr":"twi-a-tz--", "lv_size":"26.96g", "pool_lv":"", "origin":"", "data_percent":"0.00", "metadata_percent":"98.52", "move_pv":"", "mirror_log":"", "copy_percent":"", "convert_lv":""}
              ]
          }
      ]
  }
`
var mockLvsOutputThinPoolSuspended = `
{
      "report": [
          {
              "lv": [
                  {"lv_name":"thin-pool-1", "vg_name":"vg1", "lv_attr":"twi-s-tz--", "lv_size":"26.96g", "pool_lv":"", "origin":"", "data_percent":"0.00", "metadata_percent":"10.52", "move_pv":"", "mirror_log":"", "copy_percent":"", "convert_lv":""}
              ]
          }
      ]
  }
`

var mockLvsOutputRAID = `
{
      "report": [
          {
              "lv": [
                  {"lv_name":"thin-pool-1", "vg_name":"vg1", "lv_attr":"rwi-a-tz--", "lv_size":"26.96g", "pool_lv":"", "origin":"", "data_percent":"0.00", "metadata_percent":"10.52", "move_pv":"", "mirror_log":"", "copy_percent":"", "convert_lv":""}
              ]
          }
      ]
  }
`

func TestVGReconciler_validateLVs(t *testing.T) {
	type fields struct {
		executor internal.Executor
	}
	type args struct {
		volumeGroup *lvmv1alpha1.LVMVolumeGroup
	}

	lvsCommandForVG1 := []string{
		"lvs",
		"-S",
		"vgname=vg1",
		"--units",
		"g",
		"--reportformat",
		"json",
	}

	mockExecutorForLVSOutput := func(output string) internal.Executor {
		return &mockExec.MockExecutor{
			MockExecuteCommandWithOutputAsHost: func(command string, args ...string) (string, error) {
				if !slices.Equal(args, lvsCommandForVG1) {
					return "", fmt.Errorf("invalid args %q", args)
				}
				return output, nil
			},
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid LV",
			fields: fields{
				executor: mockExecutorForLVSOutput(mockLvsOutputThinPoolValid),
			},
			args: args{volumeGroup: &lvmv1alpha1.LVMVolumeGroup{
				ObjectMeta: metav1.ObjectMeta{Name: "vg1", Namespace: "default"},
				Spec: lvmv1alpha1.LVMVolumeGroupSpec{ThinPoolConfig: &lvmv1alpha1.ThinPoolConfig{
					Name: "thin-pool-1",
				}},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid LV due to Type not being Thin Pool",
			fields: fields{
				executor: mockExecutorForLVSOutput(mockLvsOutputRAID),
			},
			args: args{volumeGroup: &lvmv1alpha1.LVMVolumeGroup{
				ObjectMeta: metav1.ObjectMeta{Name: "vg1", Namespace: "default"},
				Spec: lvmv1alpha1.LVMVolumeGroupSpec{ThinPoolConfig: &lvmv1alpha1.ThinPoolConfig{
					Name: "thin-pool-1",
				}},
			}},
			wantErr: assert.Error,
		},
		{
			name: "Invalid LV due to high metadata percentage",
			fields: fields{
				executor: mockExecutorForLVSOutput(mockLvsOutputThinPoolHighMetadataUse),
			},
			args: args{volumeGroup: &lvmv1alpha1.LVMVolumeGroup{
				ObjectMeta: metav1.ObjectMeta{Name: "vg1", Namespace: "default"},
				Spec: lvmv1alpha1.LVMVolumeGroupSpec{ThinPoolConfig: &lvmv1alpha1.ThinPoolConfig{
					Name: "thin-pool-1",
				}},
			}},
			wantErr: assert.Error,
		},
		{
			name: "Invalid LV due to suspended instead of active state",
			fields: fields{
				executor: mockExecutorForLVSOutput(mockLvsOutputThinPoolSuspended),
			},
			args: args{volumeGroup: &lvmv1alpha1.LVMVolumeGroup{
				ObjectMeta: metav1.ObjectMeta{Name: "vg1", Namespace: "default"},
				Spec: lvmv1alpha1.LVMVolumeGroupSpec{ThinPoolConfig: &lvmv1alpha1.ThinPoolConfig{
					Name: "thin-pool-1",
				}},
			}},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &VGReconciler{
				Log:      testr.New(t),
				executor: tt.fields.executor,
			}
			tt.wantErr(t, r.validateLVs(tt.args.volumeGroup), fmt.Sprintf("validateLVs(%v)", tt.args.volumeGroup))
		})
	}
}
