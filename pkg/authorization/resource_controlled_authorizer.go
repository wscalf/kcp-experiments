package authorization

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kaudit "k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	kubernetesinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const (
	ResourceControlledAuditPrefix   = "resourcecontrolled.authorization.kcp.dev/"
	ResourceControlledAuditDecision = ResourceControlledAuditPrefix + "decision"
	ResourceControlledAuditReason   = ResourceControlledAuditPrefix + "reason"
)

const expectedOrgId = "123"

type ResourceControlledAuthorizer struct {
	configLister cache.GenericLister
}

func NewResourceControlledAuthorizer(informer kubernetesinformers.SharedInformerFactory) authorizer.Authorizer {
	configInformer, err := informer.ForResource(schema.GroupVersionResource{
		Group:    "stable.redhat.com",
		Version:  "v1",
		Resource: "authz-config",
	})

	if err != nil {
		return nil
	} else {
		return &ResourceControlledAuthorizer{
			configLister: configInformer.Lister(),
		}
	}
}

func (a *ResourceControlledAuthorizer) Authorize(ctx context.Context, attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	lcluster, err := genericapirequest.ClusterNameFrom(ctx)
	if err != nil {
		kaudit.AddAuditAnnotations(
			ctx,
			ResourceControlledAuditDecision, DecisionNoOpinion,
			ResourceControlledAuditReason, fmt.Sprintf("error getting cluster from request: %v", err),
		)
		return authorizer.DecisionNoOpinion, "internal error", err
	}

	switch attr.GetResource() {
	case "apiexports":
		if attr.GetVerb() != "bind" {
			return notApplicable()
		}

		return a.authorizeExportBind(lcluster, attr)
	case "apibindings":
		return a.authorizeBindingOperation(lcluster, attr)
	default:
		return notApplicable()
	}
}

func (a *ResourceControlledAuthorizer) authorizeExportBind(clusterName logicalcluster.Name, attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	user := attr.GetUser()
	org, err := getTenantId(user)
	if err != nil {
		return authorizer.DecisionNoOpinion, "internal error", err
	}

	configs, err := a.configLister.List(labels.Everything())
	for _, config := range configs {
		print(config.GetObjectKind().GroupVersionKind().Kind)
	}

	if org == expectedOrgId {
		return authorizer.DecisionAllow, fmt.Sprintf("allowed by organization id: %s", org), nil
	} else {
		return authorizer.DecisionNoOpinion, "not matched", nil
	}
}

func (a *ResourceControlledAuthorizer) authorizeBindingOperation(clusterName logicalcluster.Name, attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	println(fmt.Sprintf("Running binding-side authz to %s %s in %s workspace.", attr.GetVerb(), attr.GetResource(), clusterName.String()))
	return authorizer.DecisionNoOpinion, "not implemented", nil
}

func getTenantId(usr user.Info) (string, error) {
	for _, group := range usr.GetGroups() {
		if strings.HasPrefix(group, "org/") {
			return group[4:], nil
		}
	}
	return "", errors.New("organization group not found")
}

func notApplicable() (authorized authorizer.Decision, reason string, err error) {
	return authorizer.DecisionNoOpinion, "Not applicable", nil
}
