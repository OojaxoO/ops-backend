package kube

import (
	"net/http"
	"errors"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"ops-backend/pkg/app"
	"ops-backend/pkg/file"
	"ops-backend/pkg/e"
	"ops-backend/pkg/models"
	"ops-backend/view/kube/cluster"
)

type Kube struct {
	Config string
}

func NewClient (config string, name string) (Kube, error) {
	configPath := "/tmp/config" + name
	if len(config) == 0 {
		return Kube{}, errors.New("集群配置为空")
	} 
	if ok := file.SaveFile(configPath, config); !ok {
		return Kube{}, errors.New("配置文件保存失败")
	}
	k := Kube{Config: configPath} 
	return k, nil
}

func (this *Kube) getClient () (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", this.Config)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func GetClient (c *gin.Context) (*kubernetes.Clientset, error) {
	cluster := &cluster.Cluster{}
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(cluster); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return nil, err
	}
	if err := models.Get(cluster); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return nil, err
	}
	kube, err := NewClient(cluster.Config, cluster.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return nil, err
	}
	return kube.getClient()
} 

func List (c *gin.Context) {
	resource := c.Param("resource")
	client, err := GetClient(c)
	appG := app.Gin{C: c}
	if err != nil {
		return
	}
	labelSelector := c.DefaultQuery("labelSelector", "")
	opts := metav1.ListOptions{LabelSelector: labelSelector}
	namespace := c.DefaultQuery("namespace", "default")
	var data interface{} 
	switch resource {
	case "node":
		data, err = client.CoreV1().Nodes().List(opts)
	case "pod":
		data, err = client.CoreV1().Pods(namespace).List(opts)
	case "deployment":
		data, err = client.AppsV1().Deployments(namespace).List(opts)
	case "daemonset":
		data, err = client.AppsV1().DaemonSets(namespace).List(opts)
	case "replicaset":
		data, err = client.AppsV1().ReplicaSets(namespace).List(opts)
	case "stateful":
		data, err = client.AppsV1().StatefulSets(namespace).List(opts)
	case "namespace":
		data, err = client.CoreV1().Namespaces().List(opts)
	case "service":
		data, err = client.CoreV1().Services(namespace).List(opts)
	case "configmap":
		data, err = client.CoreV1().ConfigMaps(namespace).List(opts)
	case "pv":
		data, err = client.CoreV1().PersistentVolumes().List(opts)
	case "pvc":
		data, err = client.CoreV1().PersistentVolumeClaims(namespace).List(opts)
	case "ingress":
		data, err = client.NetworkingV1beta1().Ingresses(namespace).List(opts)
    default:
		appG.Response(http.StatusInternalServerError, e.ERROR, "不支持的资源类型")
		return
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

