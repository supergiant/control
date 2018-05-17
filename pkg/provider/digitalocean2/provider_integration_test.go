package digitalocean2
//
//import (
//	"testing"
//	"github.com/supergiant/supergiant/pkg/core"
//	"github.com/supergiant/supergiant/test/fake_core"
//	"github.com/Sirupsen/logrus"
//	"github.com/supergiant/supergiant/pkg/model"
//	"github.com/stretchr/testify/require"
//)
//
//func TestProvider_CreateAndDeleteKube(t *testing.T) {
//	c := &core.Core{
//		DB:  new(fake_core.DB),
//		Log: logrus.New(),
//	}
//	p := Provider{
//		//FIXME
//		DOClient:  NewClient("be0a1efdaee4b2301e4af61d4bb2b7741312406a9d7ebf449d570d3b273e45a4"),
//		Core:      c,
//		Provisioner: nil,
//	}
//	k := &model.Kube{
//		MasterName: "qqqq",
//		DigitalOceanConfig: &model.DOKubeConfig{
//			Region: "nyc1",
//			SSHKeyFingerprint: []string{
//				"1c:b7:65:18:23:e6:9f:28:c1:64:e9:aa:ba:ec:56:39",
//			},
//		},
//		MasterNodeSize: "s-1vcpu-1gb",
//	}
//	k.Name = k.MasterName
//
//	err := p.CreateKube(k, nil)
//	require.NoError(t, err)
//	err = p.DeleteKube(k, nil)
//	require.NoError(t, err)
//}
