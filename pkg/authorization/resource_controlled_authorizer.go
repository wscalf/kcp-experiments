package authorization

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kcp-dev/logicalcluster/v2"

	kaudit "k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/tools/cache"

	redhat "github.com/kcp-dev/kcp/pkg/apis/redhat/v1"
	kcpinformers "github.com/kcp-dev/kcp/pkg/client/informers/externalversions"
	"github.com/kcp-dev/kcp/pkg/indexers"
)

const (
	ResourceControlledAuditPrefix   = "resourcecontrolled.authorization.kcp.dev/"
	ResourceControlledAuditDecision = ResourceControlledAuditPrefix + "decision"
	ResourceControlledAuditReason   = ResourceControlledAuditPrefix + "reason"
)

type ResourceControlledAuthorizer struct {
	configIndexer cache.Indexer
}

func NewResourceControlledAuthorizer(informer kcpinformers.SharedInformerFactory) authorizer.Authorizer {
	return &ResourceControlledAuthorizer{
		configIndexer: informer.Stable().V1().AuthzConfigs().Informer().GetIndexer(),
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

	objs, err := a.configIndexer.ByIndex(indexers.ByLogicalCluster, clusterName.String())
	if err == nil {
		for _, obj := range objs {
			config := obj.(*redhat.AuthzConfig)
			if config.Spec.OrgId == org {
				return authorizer.DecisionAllow, fmt.Sprintf("allowed by organization id: %s", org), nil
			}
		}
	} else {
		print("Error retrieving authorization configs: " + err.Error())
	}

	return authorizer.DecisionNoOpinion, "not matched", nil
}

func (a *ResourceControlledAuthorizer) authorizeBindingOperation(clusterName logicalcluster.Name, attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	if attr.GetVerb() == "create" {
		println(fmt.Sprintf("Running binding-side authz to %s %s in %s workspace.", attr.GetVerb(), attr.GetResource(), clusterName.String()))
	}

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
