package bosh

import (
	"fmt"
)

type Process struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type VM struct {
	AgentID            string   `json:"agent_id"`
	VMCID              string   `json:"vm_cid"`
	ResourcePool       string   `json:"resource_pool"`
	DiskCID            string   `json:"disk_cid"`
	JobName            string   `json:"job_name"`
	Index              int      `json:"index"`
	ResurrectionPaused bool     `json:"resurrection_paused"`
	JobState           string   `json:"job_state"`
	IPs                []string `json:"ips"`
	DNS                []string `json:"dns"`
	//	Vitals             []Vital   `json:"vitals"`
	Processes []Process `json:"processes"`
}

func (t Target) GetVMs() ([]VM, error) {
	var l []VM

	r, err := t.Get("/vms")
	if err != nil {
		return l, err
	}

	if r.StatusCode != 200 {
		return l, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	if err = t.InterpretJSON(r, &l); err != nil {
		return l, err
	}
	return l, nil
}
