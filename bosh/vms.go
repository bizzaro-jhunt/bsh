package bosh

import (
	"encoding/json"
	"fmt"
	"time"
)

type Process struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Uptime struct {
		Seconds uint64 `json:"secs"`
	} `json:"uptime"`

	Memory struct {
		KB      uint64  `json:"kb"`
		Percent float32 `json:"percent"`
	} `json:"mem"`
	CPU struct {
		Total float32 `json:"total"`
	} `json:"cpu"`
}

type DiskUsage struct {
	InodePercent string `json:"inode_percent"`
	Percent      string `json:"percent"`
}

type MemUsage struct {
	KB      string `json:"kb"`
	Percent string `json:"percent"`
}

type VM struct {
	ID                 string   `json:"id"`
	AgentID            string   `json:"agent_id"`
	CID                string   `json:"vm_cid"`
	ResourcePool       string   `json:"resource_pool"`
	Type               string   `json:"vm_type"`
	AZ                 string   `json:"az"`
	DiskCID            string   `json:"disk_cid"`
	JobName            string   `json:"job_name"`
	Index              int      `json:"index"`
	ResurrectionPaused bool     `json:"resurrection_paused"`
	Ignore             bool     `json:"ignore"`
	JobState           string   `json:"job_state"`
	IPs                []string `json:"ips"`
	DNS                []string `json:"dns"`

	Vitals struct {
		CPU struct {
			Sys  string `json:"sys"`
			User string `json:"user"`
			Wait string `json:"wait"`
		} `json:"cpu"`

		Load []string `json:"load"`

		Disk struct {
			Ephemeral  DiskUsage `json:"ephemeral"`
			Persistent DiskUsage `json:"persistent"`
			System     DiskUsage `json:"system"`
		} `json:"disk"`

		Memory MemUsage `json:"mem"`
		Swap   MemUsage `json:"swap"`
	} `json:"vitals"`

	Processes []Process `json:"processes"`
}

func (vm VM) FullName() string {
	return fmt.Sprintf("%s/%d", vm.JobName, vm.Index)
}

func (t Target) GetVMsFor(deployment string) ([]VM, error) {
	var l []VM
	var tasks []Task

	r, err := t.Get(fmt.Sprintf("/deployments/%s/vms?format=full", deployment))
	if err != nil {
		return l, err
	}

	if r.StatusCode != 200 {
		return l, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	jsons, err := t.InterpretJSONList(r)
	if err != nil {
		return l, err
	}
	for _, b := range jsons {
		task := Task{}
		err = json.Unmarshal(b, &task)
		if err != nil {
			return l, err
		}
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		_, err := t.WaitTask(task.ID, 500*time.Millisecond)
		if err != nil {
			fmt.Printf("ERR: @R{%s}\n", err)
			continue
		}

		res, err := t.Get(fmt.Sprintf("/tasks/%d/output?type=result", task.ID))
		if err != nil {
			fmt.Printf("ERR: @R{%s}\n", err)
			continue
		}
		if r.StatusCode != 200 {
			fmt.Printf("BOSH API returned @R{%s}\n", r.Status)
			continue
		}

		records, err := t.InterpretJSONList(res)
		if err != nil {
			fmt.Printf("ERR: @R{%s}\n", err)
			continue
		}

		for _, record := range records {
			vm := VM{}
			err = json.Unmarshal(record, &vm)
			if err != nil {
				fmt.Printf("ERR: @R{%s}\n", err)
				continue
			}

			l = append(l, vm)
		}
	}

	return l, nil
}

func (v VM) ResurrectionState() string {
	if v.ResurrectionPaused {
		return "resurrection: off"
	}
	return "resurrection: on"
}

func (v VM) IgnoredState() string {
	if v.Ignore {
		return "ignored"
	}
	return "not ignored"
}
