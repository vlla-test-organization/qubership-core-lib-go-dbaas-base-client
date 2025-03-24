package dbaasbase

import (
	"context"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/netcracker/qubership-core-lib-go/v3/serviceloader"
	. "github.com/netcracker/qubership-core-lib-go/v3/const"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/baseproviders/tenant"
)

func BaseServiceClassifier(ctx context.Context) map[string]interface{} {
	classifier := make(map[string]interface{})
	classifier["microserviceName"] = configloader.GetKoanf().MustString(MicroserviceNameProperty)
	classifier["namespace"] = configloader.GetKoanf().MustString(NamespaceProperty)
	classifier["scope"] = "service"
	return classifier
}

func BaseTenantClassifier(ctx context.Context) map[string]interface{} {
	classifier := make(map[string]interface{})
	classifier["microserviceName"] = configloader.GetKoanf().MustString(MicroserviceNameProperty)
	classifier["namespace"] = configloader.GetKoanf().MustString(NamespaceProperty)
	classifier["scope"] = "tenant"
	tenantProvider := serviceloader.MustLoad[tenant.TenantProviderI]()
	tenantId := tenantProvider.GetTenantId(ctx)
	if tenantId == "-" {
		logger.PanicC(ctx, "Can't create tenant database, tenantId is absent")
	}
	classifier["tenantId"] = tenantId
	return classifier
}
