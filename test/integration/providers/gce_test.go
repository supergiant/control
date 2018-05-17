package providers

import (
	"context"
	"fmt"
	"github.com/supergiant/supergiant/pkg/util"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGCE(t *testing.T) {
	projectId := os.Getenv("GCE_PROJECT_ID")
	region := os.Getenv("GCE_REGION")
	zone := os.Getenv("GCE_ZONE")
	pubKey := os.Getenv("GCE_PUB_KEY")
	privateKey := os.Getenv("GCE_PRIVATE_KEY")
	email := os.Getenv("GCE_EMAIL")
	tokenUri := os.Getenv("GCE_TOKEN_URI")

	// NOTE(stgleb): When pass through env variables it gets modified
	privateKey = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCt0O2s16BRimGQ\n98yQQR9kKhtpVpGT9u2BSm4l/bc+v4ay5az8elIZhu4jfTba9mcu7Unt/H4H2qhW\nd/n0QIjvIyy+FT8shg5m091qkB4MLQ8Jz4uYx25B61c91LWcT5tFCf8SXfQy2+ST\nUJ6txNujipLUV5urk9u3enjdIQfVQg2ECwCi9rbkeIWqNQiqm3kjb2rVXBFRpt9r\nSQ9zyXR9wu9JxgHAbrB4C9j12rhz9cI1AOaZJ+Q2vZq3pG6kPaNVpCD+rasKcrnu\nuDIVS1BCHttmv3p5encxjxcPvu87EPRMj5KuOPXsGuukJxWoU/waj3qy8RGUr0VD\nwAt8CDvrAgMBAAECggEAC9alaYFM3QuPdkqPtz+IPiU+tV+nuKviefeHewpFzH30\n/wL0/kPtNZSNbFi0lMNK5ypS6mT+cdXVuKtL7fPl9Quwmh25A3T5+XedUWBm0NjS\nI7lBaR15h/8vHk9PiJ7vagtdQR3/FOZfMh1qlqvyyGJsJQ72flWAFdtp7KhkRS+x\nb7OMkage9AHAdV259DFYsZ4ORU1JxvqHli7fVXpSIHsAgEYxJeqjqj/PQeSmdlen\ntqtBFrGy7P30E/y8coGfz2icvXipp2dSeQb1Gpp1e/AdoFWMbvhzqYZWdPckhvcT\nIsjdN+ULD3FWmF4Cwd06KEqEttuNOH6sB5jBMlZyIQKBgQDavk9vHHptyxAD3Rcz\nz/sCQT0/p5wQFu5UZhT2mXN0V4i9jh3qqOWNIDMu8sPVwB3vKeqnjWUDXpj5bXMu\nIfxr6oXxSHsyQdCBnd3OildlsU8Eh6MJ6gbIj88rIvcUkO+V0NM4snGK45w6m74f\nASDVVQcHRNUgeEYF55STLQ7UEQKBgQDLa7Dir4dvJ42HmyvnDu8GIbfTAhyCS1Dz\nvgqoyUPN8IT2kXaSPU35jyAiQgdmNgGSymBqj1a7/BcvwWZHEj4GXIuf3kNYYEYL\nGjq2pt66sTg05yZCnFzWwqtyr0PXBq5VmhMuo6b6hN5QbkooiZGgxInTuLpX69cV\nqzt2wTOcOwKBgDBHqv6qOXd8R1ei99kOwac4wQ0IsJB4jzf/pAbdzbbTDzJPaNj/\nWFMy1Tk6ifDmy3SbOtiqg64ftgHvn2mCRNWI2PFtfwuTrTK+plNNA4dFgFxOl7S7\ne63O1/n8aK6YYtkdU1GDST5PiI8DCw6K0DVl4/w9vBDDmyj4eTmWy1wRAoGASdH3\n7Bu083KQGuEF6qDxvvDni8ydWe9JHlsd9Sis0YRyTCR3uhRDQshc6fG6S65XndSR\nbro7yJZwN6Vgn3QQTDCzfr2jBORTJt5K5lPiSi/b7N7hdJTX4BvfKgxOey7yfyAd\ny/QZuZoUL24GvXVHAuev+MR140gz0qpENxFf0FcCgYEA0WFxDgpwFtjDz9bL9LO7\n85TOd7FS/xy2oBcdtVdfnPojSKvF5P1gLspQhDlr6CvM+6d8lUpDRXst9dROHBvR\nOIiM7c/nkq8EndNydS/HJYlgDN/R+G/viM4+uxKKIkcDk/ZwyWQwwBNfM4jxcAC6\ngO0wwbXXVvN9y57NyUSNjzw=\n-----END PRIVATE KEY-----\n"

	k8sVersions := []string{"1.8.7", "1.7.7", "1.6.7", "1.5.7"}
	srv, err := newServer()

	if err != nil {
		t.Errorf("Unexpected error while creating client and server %v", err)
		return
	}

	requestor := createAdmin(srv.Core)
	client := srv.Core.APIClient("token", requestor.APIToken)

	go srv.Start()
	defer srv.Stop()

	cloudAccount, err := createCloudAccount(client,
		map[string]string{
			"project_id":   projectId,
			"region":       region,
			"zone":         zone,
			"private_key":  privateKey,
			"client_email": email,
			"token_uri":    tokenUri},
		"gce")

	if err != nil {
		t.Errorf("Unexpected error while creating cloud account %v", err)
		return
	}

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Run test GCE-%s", k8sVersion), func(t *testing.T) {

			kube, err := createKubeGCE(client, cloudAccount,
				fmt.Sprintf("test-%s", strings.ToLower(util.RandomString(5))),
				zone,
				pubKey,
				k8sVersion)

			if err != nil {
				t.Errorf("error creating the kube %v", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
			defer cancel()

			err = util.WaitFor(ctx, "Wait for cluster to start", time.Second*30, func() (bool, error) {
				err := client.Kubes.Get(kube.ID, kube)

				if err != nil {
					return false, err
				}

				return kube.Ready, nil
			})

			// Clean up afterwards
			err = client.Kubes.Delete(kube.ID, kube)

			if err != nil {
				t.Errorf("Error while deleting the kube %s", kube.ID)
			}
		})
	}
}
