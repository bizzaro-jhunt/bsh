package bosh

import (
	"encoding/json"
	"fmt"
	"time"
)

type Process struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type VM struct {
	AgentID            string   `json:"agent_id"`
	CID                string   `json:"vm_cid"`
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

	AZ   string `json:"az"`
	Type string `json:"vm_type"`
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
