package kubeslice

import (
	"context"

	dnsCache "bitbucket.org/realtimeai/kubeslice-dns/plugin/kubeslice/cache"
	"bitbucket.org/realtimeai/kubeslice-dns/plugin/kubeslice/slice"
	meshv1beta1 "bitbucket.org/realtimeai/kubeslice-operator/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReplicaSetReconciler is a simple ControllerManagedBy example implementation.
type ServiceImportReconciler struct {
	client.Client
	EndpointsCache dnsCache.EndpointsCache
}

// Watch the ServiceImport changes and adjust dns cache accordingly
func (r *ServiceImportReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	si := &meshv1beta1.ServiceImport{}
	err := r.Get(ctx, req.NamespacedName, si)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("got si")

	eps := []slice.Endpoint{}

	for _, ep := range si.Status.Endpoints {
		endpoint := slice.Endpoint{
			Host: ep.DNSName,
			IP:   ep.IP,
		}
		eps = append(eps, endpoint)
	}

	r.EndpointsCache.Put(si.Name, si.Spec.Slice, si.Namespace, eps)

	log.Info(r.EndpointsCache.GetAll())

	return reconcile.Result{}, nil
}

func (r *ServiceImportReconciler) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}