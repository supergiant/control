package ssh

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAtArxGzmUffkRNy4bpITg0oicUA6itrh2RumMoydra2QqRL8i
sA6xBaPHbBAOJO/gY/h/qvr8Hnb38GFJcQQy2eENb83i2u8BVnnN2IFkgyCyYCN7
DE54bQejH0xD4qMhXdyEUOyKaOBzHHBliyIR4HmobiddJho4G0Ku3onLDm+++XNG
ZVNj0drFOE0YG+s/Zy5j/8EH3b2NgNzoE9h+jcIkZaRNuYYo5e26RCaakXeJT/Iu
vRFYtSPy3tgGQ+Q/Aj1Kv6Gjv0OFTqO7mmHN+nShThtjdLaCMM8hr72OoVkjljk1
kUBG7VAy0YlgRzxlEOUyXGKY7cOU4Nx51RIYQwIDAQABAoIBAEV2d0F+vKjBoH++
nVGjJq5zoINOsj52+sMvNmB4Q/yB/8DYUYTFlkzLvJQXua1MkzFe3brU7NLAKban
glNFQG1JZAq/z4eScNyxT9b5TRM+WTO4XLAJ0nKWYLwhi4t0TtpMywwBxwDhn+fY
AYVllqoZpf8h1tFtijoSRy960En381UuXWMGUDwzpr01G1GLX0ux4sCVMB4l1SX2
yWxtS/tJKvdLUNlEH72n+w9uz4Xpt62QXubDiibFDtomDbJrIi15umMAGcSgIvLq
2T/HoyXysIKdJNzu8jTYLBHY0oKmMN+DS047rDL6HMX3OH2GTpbnWBt8vk9Ndgwh
akMvfzkCgYEA35icpBVRDS+NUPMFbWcZRPkv8mSCPxJBL2/PHhRyjf85lOiMG7Rl
FVNA1pvATkD5QRo1SP4QiMcLbtLxX5bcevmg48ziGoWcSync5PUKnSwKMR/AAolZ
HD8pL2sjXpedNy1twwIY4XWALC+jf3h2UlBD0Xc/WcHPvpjNaPH1QiUCgYEAziKB
IuagG2uUkeBW5iQiC4eGuelTnRgCBN8CO6TWL88q0IEHtgx9c/s+rQ1AeR5VbGRu
lECDKCw6aq6swAsW4sacGYxvMyivdnSKSrKTFqDbV8ccJ+g30HoVTrcfXN1pZIpd
urEFVTi7uokINLuhVHqmY9poOkLD9ugx1G1dwEcCgYEAjWLyQetcyiq0gGh7mRdl
ajDr+alGlt1TLMzVuh6R5WprHdcCqY4jkR2I1Wu9aX46XslUwmgtSmAawaRPjvNV
TcnFy+ZFXyH3l6vMC1dLs+EiPLfn8XKqT2s8/sgPoIPcnQRz8KjF1OM4/jfNehBR
OXGZwL1X6MjeNZZn5SGCxaUCgYEAwpkGShloeppT9mbQApGH5lR6FpYzzjD07v5M
0FurrBSqOY4l5nHiGRTNtXa2L+E4CYzCa4h/iPQ/7aibAu01HL8cbG4MKEK3al9Y
km4Et68BgttANFhgIJqv9NChdfy72yNYmr805qAZcV6d9ZJQGj1zSP7NuHqBH11S
dVUN1U8CgYEAw5N6ScysYb9Jsaurcykij4mn1tvXzpDcap/Lqu/QXSUJZU1D7Cac
OOJSve1MuYQbV1LEIc15yMPsWTTik2Z98r9IL+3xdofh9yFaG1nxzi9OkN6aVMAz
dZM6MSCYh9kcT0pi2FPmY9iXba9kx4XAnf+0YB5xCz9QSMk4W5xSTBs=
-----END RSA PRIVATE KEY-----`

func TestStepRunSuccess(t *testing.T) {
	step := &Step{}

	port := "22"
	user := "root"
	host := "10.20.30.40"
	timeout := 120

	config := steps.NewConfig("", "", "", profile.Profile{})
	config.SshConfig = steps.SshConfig{
		Port:                port,
		User:                user,
		Timeout:             timeout,
		BootstrapPrivateKey: privateKey,
	}

	config.AddMaster(&node.Node{
		PrivateIp: host,
		PublicIp:  host,
	})

	config.Node = node.Node{
		PublicIp: host,
	}

	err := step.Run(context.Background(), ioutil.Discard, config)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if config.Runner == nil {
		t.Error("Runner must not be nil")
		return
	}
}

func TestStepRunError(t *testing.T) {
	step := &Step{}

	port := "22"
	user := "root"
	host := "10.20.30.40"
	timeout := 120

	config := steps.NewConfig("", "", "", profile.Profile{})
	config.Node = node.Node{
		PrivateIp: host,
	}
	config.SshConfig = steps.SshConfig{
		Port:                port,
		User:                user,
		Timeout:             timeout,
		BootstrapPrivateKey: "",
	}

	err := step.Run(context.Background(), ioutil.Discard, config)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}
}

func TestStepName(t *testing.T) {
	s := Step{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s", StepName, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := Step{}

	if len(s.Depends()) != 1 && s.Depends()[0] != "node" {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{"node"})
	}
}

func TestStep_Rollback(t *testing.T) {
	s := Step{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
