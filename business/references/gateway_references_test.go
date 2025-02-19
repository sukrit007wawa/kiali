package references

import (
	"testing"

	"github.com/stretchr/testify/assert"
	networking_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/tests/data"
)

func prepareTestForGateway(gw *networking_v1alpha3.Gateway, vss []networking_v1alpha3.VirtualService) models.IstioReferences {
	gwReferences := GatewayReferences{
		Gateways:        []networking_v1alpha3.Gateway{*gw},
		VirtualServices: vss,
		WorkloadsPerNamespace: map[string]models.WorkloadList{
			"test": data.CreateWorkloadList("istio-system",
				data.CreateWorkloadListItem("istio-ingressgateway", map[string]string{"istio": "ingressgateway"})),
		},
	}
	return *gwReferences.References()[models.IstioReferenceKey{ObjectType: "gateway", Namespace: gw.Namespace, Name: gw.Name}]
}

func TestGatewayReferences(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	references := prepareTestForGateway(fakeGateway(t), fakeVirtualServices(t))

	// Check Workload references
	assert.Len(references.WorkloadReferences, 1)
	assert.Equal(references.WorkloadReferences[0].Name, "istio-ingressgateway")
	assert.Equal(references.WorkloadReferences[0].Namespace, "istio-system")

	// Check VS references
	assert.Len(references.ObjectReferences, 3)
	assert.Equal(references.ObjectReferences[0].Name, "reviews1")
	assert.Equal(references.ObjectReferences[0].Namespace, "bookinfo")
	assert.Equal(references.ObjectReferences[0].ObjectType, "virtualservice")
	assert.Equal(references.ObjectReferences[1].Name, "reviews2")
	assert.Equal(references.ObjectReferences[1].Namespace, "bookinfo")
	assert.Equal(references.ObjectReferences[1].ObjectType, "virtualservice")
	assert.Equal(references.ObjectReferences[2].Name, "reviews3")
	assert.Equal(references.ObjectReferences[2].Namespace, "bookinfo")
	assert.Equal(references.ObjectReferences[2].ObjectType, "virtualservice")
}

func TestGatewayNoWorkloadReferences(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	references := prepareTestForGateway(data.CreateEmptyGateway("reviews-empty", "bookinfo", map[string]string{"wrong": "selector"}), fakeVirtualServices(t))
	assert.Empty(references.WorkloadReferences)
}

func fakeGateway(t *testing.T) *networking_v1alpha3.Gateway {
	gwObject := data.CreateEmptyGateway("gateway", "istio-system", map[string]string{
		"istio": "ingressgateway",
	})

	return gwObject
}

func fakeVirtualServices(t *testing.T) []networking_v1alpha3.VirtualService {
	loader := yamlFixtureLoader("multiple-vs-gateways.yaml")
	err := loader.Load()
	if err != nil {
		t.Error("Error loading test data.")
	}

	return loader.FindVirtualServiceIn("bookinfo")
}
