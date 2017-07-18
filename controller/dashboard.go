package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"gopkg.in/urfave/cli.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	Dashboard struct {
		ProjectConfigPath string
		listen            string
		templates         string
		demo              bool
	}

	// KubeClient is our configured client to access kubernetes
	KubeClient struct {
		namespace  string
		clientset  *kubernetes.Clientset
		kubeconfig clientcmd.ClientConfig
		restconfig *rest.Config
	}

	Deployment struct {
		Name    string
		Exists  bool
		Alive   string
		Ingress []Ingress
		Version []string
		K8s     apps.Deployment
	}

	Ingress struct {
		Url    string
		Alive  bool
		Status string
	}
)

var (
	DashboardCommand = cli.Command{
		Name:   "dashboard",
		Usage:  "Run a HTTP based dashboarde",
		Action: (&Dashboard{}).Server,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen",
				Value: ":8080",
				Usage: "Listen Address",
			},
			cli.StringFlag{
				Name:  "templates",
				Value: "templates/dashboard",
				Usage: "Dashboard Templates (dashboard.html + static/)",
			},
			cli.BoolFlag{
				Name:  "demo",
				Usage: "Demo Dashboard",
			},
		},
	}
)

// KubeClientFromConfig loads a new KubeClient from the usual configuration
// (KUBECONFIG env param / selfconfigured in kubernetes)
func KubeClientFromConfig() (*KubeClient, error) {
	var client = new(KubeClient)
	var err error

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	configOverrides := &clientcmd.ConfigOverrides{}

	client.kubeconfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	client.restconfig, err = client.kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client.clientset, err = kubernetes.NewForConfig(client.restconfig)
	if err != nil {
		return nil, err
	}

	client.namespace, _, err = client.kubeconfig.Namespace()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func checkAlive(wg *sync.WaitGroup, d *Deployment) {
	for i, ing := range d.Ingress {
		http.DefaultClient.Timeout = 2 * time.Second
		r, err := http.Get("https://" + ing.Url)
		if err != nil {
			d.Ingress[i].Status = err.Error()
		} else {
			if r.StatusCode < 500 {
				d.Ingress[i].Alive = true
			}
			d.Ingress[i].Status = r.Status
		}
	}
	wg.Done()
}

func (d *Dashboard) load() ([]Deployment, error) {
	project := loadProject(d.ProjectConfigPath)

	var deployments *apps.DeploymentList
	var ingresses *extensions.IngressList

	if d.demo {
		deployments = &apps.DeploymentList{
			Items: []apps.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "flamingo",
					},
					Spec: apps.DeploymentSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{Image: "docker.aoe.com/flamingo:v1.0.0"},
								},
							},
						},
					},
					Status: apps.DeploymentStatus{
						AvailableReplicas:  3,
						Replicas:           5,
						ObservedGeneration: 132,
						Conditions: []apps.DeploymentCondition{
							{Status: v1.ConditionTrue, Type: "TestCondition", Message: "Test Condition is feeling good!"},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "akeneo",
					},
					Spec: apps.DeploymentSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{Image: "docker.aoe.com/akeneo:v1.2.3"},
								},
							},
						},
					},
					Status: apps.DeploymentStatus{
						AvailableReplicas:  1,
						Replicas:           1,
						ObservedGeneration: 32,
						Conditions: []apps.DeploymentCondition{
							{Status: v1.ConditionTrue, Type: "TestCondition", Message: "Test Condition is feeling good!"},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keycloak",
					},
					Spec: apps.DeploymentSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{Image: "docker.aoe.com/keycloak:v1.0.0"},
									{Image: "docker.aoe.com/keycloak-support:v1.0.0"},
								},
							},
						},
					},
					Status: apps.DeploymentStatus{
						AvailableReplicas:  2,
						Replicas:           2,
						ObservedGeneration: 12,
						Conditions: []apps.DeploymentCondition{
							{Status: v1.ConditionTrue, Type: "TestCondition", Message: "Test Condition is feeling good!"},
						},
					},
				},
			},
		}

		ingresses = &extensions.IngressList{
			Items: []extensions.Ingress{
				{
					Spec: extensions.IngressSpec{
						Rules: []extensions.IngressRule{
							{
								Host: "google.com",
								IngressRuleValue: extensions.IngressRuleValue{
									HTTP: &extensions.HTTPIngressRuleValue{
										Paths: []extensions.HTTPIngressPath{
											{Backend: extensions.IngressBackend{ServiceName: "flamingo"}, Path: "/"},
										},
									},
								},
							},
						},
					},
				},
				{
					Spec: extensions.IngressSpec{
						Rules: []extensions.IngressRule{
							{
								Host: "google.com",
								IngressRuleValue: extensions.IngressRuleValue{
									HTTP: &extensions.HTTPIngressRuleValue{
										Paths: []extensions.HTTPIngressPath{
											{Backend: extensions.IngressBackend{ServiceName: "akeneo"}, Path: "/akeneo"},
										},
									},
								},
							},
						},
					},
				},
				{
					Spec: extensions.IngressSpec{
						Rules: []extensions.IngressRule{
							{
								Host: "keycloak.bla",
								IngressRuleValue: extensions.IngressRuleValue{
									HTTP: &extensions.HTTPIngressRuleValue{
										Paths: []extensions.HTTPIngressPath{
											{Backend: extensions.IngressBackend{ServiceName: "keycloak"}, Path: "/blabla"},
										},
									},
								},
							},
						},
					},
				},
				{
					Spec: extensions.IngressSpec{
						Rules: []extensions.IngressRule{
							{
								Host: "keycloak.om3",
								IngressRuleValue: extensions.IngressRuleValue{
									HTTP: &extensions.HTTPIngressRuleValue{
										Paths: []extensions.HTTPIngressPath{
											{Backend: extensions.IngressBackend{ServiceName: "keycloak"}, Path: "/"},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	} else {
		client, err := KubeClientFromConfig()
		if err != nil {
			return nil, err
		}

		deploymentClient := client.clientset.AppsV1beta1().Deployments(client.namespace)
		deployments, err = deploymentClient.List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		ingressClient := client.clientset.ExtensionsV1beta1().Ingresses(client.namespace)
		ingresses, err = ingressClient.List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	}

	deploymentIndex := make(map[string]apps.Deployment, len(deployments.Items))
	for _, deployment := range deployments.Items {
		deploymentIndex[deployment.Name] = deployment
	}

	ingressIndex := make(map[string][]Ingress)
	for _, ingress := range ingresses.Items {
		for _, rule := range ingress.Spec.Rules {
			for _, p := range rule.HTTP.Paths {
				name := p.Backend.ServiceName
				ingressIndex[name] = append(ingressIndex[name], Ingress{Url: rule.Host + p.Path})
			}
		}
	}

	var deploymentlist []Deployment
	var wg = new(sync.WaitGroup)

	for _, application := range project.Applications {
		name := application.Name
		if d, ok := application.Properties["deployment"]; !ok || d != "kubernetes" {
			continue
		}
		if n, ok := application.Properties["kubernetes-name"]; ok && n != "" {
			name = n
		}
		deployment, exists := deploymentIndex[name]

		d := Deployment{
			Name:   name,
			Exists: exists,
			Alive:  "0",
		}

		if exists {
			d.K8s = deployment

			var v []string
			for _, c := range deployment.Spec.Template.Spec.Containers {
				v = append(v, c.Image)
			}
			d.Version = v

			if len(ingressIndex[name]) > 0 {
				d.Ingress = ingressIndex[name]
				wg.Add(1)
				go checkAlive(wg, &d)
			}
		}

		deploymentlist = append(deploymentlist, d)
	}

	wg.Wait()

	return deploymentlist, nil
}

func e(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("content-type", "text/plain")
	fmt.Fprintf(rw, "%+v", err)
}

func (d *Dashboard) handler(rw http.ResponseWriter, r *http.Request) {
	deployments, err := d.load()

	if err != nil {
		e(rw, err)
		return
	}

	tpl, err := template.ParseFiles(path.Join(d.templates, "dashboard.html"))
	if err != nil {
		e(rw, err)
		return
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, deployments)
	if err != nil {
		e(rw, err)
		return
	}

	rw.Header().Set("content-type", "text/html")
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, buf)
}

// AnalyzeAction controller action
func (d *Dashboard) Server(context *cli.Context) error {
	d.ProjectConfigPath = context.GlobalString("config")
	d.demo = context.Bool("demo")
	d.listen = context.String("listen")
	d.templates = context.String("templates")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(d.templates, "static")))))
	http.HandleFunc("/", d.handler)

	log.Println("Listening on http://" + d.listen + "/")
	return http.ListenAndServe(d.listen, nil)
}
