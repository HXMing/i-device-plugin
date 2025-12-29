# i-device-plugin

---
kubernetes device plugin demo

+ i-device-plugin会新增资源hxm.com/gopher
+ 将会扫描获取 /etc/gophers 目录下的文件作为对应的设备。
+ 将设备分配给 Pod 后，会在 Pod 中新增环境变量Gopher=$deviceId

## 构建镜像
```sh
make build-image
```

## 部署
使用daemonset来部署i-device-plugin，以便于在每个节点上安装
```sh
kubectl apply -f deploy/daemonset.yaml
```

检查Pod运行情况
```sh
kubectl -n kube-system get pods
```