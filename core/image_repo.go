package core

import (
	"path"

	"github.com/supergiant/guber"
)

type ImageRepoResource struct {
	c *Client
}

type ImageRepo struct {
	r    *ImageRepoResource
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ImageRepoList struct {
	Items []*ImageRepo `json:"items"`
}

// Model (Entity)
//==============================================================================

// Util-------------------------------------------------------------------------
// AsKubeSecret returns a Kubernetes Secret definition
func (m *ImageRepo) AsKubeSecret() *guber.Secret {
	return &guber.Secret{
		Metadata: &guber.Metadata{
			Name: m.Name,
		},
		Type: "kubernetes.io/dockercfg",
		Data: map[string]string{
			".dockercfg": m.Key,
		},
	}
}

func (m *ImageRepo) AsKubeImagePullSecret() *guber.ImagePullSecret {
	return &guber.ImagePullSecret{
		Name: m.Name,
	}
}

// Resource (Collection)
//==============================================================================
func (r *ImageRepoResource) EtcdKey(modelName string) string {
	return path.Join("/image_repos/dockerhub", modelName)
}

// InitializeModel is a part of Resource interface
func (r *ImageRepoResource) InitializeModel(m Model) {
	model := m.(*ImageRepo)
	model.r = r
}

func (r *ImageRepoResource) List() (*ImageRepoList, error) {
	list := new(ImageRepoList)
	err := r.c.DB.List(r, list)
	return list, err
}

func (r *ImageRepoResource) Create(m *ImageRepo) (*ImageRepo, error) {
	err := r.c.DB.Create(r, m.Name, m)
	return m, err
}

func (r *ImageRepoResource) Get(name string) (*ImageRepo, error) {
	m := new(ImageRepo)
	err := r.c.DB.Get(r, name, m)
	return m, err
}

// No update for ImageRepo (yet)

func (r *ImageRepoResource) Delete(name string) error {
	return r.c.DB.Delete(r, name)
}
