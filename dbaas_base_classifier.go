package dbaasbase

import (
	"context"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	. "github.com/netcracker/qubership-core-lib-go/v3/const"
	"github.com/netcracker/qubership-core-lib-go/v3/security"
	"github.com/netcracker/qubership-core-lib-go/v3/serviceloader"
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
	tenantProvider := serviceloader.MustLoad[security.TenantProvider]()
	tenantId, err := tenantProvider.GetTenantId(ctx)
	if err != nil {
		logger.PanicC(ctx, "Can't create tenant database, tenantId is absent")
	}
	classifier["tenantId"] = tenantId
	return classifier
}
