package runner

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"hash/fnv"
	"k8s.io/client-go/util/retry"
	"os"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type lvmlockdRunner struct {
	nodeName string
}

var _ manager.Runnable = &lvmlockdRunner{}

func NewLvmlockdRunner(nodeName string) manager.Runnable {
	return &lvmlockdRunner{
		nodeName: nodeName,
	}
}

// Start implements controller-runtime's manager.Runnable.
func (l *lvmlockdRunner) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("Starting lvmlockd runnable")

	errs, ctx := errgroup.WithContext(ctx)

	errs.Go(func() error {
		hostID := hashToNumber(l.nodeName, 2000)
		args := []string{"--daemon-debug", "--gl-type=sanlock", fmt.Sprintf("--host-id=%d", hostID)}

		err := retry.OnError(retry.DefaultRetry, func(err error) bool { return true }, func() error {
			cmd := exec.Command("lvmlockd", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run lvmlockd daemon: %w", err)
			}
			return nil
		})

		return err
	})

	return errs.Wait()
}

func hashToNumber(nodeName string, maxNumber int) int {
	hash := fnv.New32a()
	hash.Write([]byte(nodeName))
	hashSum := hash.Sum32()
	return int(hashSum)%maxNumber + 1
}
