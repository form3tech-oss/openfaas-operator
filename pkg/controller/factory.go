package controller

import (
	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	"github.com/openfaas/faas-netes/k8s"
	"github.com/openfaas/faas-provider/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
)

// FunctionFactory wraps faas-netes factory
type FunctionFactory struct {
	Factory k8s.FunctionFactory
}

func NewFunctionFactory(clientset kubernetes.Interface, config k8s.DeploymentConfig) FunctionFactory {
	return FunctionFactory{
		k8s.FunctionFactory{
			Client: clientset,
			Config: config,
		},
	}
}

func functionToFunctionRequest(in *faasv1.Function) types.FunctionDeployment {
	env := make(map[string]string)
	if in.Spec.Environment != nil {
		env = *in.Spec.Environment
	}
	lim, req := functionToFunctionResources(in)
	return types.FunctionDeployment{
		Annotations:            in.Spec.Annotations,
		Service:                in.Name,
		Labels:                 &in.Labels,
		Constraints:            in.Spec.Constraints,
		EnvProcess:             in.Spec.Handler,
		EnvVars:                env,
		Image:                  in.Spec.Image,
		Limits:                 lim,
		Requests:               req,
		ReadOnlyRootFilesystem: true, // Force the usage of a read-only root filesystem.
	}
}

func functionToFunctionResources(in *faasv1.Function) (l *types.FunctionResources, r *types.FunctionResources) {
	if in.Spec.Limits != nil {
		l = &types.FunctionResources{
			Memory: in.Spec.Limits.Memory,
			CPU:    in.Spec.Limits.CPU,
		}
	}
	if in.Spec.Requests != nil {
		r = &types.FunctionResources{
			Memory: in.Spec.Requests.Memory,
			CPU:    in.Spec.Requests.CPU,
		}
	}
	return
}

func (f *FunctionFactory) MakeProbes(function *faasv1.Function) (*k8s.FunctionProbes, error) {
	req := functionToFunctionRequest(function)
	return f.Factory.MakeProbes(req)
}

func (f *FunctionFactory) ConfigureReadOnlyRootFilesystem(function *faasv1.Function, deployment *appsv1.Deployment) {
	req := functionToFunctionRequest(function)
	f.Factory.ConfigureReadOnlyRootFilesystem(req, deployment)
}

func (f *FunctionFactory) ConfigureContainerUserID(deployment *appsv1.Deployment) {
	f.Factory.ConfigureContainerUserID(deployment)
}

func (f *FunctionFactory) ConfigureSecurityContext(deployment *appsv1.Deployment) {
	allowPrivilegeEscalation := false
	deployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation = &allowPrivilegeEscalation
	deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities = &corev1.Capabilities{
		Drop: []corev1.Capability{"ALL"},
	}
	runAsNonRoot := true
	deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsGroup = deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser
	deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsNonRoot = &runAsNonRoot
	deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
		FSGroup: deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser,
	}
}
