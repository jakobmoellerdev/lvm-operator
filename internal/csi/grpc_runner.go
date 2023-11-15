package csi

import (
	"context"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type gRPCServerRunner struct {
	srv            *grpc.Server
	sockFile       string
	leaderElection bool
}

var _ manager.LeaderElectionRunnable = gRPCServerRunner{}

// NewGRPCRunner creates controller-runtime's manager.Runnable for a gRPC server.
// The server will listen on UNIX domain socket at sockFile.
// If leaderElection is true, the server will run only when it is elected as leader.
func NewGRPCRunner(srv *grpc.Server, sockFile string, leaderElection bool) manager.Runnable {
	return gRPCServerRunner{srv, sockFile, leaderElection}
}

// Start implements controller-runtime's manager.Runnable.
func (r gRPCServerRunner) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("Starting gRPC server", "sockFile", r.sockFile)
	err := os.Remove(r.sockFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	lis, err := net.Listen("unix", r.sockFile)
	if err != nil {
		return err
	}

	go func() {
		if err := r.srv.Serve(lis); err != nil {
			logger.Error(err, "gRPC server error")
		}
	}()
	<-ctx.Done()
	logger.Info("Stopping gRPC server")

	start := time.Now()
	end := make(chan any, 1)
	go func() {
		r.srv.GracefulStop()
		end <- nil
	}()
	select {
	case <-end:
		logger.Info("Stopped gRPC server gracefully", "duration", time.Since(start))
	case <-time.After(10 * time.Second):
		r.srv.Stop()
		logger.Info("Stopped gRPC server forcibly", "duration", time.Since(start))
	}
	return nil
}

// NeedLeaderElection implements controller-runtime's manager.LeaderElectionRunnable.
func (r gRPCServerRunner) NeedLeaderElection() bool {
	return r.leaderElection
}
