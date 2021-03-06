package hotupdate

import (
	"encoding/json"
	clientset "github.com/ZDragon/coredns-hot-update/pkg/generated/clientset/versioned"
	samplescheme "github.com/ZDragon/coredns-hot-update/pkg/generated/clientset/versioned/scheme"
	informers "github.com/ZDragon/coredns-hot-update/pkg/generated/informers/externalversions"
	listers "github.com/ZDragon/coredns-hot-update/pkg/generated/listers/federation/v1alpha1"
	"github.com/ZDragon/coredns-hot-update/pkg/signals"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"net/http"
	"os"
	"time"
)

const IsDev = false

// init registers this plugin.
func init() { plugin.Register("hotupdate", setup) }

// setup is the function that gets called when the config parser see the token "example". Setup is responsible
// for parsing any extra options the example plugin may have. The first token this function sees is "example".
func setup(c *caddy.Controller) error {
	c.Next() // Ignore "example" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("hotupdate", c.ArgErr())
	}

	re := New()

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		re.Next = next
		return re
	})

	config, err := getConfig()
	if err != nil {
		klog.Fatalf("Error connect to cluster %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	re.LoadAllEntries(exampleClient)
	go startKubeAPI(re, exampleClient)

	// All OK, return a nil error. very useful comment
	return nil
}

func getConfig() (*rest.Config, error) {
	if IsDev {
		return clientcmd.BuildConfigFromFlags("", "/Users/u17908803/.kube/config")
	} else {
		return rest.InClusterConfig()
	}
}

func startKubeAPI(re *HotUpdate, exampleClient *clientset.Clientset) {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))

	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(exampleClient,
		exampleInformerFactory.Federation().V1alpha1().HostEntries(),
		exampleInformerFactory.Federation().V1alpha1().HostEntriesSlices(),
		re)

	go startRestAPI(re, controller.singleDNSLister, controller.sliceDNSLister)

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	exampleInformerFactory.Start(stopCh)

	err := controller.Run(1, stopCh)
	if err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

	klog.Info("KubeAPI Controller started")
}

func startRestAPI(re *HotUpdate, client listers.HostEntryLister, lister listers.HostEntriesSliceLister) {
	port := os.Getenv("PORT") //???????????????? ???????? ???? ?????????? .env; ???? ???? ?????????????? ????????, ?????????????? ?????? ?????????????????? ???????????????????????? ???????????? ???????????????????????? ???????????? ????????????
	if port == "" {
		port = "8000" //localhost
	}

	klog.Info("RestAPI starting with port " + port)

	http.HandleFunc("/api/dns/check", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			klog.Fatal("Fatal error with parse form " + err.Error())
			return
		}

		qname := r.PostFormValue("host")
		klog.Info("Get req for check with host " + qname)
		sendResponse(w, err, re.CheckInDB(client, lister, qname))
	})
	err := http.ListenAndServe(":"+port, nil) //?????????????????? ????????????????????, ???????????????? localhost:8000/api

	if err != nil {
		klog.Fatal(err.Error())
	}

	klog.Info("RestAPI started")
}

func sendResponse(w http.ResponseWriter, err error, response bool) {
	errParse := json.NewEncoder(w).Encode(map[string]bool{"dns_record_found": response})
	if errParse != nil {
		klog.Fatal("Fatal error with answer " + err.Error())
	}
}
