package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"supergiant/core/model"
)

type JobStorage struct {
	client *Client
}

func parseJobID(key string) string {
	strs := strings.Split(key, "/")
	return strs[len(strs)-1]
}

func (store *JobStorage) Create(e *model.Job) (*model.Job, error) {
	value, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	resp, err := store.client.CreateInOrder("/jobs", string(value))
	if err != nil {
		return nil, err
	}

	e.ID = parseJobID(resp.Node.Key)

	return e, nil
}

func (store *JobStorage) List() ([]*model.Job, error) {
	resp, err := store.client.GetInOrder("/jobs")
	if err != nil {
		return nil, err
	}

	apps := make([]*model.Job, 0)

	for _, node := range resp.Node.Nodes {
		e := new(model.Job)
		e.ID = parseJobID(node.Key)
		if err := json.Unmarshal([]byte(node.Value), e); err != nil {
			return nil, err
		}
		apps = append(apps, e)
	}
	return apps, nil
}

// func (store *JobStorage) Get(id int) (*model.Job, error) {
// 	// TODO repeated, move to method
// 	key := fmt.Sprintf("/jobs/%d", id)
// 	resp, err := store.client.Get(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	value := resp.Node.Value
//
// 	e := new(model.Job)
// 	e.ID = id
// 	if err := json.Unmarshal([]byte(value), e); err != nil {
// 		return nil, err
// 	}
// 	return e, nil
// }

func (store *JobStorage) Update(id string, e *model.Job) (*model.Job, error) {
	e.ID = id

	key := fmt.Sprintf("/jobs/%s", id)
	value, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Update(key, string(value))
	return e, err
}

func (store *JobStorage) Delete(id string) error {
	// TODO repeated
	_, err := store.client.Delete(fmt.Sprintf("/jobs/%s", id))
	return err
}
