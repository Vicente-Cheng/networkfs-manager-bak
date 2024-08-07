package endpoint

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	ctlendpoint "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	networkfsv1 "github.com/Vicente-Cheng/networkfs-manager/pkg/apis/harvesterhci.io/v1beta1"
	ctlntefsv1 "github.com/Vicente-Cheng/networkfs-manager/pkg/generated/controllers/harvesterhci.io/v1beta1"
	"github.com/Vicente-Cheng/networkfs-manager/pkg/utils"
)

type Controller struct {
	namespace string
	nodeName  string

	EndpointCache     ctlendpoint.EndpointsCache
	Endpoints         ctlendpoint.EndpointsController
	NetworkFSCache    ctlntefsv1.NetworkFilesystemCache
	NetworkFilsystems ctlntefsv1.NetworkFilesystemController
}

const (
	netFSEndpointHandlerName = "harvester-netfs-endpoint-handler"
)

// Register register the longhorn node CRD controller
func Register(ctx context.Context, endpoint ctlendpoint.EndpointsController, netfilesystems ctlntefsv1.NetworkFilesystemController, opt *utils.Option) error {

	c := &Controller{
		namespace:         opt.Namespace,
		nodeName:          opt.NodeName,
		Endpoints:         endpoint,
		EndpointCache:     endpoint.Cache(),
		NetworkFilsystems: netfilesystems,
		NetworkFSCache:    netfilesystems.Cache(),
	}

	c.Endpoints.OnChange(ctx, netFSEndpointHandlerName, c.OnEndpointChange)
	c.Endpoints.OnRemove(ctx, netFSEndpointHandlerName, c.OnEndpointDelete)
	return nil
}

// OnChange watch the node CR on change and sync up to block device CR
func (c *Controller) OnEndpointChange(_ string, endpoint *corev1.Endpoints) (*corev1.Endpoints, error) {
	if endpoint == nil || endpoint.DeletionTimestamp != nil {
		logrus.Infof("Skip this round because the endpoint %s is deleting", endpoint.Name)
		return nil, nil
	}

	// we only care about the endpoint with name prefix "pvc-"
	if !strings.HasPrefix(endpoint.Name, "pvc-") {
		return nil, nil
	}

	logrus.Infof("Handling endpoint %s change event", endpoint.Name)
	if len(endpoint.Subsets) == 0 || len(endpoint.Subsets[0].Addresses) == 0 {
		return nil, fmt.Errorf("endpoint %s is not ready", endpoint.Name)
	}

	networkFS, err := c.NetworkFilsystems.Get(c.namespace, endpoint.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		logrus.Errorf("Failed to get networkFS %s: %v", endpoint.Name, err)
		return nil, err
	}
	if errors.IsNotFound(err) {
		logrus.Infof("Creating networkFilesystem CRD %s", endpoint.Name)
		return c.createNetworkFS(endpoint)
	}

	networkFSCpy := networkFS.DeepCopy()
	if networkFSCpy.Status.Endpoint != endpoint.Subsets[0].Addresses[0].IP {
		networkFSCpy.Status.Endpoint = endpoint.Subsets[0].Addresses[0].IP
		networkFSCpy.Status.State = networkfsv1.NetworkFSStateEnabled
	}

	if !reflect.DeepEqual(networkFS, networkFSCpy) {
		if _, err := c.NetworkFilsystems.Update(networkFSCpy); err != nil {
			logrus.Errorf("Failed to update networkFS %s: %v", networkFS.Name, err)
			return nil, err
		}
	}

	return nil, nil
}

// OnNodeDelete watch the node CR on remove and delete node related block devices
func (c *Controller) OnEndpointDelete(_ string, endpoint *corev1.Endpoints) (*corev1.Endpoints, error) {
	if endpoint == nil || endpoint.DeletionTimestamp == nil {
		return nil, nil
	}
	logrus.Infof("Handling endpoint %s delete event", endpoint.Name)
	return nil, nil
}

func (c *Controller) createNetworkFS(endpoint *corev1.Endpoints) (*corev1.Endpoints, error) {
	networkFS := &networkfsv1.NetworkFilesystem{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpoint.Name,
			Namespace: c.namespace,
		},
		Spec: networkfsv1.NetworkFSSpec{
			NetworkFSName: endpoint.Name,
		},
		Status: networkfsv1.NetworkFSStatus{
			Endpoint: endpoint.Subsets[0].Addresses[0].IP,
			State:    networkfsv1.NetworkFSStateDisabled,
		},
	}

	networkFSType := networkfsv1.NetworkFSTypeNFS
	networkFS.Status.Type = networkFSType

	if _, err := c.NetworkFilsystems.Create(networkFS); err != nil {
		logrus.Errorf("Failed to create networkFS %s: %v", networkFS.Name, err)
		return nil, err
	}

	return nil, nil
}
