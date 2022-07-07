package cloudscale

import (
	cloudscalev1 "github.com/vshn/appcat-service-s3/apis/cloudscale/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

// SetupController adds a controller that reconciles cloudscalev1.ObjectsUser managed resources.
func SetupController(mgr ctrl.Manager) error {
	name := strings.ToLower(cloudscalev1.ObjectsUserGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&cloudscalev1.ObjectsUser{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{})).
		Complete(&PostgresStandaloneReconciler{
			client: mgr.GetClient(),
		})
}