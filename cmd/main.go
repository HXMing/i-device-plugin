package main

import (
	"github.com/HXMing/i-device-plugin/pkg/deviceplugin"
	"github.com/HXMing/i-device-plugin/pkg/utils"
	"k8s.io/klog/v2"
)

func main() {
	klog.Infof("device plugin starting...")
	dp := deviceplugin.NewGopherDevicePlugin()
	go dp.Run()

	if err := dp.Register(); err != nil {
		klog.Fatalf("register to kubelet failed: %v", err)
	}

	stop := make(chan struct{})
	fw, err := utils.WatchKubelet(stop)
	if err != nil {
		klog.Fatalf("start to kubelet failed: %v", err)
	}
	defer fw.Close()

	<-stop
	klog.Infof("kubelet restart, exiting")
}
