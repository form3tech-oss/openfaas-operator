package controller

import (
	"testing"

	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

func Test_makeAnnotations_NoKeys(t *testing.T) {
	annotationVal := `{"name":"","image":"","replicas":null,"handler":"","annotations":null,"labels":null,"environment":null,"envFrom":null,"constraints":null,"secrets":null,"limits":null,"requests":null}`

	spec := faasv1.Function{
		Spec: faasv1.FunctionSpec{},
	}

	annotations := makeAnnotations(&spec)

	if _, ok := annotations["prometheus.io.scrape"]; !ok {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + " to be added")
		t.Fail()
	}
	if val, _ := annotations["prometheus.io.scrape"]; val != "false" {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + ` to equal "false"`)
		t.Fail()
	}

	if _, ok := annotations[annotationFunctionSpec]; !ok {
		t.Errorf("wanted annotation " + annotationFunctionSpec)
		t.Fail()
	}

	if val, _ := annotations[annotationFunctionSpec]; val != annotationVal {
		t.Errorf("Annotation " + annotationFunctionSpec + "\nwant: '" + annotationVal + "'\ngot: '" + val + "'")
		t.Fail()
	}
}

func Test_makeAnnotations_WithKeyAndValue(t *testing.T) {
	annotationVal := `{"name":"","image":"","replicas":null,"handler":"","annotations":{"key":"value","key2":"value2"},"labels":null,"environment":null,"envFrom":null,"constraints":null,"secrets":null,"limits":null,"requests":null}`

	spec := faasv1.Function{
		Spec: faasv1.FunctionSpec{
			Annotations: &map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		},
	}

	annotations := makeAnnotations(&spec)

	if _, ok := annotations["prometheus.io.scrape"]; !ok {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + " to be added")
		t.Fail()
	}
	if val, _ := annotations["prometheus.io.scrape"]; val != "false" {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + ` to equal "false"`)
		t.Fail()
	}

	if _, ok := annotations[annotationFunctionSpec]; !ok {
		t.Errorf("wanted annotation " + annotationFunctionSpec)
		t.Fail()
	}

	if val, _ := annotations[annotationFunctionSpec]; val != annotationVal {
		t.Errorf("Annotation " + annotationFunctionSpec + "\nwant: '" + annotationVal + "'\ngot: '" + val + "'")
		t.Fail()
	}
}

func Test_makeAnnotationsDoesNotModifyOriginalSpec(t *testing.T) {
	specAnnotations := map[string]string{
		"test.foo": "bar",
	}
	function := &faasv1.Function{
		Spec: faasv1.FunctionSpec{
			Name:        "testfunc",
			Annotations: &specAnnotations,
		},
	}

	expectedAnnotations := map[string]string{
		"prometheus.io.scrape": "false",
		"test.foo":             "bar",
		annotationFunctionSpec: `{"name":"testfunc","image":"","replicas":null,"handler":"","annotations":{"test.foo":"bar"},"labels":null,"environment":null,"envFrom":null,"constraints":null,"secrets":null,"limits":null,"requests":null}`,
	}

	makeAnnotations(function)
	annotations := makeAnnotations(function)

	if len(specAnnotations) != 1 {
		t.Errorf("length of original spec annotations has changed, expected 1, got %d", len(specAnnotations))
	}

	if specAnnotations["test.foo"] != "bar" {
		t.Errorf("original spec annotation has changed")
	}

	for name, expectedValue := range expectedAnnotations {
		actualValue := annotations[name]
		if actualValue != expectedValue {
			t.Fatalf("incorrect annotation for '%s': \nexpected %s,\ngot %s", name, expectedValue, actualValue)
		}
	}
}
