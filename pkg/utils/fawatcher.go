package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	kubeletDir string = "/var/lib/kubelet/device-plugins"
)

func WatchKubelet(stop chan<- struct{}) (*fsnotify.Watcher, error) {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create fs watcher")
	}

	go func() {
		// Start listening for events.
		for {
			select {
			case event, ok := <-fswatcher.Events:
				if !ok {
					return
				}
				klog.Infof("fs event: %s %v", event.Name, event.Op.String())
				if event.Name == pluginapi.KubeletSocket && event.Op == fsnotify.Create {
					klog.Warning("inotify: kubelet.sock created, restarting.")
					stop <- struct{}{}
				}
			case err, ok := <-fswatcher.Errors:
				if !ok {
					return
				}
				klog.Errorf("fsnotify failed restarting,detail:%v", err)
			}
		}
	}()
	err = fswatcher.Add(kubeletDir)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to add kubelet socket to fs watcher")
	}
	return fswatcher, nil
}
