package driver

import (
	"fmt"

	esv1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/deployment"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/version"
	"github.com/elastic/cloud-on-k8s/pkg/controller/elasticsearch/label"
	"github.com/elastic/cloud-on-k8s/pkg/controller/elasticsearch/services"
	"github.com/elastic/cloud-on-k8s/pkg/controller/elasticsearch/user"
	"github.com/elastic/cloud-on-k8s/pkg/utils/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ExporterDeploymentName(esName string) string {
	return esv1.ExporterDeployment(esName)
}

func newExporterDeployment(es esv1.Elasticsearch) (appsv1.Deployment, error) {
	nsn := k8s.ExtractNamespacedName(&es)
	nm := ExporterDeploymentName(es.Name)
	ver, err := version.Parse(es.Spec.Version)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	lbs := label.NewExporterDeploymentLabels(nsn, nm, ver)
	sel := label.NewLabelSelectorForExporterDeployment(nsn, nm)
	d := deployment.New(deployment.Params{
		Name:            nm,
		Namespace:       es.Namespace,
		Selector:        sel,
		Labels:          lbs,
		PodTemplateSpec: newExporterPodTemplateSpec(es, lbs),
		Replicas:        1,
		Strategy:        appsv1.DeploymentStrategy{},
	})

	return d, nil
}

func esURI(es esv1.Elasticsearch) string {
	res := fmt.Sprint("--es.uri=", es.Spec.HTTP.Protocol(), "://", user.ExporterUserName, ":$(ES_PASSWORD)@", services.InternalServiceName(es.Name), ":9200")
	return res
}

func newExporterPodTemplateSpec(es esv1.Elasticsearch, lbs map[string]string) corev1.PodTemplateSpec {
	imgUrl := "quay.io/prometheuscommunity/elasticsearch-exporter:v1.3.0"
	p := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "elasticsearch-exporter",
			Namespace: es.Namespace,
			Labels:    lbs,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "elasticsearch-exporter",
					Image: imgUrl,
					Env: []corev1.EnvVar{
						{
							Name: "ES_PASSWORD",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: esv1.ExporterUserSecret(es.Name),
									},
									Key: user.ExporterUserName,
								},
							},
						},
					},
					Command: []string{
						"elasticsearch_exporter",
						"--log.format=logfmt",
						"--log.level=info",
						esURI(es),
						"--es.all",
						"--es.indices",
						"--es.shards",
						"--es.snapshots",
						"--es.timeout=30s",
						"--es.ssl-skip-verify",
						"--web.listen-address=:9108",
						"--web.telemetry-path=/metrics",
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 9108,
						},
					},
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/healthz",
								Port: intstr.IntOrString{
									Type:   intstr.String,
									StrVal: "http",
								},
							},
						},
						InitialDelaySeconds: 5,
						TimeoutSeconds:      5,
						PeriodSeconds:       5,
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/healthz",
								Port: intstr.IntOrString{
									Type:   intstr.String,
									StrVal: "http",
								},
							},
						},
						InitialDelaySeconds: 1,
						TimeoutSeconds:      5,
						PeriodSeconds:       5,
					},
					Lifecycle: &corev1.Lifecycle{
						PreStop: &corev1.LifecycleHandler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"/bin/ash", "-c", "sleep 20",
								},
							},
						},
					},
				},
			},
		},
	}
	return p
}
