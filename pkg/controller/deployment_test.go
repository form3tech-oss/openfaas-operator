package controller

import (
	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	"github.com/openfaas/faas-netes/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"reflect"
	"testing"
)

func Test_newDeployment(t *testing.T) {
	function := &faasv1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubesec",
		},
		Spec: faasv1.FunctionSpec{
			Name:  "kubesec",
			Image: "docker.io/kubesec/kubesec",
			Annotations: &map[string]string{
				"com.openfaas.serviceaccount":           "kubesec",
				"com.openfaas.health.http.initialDelay": "2m",
				"com.openfaas.health.http.path":         "/healthz",
			},
			ReadOnlyRootFilesystem: true,
			EnvFrom: []corev1.EnvFromSource{
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "foo",
						},
					},
				},
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "bar",
						},
					},
				},
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "baz",
						},
					},
				},
			},
		},
	}

	factory := NewFunctionFactory(fake.NewSimpleClientset(),
		k8s.DeploymentConfig{
			HTTPProbe:      false,
			SetNonRootUser: true,
			LivenessProbe: &k8s.ProbeConfig{
				PeriodSeconds:       1,
				TimeoutSeconds:      3,
				InitialDelaySeconds: 0,
			},
			ReadinessProbe: &k8s.ProbeConfig{
				PeriodSeconds:       1,
				TimeoutSeconds:      3,
				InitialDelaySeconds: 0,
			},
		})

	secrets := map[string]*corev1.Secret{}

	deployment := newDeployment(function, secrets, factory)

	if deployment.Spec.Template.Spec.ServiceAccountName != "kubesec" {
		t.Errorf("ServiceAccountName should be %s", "kubesec")
		t.Fail()
	}

	if deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Path != "/healthz" {
		t.Errorf("Readiness probe should have HTTPGet handler set to %s", "/healthz")
		t.Fail()
	}

	if deployment.Spec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds != 120 {
		t.Errorf("Liveness probe should have initial delay seconds set to %s", "2m")
		t.Fail()
	}

	if !*(deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem) {
		t.Errorf("ReadOnlyRootFilesystem should be true")
		t.Fail()
	}

	if *(deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser) != k8s.SecurityContextUserID {
		t.Errorf("RunAsUser should be %v", k8s.SecurityContextUserID)
		t.Fail()
	}

	if !reflect.DeepEqual(deployment.Spec.Template.Spec.Containers[0].EnvFrom, function.Spec.EnvFrom) {
		t.Errorf("EnvFrom should be %+#v", function.Spec.EnvFrom)
		t.Fail()
	}
}
