package vmwareconf

// Use govc source code to get esxi name

// # Reference

//  https://github.com/vmware/govmomi/tree/master/govc

// govc is a vSphere CLI built on top of govmomi.

// The CLI is designed to be a user friendly CLI alternative to the GUI and well suited for automation tasks.
// It also acts as a [test harness](test) for the govmomi APIs and provides working examples of how to use the APIs.

// ## License

// govc is available under the [Apache 2 license](../LICENSE).

// ## Name

// Pronounced "go-v-c", short for "Go(lang) vCenter CLI".

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type info struct {
	WaitForIP       bool
	General         bool
	ExtraConfig     bool
	Resources       bool
	ToolsConfigInfo bool
}

type infoResult struct {
	VirtualMachines []mo.VirtualMachine
	objects         []*object.VirtualMachine
	entities        map[types.ManagedObjectReference]string
	cmd             *info
}

// collectReferences builds a unique set of MORs to the set of VirtualMachines,
// so we can collect properties in a single call for each reference type {host,datastore,network}.
func (r *infoResult) collectReferences(pc *property.Collector, ctx context.Context) error {
	// MOR -> Name map
	r.entities = make(map[types.ManagedObjectReference]string)

	var host []mo.HostSystem
	var network []mo.Network
	var opaque []mo.OpaqueNetwork
	var dvp []mo.DistributedVirtualPortgroup
	var datastore []mo.Datastore
	// Table to drive inflating refs to their mo.* counterparts (dest)
	// and save() the Name to r.entities w/o using reflection here.
	// Note that we cannot use a []mo.ManagedEntity here, since mo.Network has its own 'Name' field,
	// the mo.Network.ManagedEntity.Name field will not be set.
	vrefs := map[string]*struct {
		dest interface{}
		refs []types.ManagedObjectReference
		save func()
	}{
		"HostSystem": {
			&host, nil, func() {
				for _, e := range host {
					r.entities[e.Reference()] = e.Name
				}
			},
		},
		"Network": {
			&network, nil, func() {
				for _, e := range network {
					r.entities[e.Reference()] = e.Name
				}
			},
		},
		"OpaqueNetwork": {
			&opaque, nil, func() {
				for _, e := range opaque {
					r.entities[e.Reference()] = e.Name
				}
			},
		},
		"DistributedVirtualPortgroup": {
			&dvp, nil, func() {
				for _, e := range dvp {
					r.entities[e.Reference()] = e.Name
				}
			},
		},
		"Datastore": {
			&datastore, nil, func() {
				for _, e := range datastore {
					r.entities[e.Reference()] = e.Name
				}
			},
		},
	}

	xrefs := make(map[types.ManagedObjectReference]bool)
	// Add MOR to vrefs[kind].refs avoiding any duplicates.
	addRef := func(refs ...types.ManagedObjectReference) {
		for _, ref := range refs {
			if _, exists := xrefs[ref]; exists {
				return
			}
			xrefs[ref] = true
			vref := vrefs[ref.Type]
			vref.refs = append(vref.refs, ref)
		}
	}

	for _, vm := range r.VirtualMachines {
		if r.cmd.General {
			if ref := vm.Summary.Runtime.Host; ref != nil {
				addRef(*ref)
			}
		}

		if r.cmd.Resources {
			addRef(vm.Datastore...)
			addRef(vm.Network...)
		}
	}

	for _, vref := range vrefs {
		if vref.refs == nil {
			continue
		}
		err := pc.Retrieve(ctx, vref.refs, []string{"name"}, vref.dest)
		if err != nil {
			return err
		}
		vref.save()
	}

	return nil
}
