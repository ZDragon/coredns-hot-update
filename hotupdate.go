// Package hotupdate is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package hotupdate

import (
	"context"
	versioned "github.com/ZDragon/coredns-hot-update/pkg/generated/clientset/versioned"
	listers "github.com/ZDragon/coredns-hot-update/pkg/generated/listers/federation/v1alpha1"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/file"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
	"sync"
	"time"

	//"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

const (
	FederationNs     = "supermesh"
	StatusProcessed  = "Processed"
	StatusNotStarted = "NotStarted"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("hotupdate")
var mu sync.Mutex

// HotUpdate Example is an example plugin to show how to write a plugin.
type HotUpdate struct {
	Next plugin.Handler
	file file.File
	mux  sync.Mutex
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (re *HotUpdate) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	re.mux.Lock()
	defer re.mux.Unlock()
	return re.file.ServeDNS(ctx, w, r)
}

// Name implements the Handler interface.
func (re *HotUpdate) Name() string { return "hotupdate" }

func (re *HotUpdate) CheckInDB(client listers.HostEntryLister, sliceDNS listers.HostEntriesSliceLister, qname string) bool {
	log.Infof("Call CheckInDB with check" + qname)

	list, err := client.HostEntries(FederationNs).List(labels.Everything())
	if err != nil {
		log.Errorf("Call CheckInDB Error %s", err)
		return false
	}

	for _, v := range list {
		if strings.ToLower(v.Spec.Host) == strings.ToLower(qname) {
			return true
		}
	}

	listSlices, err := sliceDNS.HostEntriesSlices(FederationNs).List(labels.Everything())
	if err != nil {
		log.Errorf("Call ReCalculateDB Error %s", err)
		return false
	}

	for _, v := range listSlices {
		for _, ii := range v.Spec.Items {
			if strings.ToLower(ii.Host) == strings.ToLower(qname) {
				return true
			}
		}
	}

	return false
}

func (re *HotUpdate) ReCalculateDB(cl versioned.Interface,
	singleDNS listers.HostEntryLister, sliceDNS listers.HostEntriesSliceLister, forceMode bool) {
	start := time.Now()
	log.Infof("Call ReCalculateDB")
	re.mux.Lock()

	//re.file = file.File{Zones: file.Zones{Z: make(map[string]*file.Zone), Names: []string{}}}

	list, err := singleDNS.HostEntries(FederationNs).List(labels.Everything())
	if err != nil {
		log.Errorf("Call ReCalculateDB Error %s", err)
		return
	}

	for _, v := range list {
		if v.Status.Process != StatusProcessed || forceMode {
			err := re.Add(context.TODO(), v.Name, v.Spec.Host, v.Spec.RR)
			if err != nil {
				log.Errorf("Call ReCalculateDB Add RR Error %s", err)
				return
			}
			vCopy := v.DeepCopy()
			vCopy.Status.Process = StatusProcessed
			vCopy.Status.LastUpdate = v1.NewTime(time.Now())
			status, err := cl.FederationV1alpha1().HostEntries(FederationNs).UpdateStatus(context.TODO(), vCopy, v1.UpdateOptions{})
			if err != nil {
				log.Errorf("Call ReCalculateDB Add RR Error %s", status)
				return
			}
		}
	}

	listSlices, err := sliceDNS.HostEntriesSlices(FederationNs).List(labels.Everything())
	if err != nil {
		log.Errorf("Call ReCalculateDB Error %s", err)
		return
	}

	for _, v := range listSlices {
		if v.Status.Process != StatusProcessed || forceMode {
			for _, ii := range v.Spec.Items {
				err := re.Add(context.TODO(), v.Name, ii.Host, ii.RR)
				if err != nil {
					log.Errorf("Call ReCalculateDB Add RR Error %s", err)
					return
				}
			}
			vCopy := v.DeepCopy()
			vCopy.Status.Process = StatusProcessed
			vCopy.Status.LastUpdate = v1.NewTime(time.Now())
			status, err := cl.FederationV1alpha1().HostEntriesSlices(FederationNs).UpdateStatus(context.TODO(), vCopy, v1.UpdateOptions{})
			if err != nil {
				log.Errorf("Call ReCalculateDB Add RR Error %s", status)
				return
			}
		}
	}
	re.mux.Unlock()
	log.Infof("END ReCalculateDB. Time %s", time.Since(start))
}

func (re *HotUpdate) Add(ctx context.Context, name string, host string, rr []string) error {
	log.Infof("Add dns record: %v %v %v", name, host, rr)
	qname := plugin.Host(host).NormalizeExact()[0]
	//log.Infof("Origins len: %v", len(re.file.Zones.Names))
	zone := plugin.Zones(re.file.Zones.Names).Matches(qname)
	if zone == "" {
		//log.Infof("Zone %v empty, try add qname %v", zone, qname)
		z := file.NewZone(".", "")

		for _, v := range rr {
			//log.Infof("RR %v", v)
			rr, err := dns.NewRR("$ORIGIN " + qname + "\n" + v + "\n")
			if err != nil {
				return err
			}
			rr.Header().Name = strings.ToLower(rr.Header().Name)
			if err := z.Insert(rr); err != nil {
				return err
			}
			//log.Infof("Log rr: %v", rr)
		}

		re.file.Zones.Z["."] = z
		re.file.Zones.Names = append(re.file.Zones.Names, ".")
	} else {
		//log.Infof("Zone %v found, try add qname %v", zone, qname)
		z := re.file.Zones.Z["."]
		_, result := z.Search(qname)
		if result {
			log.Infof("QNAME %v found, duplicate host entry incorrect by default", qname)
		} else {
			for _, v := range rr {
				rr, err := dns.NewRR("$ORIGIN " + qname + "\n" + v + "\n")
				if err != nil {
					return err
				}
				rr.Header().Name = strings.ToLower(rr.Header().Name)
				if err := z.Insert(rr); err != nil {
					return err
				}
				//log.Infof("Log rr: %v", rr)
			}
			re.file.Zones.Z["."] = z
		}
	}
	/*
		for s2, z := range re.file.Zones.Z {
			log.Infof("Zone %v", s2)
			for _, i2 := range z.All() {
				log.Infof("RR records %v", i2.All())
			}
		}*/

	return nil
}

// New returns a pointer to a new and intialized Records.
func New() *HotUpdate {
	re := new(HotUpdate)
	re.file = file.File{Zones: file.Zones{Z: make(map[string]*file.Zone), Names: []string{}}}
	return re
}
