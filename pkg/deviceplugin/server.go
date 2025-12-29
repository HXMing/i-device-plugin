package deviceplugin

import (
	"context"
	"net"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/HXMing/i-device-plugin/pkg/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type GopherDevicePlugin struct {
	server *grpc.Server
	stop   chan struct{}
	dm     *DeviceMonitor
	pluginapi.UnimplementedDevicePluginServer
}

func NewGopherDevicePlugin() *GopherDevicePlugin {
	return &GopherDevicePlugin{
		server: grpc.NewServer(grpc.EmptyServerOption{}),
		stop:   make(chan struct{}),
		dm:     NewDeviceMonitor(common.DevicePath),
	}
}

// Run start gRPC server and watcher
func (d *GopherDevicePlugin) Run() error {
	err := d.dm.List()
	if err != nil {
		klog.Fatalf("list device error: %v", err)
	}

	go func() {
		if err = d.dm.Watch(); err != nil {
			klog.Infoln("watch device error")
		}
	}()

	pluginapi.RegisterDevicePluginServer(d.server, d)

	// delete old unix socket before start
	socket := path.Join(pluginapi.DevicePluginPath, common.DeviceSocket)
	err = syscall.Unlink(socket)

	if err != nil && !os.IsNotExist(err) {
		return errors.WithMessagef(err, "delete socket: %s failed", socket)
	}

	sock, err := net.Listen("unxi", socket)
	if err != nil {
		return errors.WithMessagef(err, "listen unix %s failed", socket)
	}

	go d.server.Serve(sock)

	// wait for server to start by launching a block connection
	conn, err := connect(common.DeviceSocket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// dial establishes the gRPC communication with the registered device plugin.
func connect(socketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c, err := grpc.DialContext(ctx, socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return net.DialTimeout("unix", addr, time.Until(deadline))
			}
			return net.DialTimeout("unix", addr, common.ConnectTimeout)
		}),
	)

	if err != nil {
		return nil, err
	}
	return c, nil
}
