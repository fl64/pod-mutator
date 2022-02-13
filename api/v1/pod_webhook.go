package v1

import (
	"context"
	"encoding/json"
	"github.com/fl64/pod-mutator/internal/cfg"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
//var podlog = logf.Log.WithName("pod-resource")

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-core-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups=core,resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1
type PodMutator struct {
	Client  client.Client
	decoder *admission.Decoder
	cfg     *cfg.Cfg
}

func (p *PodMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	podlog := logf.FromContext(ctx).WithName("pod-resource")
	podlog.Info("Start mutator")
	pod := &corev1.Pod{}
	err := p.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	sidecar := corev1.Container{
		Name:            "sidecar-echo",
		Image:           "fl64/echo-http:latest",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 8000,
			},
		},
	}
	pod.Spec.Containers = append(pod.Spec.Containers, sidecar)
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (p *PodMutator) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}
