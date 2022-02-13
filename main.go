/*
Copyright 2022 fl64 <flsixtyfour@gmail.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"github.com/fl64/pod-mutator/internal/cfg"
	"github.com/fl64/pod-mutator/internal/mutator"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config.yaml", "Path to yaml config file")
	flag.Parse()
	config, err := cfg.GetCfg(configPath)
	if err != nil {
		setupLog.Error(err, "unable to get config")
		os.Exit(1)
		return
	}

	opts := zap.Options{
		Development: config.LoggerCfg.DevMode,
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	configJson, err := json.Marshal(config)
	if err != nil {
		setupLog.Error(err, "unable to marshal config")
		os.Exit(1)
		return
	}
	setupLog.Info(string(configJson))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricAddr,
		Port:                   9443,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.LeaderElect,
		LeaderElectionID:       "7cb46bcd.my.domain",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	m := &mutator.PodMutator{
		Client: mgr.GetClient(),
		Cfg:    config,
	}
	setupLog.Info("register webhook")
	mgr.GetWebhookServer().Register("/mutate-core-v1-pod", &webhook.Admission{Handler: m})
	setupLog.Info("register finished")

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
