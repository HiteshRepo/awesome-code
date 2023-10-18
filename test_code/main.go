package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/catinello/base62"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"sigs.k8s.io/yaml"

	"hash"
	"hash/crc32"
	"hash/fnv"
)

type Hardware struct {
	Hostname    string `json:"hostname"`
	BmcIP       string `json:"bmc_ip"`
	BmcUserName string `json:"bmc_username"`
	BmcPassword string `json:"bmc_password"`
	MAC         string `json:"mac"`
	IpAddress   string `json:"ip_address"`
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
	Nameservers string `json:"nameservers"`
	Labels      string `json:"labels"`
	Disk        string `json:"disk"`
}

var NonMandatoryFields = []string{"bmc_ip", "bmc_username", "bmc_password"}

const ActionsPath = "template.tasks.0.actions"

func getActionsFromTemplate(tmpJson json.RawMessage) (gjson.Result, error) {
	actions := gjson.Get(string(tmpJson), ActionsPath)
	if !actions.Exists() {
		return gjson.Result{}, errors.New("actions not found")
	}
	return actions, nil
}

func getTmplJsonFromTmplStr(tmplStr string) (json.RawMessage, error) {
	var tmpJson json.RawMessage
	err := yaml.Unmarshal([]byte(tmplStr), &tmpJson)
	if err != nil {
		return nil, err
	}
	return tmpJson, nil
}

const ActionFieldName = "name"
const EnvFieldName = "environment"
const ImageFieldName = "image"

func getActionIndex(actionName string, actions gjson.Result) (int, error) {
	idx := -1

	for i, a := range actions.Array() {
		res := a.Get(ActionFieldName)
		res2 := a.Get(EnvFieldName)
		if res.Exists() && res.String() == actionName && res2.Exists() && res2.IsObject() {
			idx = i
			break
		}
	}

	if idx == -1 {
		return -1, errors.New("action not found")
	}

	return idx, nil
}

const EquinixBottleRocketHegelPath = "environment.HEGEL_URLS"
const DefaultBottleRocketHegelPath = "environment.HEGEL_URL"
const DefaultUbuntuHegelPath = "environment.CONTENTS"
const EquinixPlatformName = "equinix"
const DefaultPlatformName = "default"
const BottleRocketOsName = "bottlerocket"
const UbuntuOsName = "ubuntu"

func getHegelPathByPlatformOs(platform, os string, actionIndex int) string {
	if platform == EquinixPlatformName && os == BottleRocketOsName {
		return fmt.Sprintf("%s.%d.%s", ActionsPath, actionIndex, EquinixBottleRocketHegelPath)
	}

	if platform == DefaultPlatformName && os == BottleRocketOsName {
		return fmt.Sprintf("%s.%d.%s", ActionsPath, actionIndex, DefaultBottleRocketHegelPath)
	}

	if platform == DefaultPlatformName && os == UbuntuOsName {
		return fmt.Sprintf("%s.%d.%s", ActionsPath, actionIndex, DefaultUbuntuHegelPath)
	}

	return ""
}

const EquinixBottleRocketActionName = "write-user-data"

func replaceHegelUrlEquinixTmpl() {
	tmplStr := "template:\n    global_timeout: 6000\n    id: \"\"\n    name: dp-hiteshp-eksa-cluster-m3-small-x86\n    tasks:\n    - actions:\n      - environment:\n          COMPRESSED: \"true\"\n          DEST_DISK: /dev/sda\n          IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/17/artifacts/raw/1-23/bottlerocket-v1.23.9-eks-d-1-23-5-eks-a-17-amd64.img.gz\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n        name: stream-image\n        timeout: 600\n      - environment:\n          CONTENTS: |\n            # Version is required, it will change as we support\n            # additional settings\n            version = 1\n\n            # \"eno1\" is the interface name\n            # Users may turn on dhcp4 and dhcp6 via boolean\n            [enp1s0f0np0]\n            dhcp4 = true\n            dhcp6 = false\n            # Define this interface as the \"primary\" interface\n            # for the system.  This IP is what kubelet will use\n            # as the node IP.  If none of the interfaces has\n            # \"primary\" set, we choose the first interface in\n            # the file\n            primary = true\n          DEST_DISK: /dev/sda12\n          DEST_PATH: /net.toml\n          DIRMODE: \"0755\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n        name: write-netplan\n        pid: host\n        timeout: 90\n      - environment:\n          BOOTCONFIG_CONTENTS: |\n            kernel {\n                console = \"ttyS1,115200n8\"\n            }\n          DEST_DISK: /dev/sda12\n          DEST_PATH: /bootconfig.data\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n        name: write-bootconfig\n        pid: host\n        timeout: 90\n      - environment:\n          DEST_DISK: /dev/sda12\n          DEST_PATH: /user-data.toml\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          HEGEL_URLS: http://147.75.88.178:50061,http://147.75.88.189:50061\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n        name: write-user-data\n        pid: host\n        timeout: 90\n      - image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n        name: reboot-image\n        pid: host\n        timeout: 90\n        volumes:\n        - /worker:/worker\n      name: dp-hiteshp-eksa-cluster-m3-small-x86\n      volumes:\n        - /dev:/dev\n        - /dev/console:/dev/console\n        - /lib/firmware:/lib/firmware:ro\n      worker: '{{.device_1}}'\n    version: \"0.1\""

	var tmpJson json.RawMessage
	var err error
	if tmpJson, err = getTmplJsonFromTmplStr(tmplStr); err != nil {
		return
	}

	actions, err := getActionsFromTemplate(tmpJson)
	if err != nil {
		return
	}

	idx, err := getActionIndex(EquinixBottleRocketActionName, actions)
	if err != nil {
		return
	}

	hegelPath := getHegelPathByPlatformOs(EquinixPlatformName, BottleRocketOsName, idx)
	if hegelPath == "" {
		return
	}

	existingHegelUrls := gjson.Get(string(tmpJson), hegelPath)
	finalHegelUrls := "http://new_hegel_urls:56001"
	if existingHegelUrls.Exists() {
		if strings.Contains(existingHegelUrls.String(), finalHegelUrls) {
			finalHegelUrls = existingHegelUrls.String()
		} else {
			finalHegelUrls = fmt.Sprintf("%s,%s", existingHegelUrls.String(), finalHegelUrls)
		}
	}
	value, err := sjson.Set(string(tmpJson), hegelPath, finalHegelUrls)
	if err != nil {
		return
	}

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println("yaml:", string(b))
}

const DefaultUbuntuActionName = "add-tink-cloud-init-config"

func replaceHegelUrlDefaultUbuntuTmpl() {
	tmplStr := "template:\n    global_timeout: 6000\n    id: \"\"\n    name: my-cluster-name\n    tasks:\n    - actions:\n      - environment:\n          COMPRESSED: \"true\"\n          DEST_DISK: /dev/sda\n          IMG_URL: https://my-file-server/ubuntu-v1.23.7-eks-a-12-amd64.gz\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: stream-image\n        timeout: 360\n      - environment:\n          DEST_DISK: /dev/sda2\n          DEST_PATH: /etc/netplan/config.yaml\n          STATIC_NETPLAN: true\n          DIRMODE: \"0755\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: write-netplan\n        timeout: 90\n      - environment:\n          CONTENTS: |\n            datasource:\n              Ec2:\n                metadata_urls: [<admin-machine-ip>, <tinkerbell-ip-from-cluster-config>]\n                strict_id: false\n            manage_etc_hosts: localhost\n            warnings:\n              dsid_missing_source: off            \n          DEST_DISK: /dev/sda2\n          DEST_PATH: /etc/cloud/cloud.cfg.d/10_tinkerbell.cfg\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0600\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: add-tink-cloud-init-config\n        timeout: 90\n      - environment:\n          CONTENTS: |\n            network:\n              config: disabled            \n          DEST_DISK: /dev/sda2\n          DEST_PATH: /etc/cloud/cloud.cfg.d/99-disable-network-config.cfg\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0600\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: disable-cloud-init-network-capabilities\n        timeout: 90\n      - environment:\n          CONTENTS: | \n            datasource: Ec2\n          DEST_DISK: /dev/sda2\n          DEST_PATH: /etc/cloud/ds-identify.cfg\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0600\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: add-tink-cloud-init-ds-config\n        timeout: 90\n      - environment:\n          BLOCK_DEVICE: /dev/sda2\n          FS_TYPE: ext4\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/kexec:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: kexec-image\n        pid: host\n        timeout: 90\n      name: my-cluster-name\n      volumes:\n      - /dev:/dev\n      - /dev/console:/dev/console\n      - /lib/firmware:/lib/firmware:ro\n      worker: '{{.device_1}}'\n    version: \"0.1\"\n"

	var tmpJson json.RawMessage
	var err error
	if tmpJson, err = getTmplJsonFromTmplStr(tmplStr); err != nil {
		return
	}

	actions, err := getActionsFromTemplate(tmpJson)
	if err != nil {
		return
	}

	idx, err := getActionIndex(DefaultUbuntuActionName, actions)
	if err != nil {
		return
	}

	hegelPath := getHegelPathByPlatformOs(DefaultPlatformName, UbuntuOsName, idx)
	if hegelPath == "" {
		return
	}

	content := gjson.Get(string(tmpJson), hegelPath)

	if !content.Exists() {
		return
	}

	var contentJson json.RawMessage
	err = yaml.Unmarshal([]byte(content.String()), &contentJson)
	if err != nil {
		return
	}
	datasource := gjson.Get(string(contentJson), "datasource")
	metadataUrls := gjson.Get(datasource.String(), "Ec2.metadata_urls")
	lastIdx := len(metadataUrls.Array())

	value, err := sjson.Set(datasource.String(), fmt.Sprintf("Ec2.metadata_urls.%d", lastIdx), "http://new_hegel_urls:56001")
	if err != nil {
		return
	}

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}

	value, err = sjson.Set(string(contentJson), "datasource", string(b))
	if err != nil {
		return
	}

	b, err = yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}

	value, err = sjson.Set(string(tmpJson), hegelPath, string(b))
	if err != nil {
		return
	}

	b, err = yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println("yaml:", string(b))
}

const DefaultBottleRocketActionName = "write-user-data"

func replaceHegelUrlDefaultBottleRocketTmpl() {
	tmplStr := "template:\n    global_timeout: 6000\n    id: \"\"\n    name: my-cluster-name\n    tasks:\n    - actions:\n      - environment:\n          COMPRESSED: \"true\"\n          DEST_DISK: /dev/sda\n          IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/11/artifacts/raw/1-22/bottlerocket-v1.22.10-eks-d-1-22-8-eks-a-11-amd64.img.gz\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: stream-image\n        timeout: 360\n      - environment:\n          # An example console declaration that will send all kernel output to both consoles, and systemd output to ttyS0.\n          # kernel {\n          #     console = \"tty0\", \"ttyS0,115200n8\"\n          # }\n          BOOTCONFIG_CONTENTS: |\n                        kernel {}\n          DEST_DISK: /dev/sda12\n          DEST_PATH: /bootconfig.data\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: write-bootconfig\n        timeout: 90\n      - environment:\n          CONTENTS: |\n            # Version is required, it will change as we support\n            # additional settings\n            version = 1\n            # \"eno1\" is the interface name\n            # Users may turn on dhcp4 and dhcp6 via boolean\n            [eno1]\n            dhcp4 = true\n            # Define this interface as the \"primary\" interface\n            # for the system.  This IP is what kubelet will use\n            # as the node IP.  If none of the interfaces has\n            # \"primary\" set, we choose the first interface in\n            # the file\n            primary = true            \n          DEST_DISK: /dev/sda12\n          DEST_PATH: /net.toml\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: write-netconfig\n        timeout: 90\n      - environment:\n          HEGEL_URL: http://<hegel-ip>:50061\n          DEST_DISK: /dev/sda12\n          DEST_PATH: /user-data.toml\n          DIRMODE: \"0700\"\n          FS_TYPE: ext4\n          GID: \"0\"\n          MODE: \"0644\"\n          UID: \"0\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        name: write-user-data\n        timeout: 90\n      - name: \"reboot\"\n        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15\n        timeout: 90\n        volumes:\n          - /worker:/worker\n      name: my-cluster-name\n      volumes:\n      - /dev:/dev\n      - /dev/console:/dev/console\n      - /lib/firmware:/lib/firmware:ro\n      worker: '{{.device_1}}'\n    version: \"0.1\"\n"

	var tmpJson json.RawMessage
	var err error
	if tmpJson, err = getTmplJsonFromTmplStr(tmplStr); err != nil {
		return
	}

	actions, err := getActionsFromTemplate(tmpJson)
	if err != nil {
		return
	}

	idx, err := getActionIndex(DefaultBottleRocketActionName, actions)
	if err != nil {
		return
	}

	hegelPath := getHegelPathByPlatformOs(DefaultPlatformName, BottleRocketOsName, idx)
	if hegelPath == "" {
		return
	}

	value, err := sjson.Set(string(tmpJson), hegelPath, "http://new_hegel_urls:56001")
	if err != nil {
		return
	}

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println("yaml:", string(b))
}

const DefaultUbuntuTmpl = `
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: TinkerbellTemplateConfig
metadata:
  name: my-cluster-name
spec:
  template:
    global_timeout: 6000
    id: ""
    name: my-cluster-name
    tasks:
    - actions:
      - environment:
          COMPRESSED: "true"
          DEST_DISK: /dev/sda
          IMG_URL: https://my-file-server/ubuntu-v1.23.7-eks-a-12-amd64.gz
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: stream-image
        timeout: 360
      - environment:
          DEST_DISK: /dev/sda2
          DEST_PATH: /etc/netplan/config.yaml
          STATIC_NETPLAN: true
          DIRMODE: "0755"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: write-netplan
        timeout: 90
      - environment:
          CONTENTS: |
            datasource:
              Ec2:
                metadata_urls: [<admin-machine-ip>, <tinkerbell-ip-from-cluster-config>]
                strict_id: false
            manage_etc_hosts: localhost
            warnings:
              dsid_missing_source: off            
          DEST_DISK: /dev/sda2
          DEST_PATH: /etc/cloud/cloud.cfg.d/10_tinkerbell.cfg
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0600"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: add-tink-cloud-init-config
        timeout: 90
      - environment:
          CONTENTS: |
            network:
              config: disabled            
          DEST_DISK: /dev/sda2
          DEST_PATH: /etc/cloud/cloud.cfg.d/99-disable-network-config.cfg
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0600"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: disable-cloud-init-network-capabilities
        timeout: 90
      - environment:
          CONTENTS: | 
            datasource: Ec2
          DEST_DISK: /dev/sda2
          DEST_PATH: /etc/cloud/ds-identify.cfg
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0600"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: add-tink-cloud-init-ds-config
        timeout: 90
      - environment:
          BLOCK_DEVICE: /dev/sda2
          FS_TYPE: ext4
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/kexec:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: kexec-image
        pid: host
        timeout: 90
      name: my-cluster-name
      volumes:
      - /dev:/dev
      - /dev/console:/dev/console
      - /lib/firmware:/lib/firmware:ro
      worker: '{{.device_1}}'
    version: "0.1"

`

const DefaultBottleRocketTmpl = `
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: TinkerbellTemplateConfig
metadata:
  name: my-cluster-name
spec:
  template:
    global_timeout: 6000
    id: ""
    name: my-cluster-name
    tasks:
    - actions:
      - environment:
          COMPRESSED: "true"
          DEST_DISK: /dev/sda
          IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/11/artifacts/raw/1-22/bottlerocket-v1.22.10-eks-d-1-22-8-eks-a-11-amd64.img.gz
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: stream-image
        timeout: 360
      - environment:
          # An example console declaration that will send all kernel output to both consoles, and systemd output to ttyS0.
          # kernel {
          #     console = "tty0", "ttyS0,115200n8"
          # }
          BOOTCONFIG_CONTENTS: |
                        kernel {}
          DEST_DISK: /dev/sda12
          DEST_PATH: /bootconfig.data
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: write-bootconfig
        timeout: 90
      - environment:
          CONTENTS: |
            # Version is required, it will change as we support
            # additional settings
            version = 1
            # "eno1" is the interface name
            # Users may turn on dhcp4 and dhcp6 via boolean
            [NIC_NAME]
            dhcp4 = true
            # Define this interface as the "primary" interface
            # for the system.  This IP is what kubelet will use
            # as the node IP.  If none of the interfaces has
            # "primary" set, we choose the first interface in
            # the file
            primary = true            
          DEST_DISK: /dev/sda12
          DEST_PATH: /net.toml
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: write-netconfig
        timeout: 90
      - environment:
          HEGEL_URL: http://<hegel-ip>:50061
          DEST_DISK: /dev/sda12
          DEST_PATH: /user-data.toml
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: write-user-data
        timeout: 90
      - name: "reboot"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        timeout: 90
        volumes:
          - /worker:/worker
      name: my-cluster-name
      volumes:
      - /dev:/dev
      - /dev/console:/dev/console
      - /lib/firmware:/lib/firmware:ro
      worker: '{{.device_1}}'
    version: "0.1"
`

const EquinixBottleRocketTmpl = `
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: TinkerbellTemplateConfig
metadata:
  name: cp-hiteshp-eksa-cluster-m3-small-x86
spec:
  template:
    global_timeout: 6000
    id: ""
    name: cp-hiteshp-eksa-cluster-m3-small-x86
    tasks:
    - actions:
      - environment:
          COMPRESSED: "true"
          DEST_DISK: /dev/sda
          IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/17/artifacts/raw/1-23/bottlerocket-v1.23.9-eks-d-1-23-5-eks-a-17-amd64.img.gz
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17
        name: stream-image
        timeout: 600
      - environment:
          CONTENTS: |
            # Version is required, it will change as we support
            # additional settings
            version = 1

            # "eno1" is the interface name
            # Users may turn on dhcp4 and dhcp6 via boolean
            [NIC_NAME]
            dhcp4 = true
            dhcp6 = false
            # Define this interface as the "primary" interface
            # for the system.  This IP is what kubelet will use
            # as the node IP.  If none of the interfaces has
            # "primary" set, we choose the first interface in
            # the file
            primary = true
          DEST_DISK: /dev/sda12
          DEST_PATH: /net.toml
          DIRMODE: "0755"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17
        name: write-netplan
        pid: host
        timeout: 90
      - environment:
          BOOTCONFIG_CONTENTS: |
            kernel {
                console = "ttyS1,115200n8"
            }
          DEST_DISK: /dev/sda12
          DEST_PATH: /bootconfig.data
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17
        name: write-bootconfig
        pid: host
        timeout: 90
      - environment:
          DEST_DISK: /dev/sda12
          DEST_PATH: /user-data.toml
          DIRMODE: "0700"
          FS_TYPE: ext4
          GID: "0"
          HEGEL_URLS: http://147.75.88.178:50061,http://147.75.88.189:50061
          MODE: "0644"
          UID: "0"
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17
        name: write-user-data
        pid: host
        timeout: 90
      - image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17
        name: reboot-image
        pid: host
        timeout: 90
        volumes:
        - /worker:/worker
      name: cp-hiteshp-eksa-cluster-m3-small-x86
      volumes:
        - /dev:/dev
        - /dev/console:/dev/console
        - /lib/firmware:/lib/firmware:ro
      worker: '{{.device_1}}'
    version: "0.1"
`

const BottleRocketStreamImageActionName = "stream-image"
const BottleRocketWriteBootCnfActionName = "write-bootconfig"
const BottleRocketWriteNetCnfActionName = "write-netconfig"
const BottleRocketWriteUserDataActionName = "write-user-data"
const BottleRocketRebootImageActionName = "reboot"

const ActionsPathFromRoot = "spec.template.tasks.0.actions"
const UbuntuStreamImageActionName = "stream-image"
const UbuntuWriteNetPlanActionName = "write-netplan"
const UbuntuAddTinkCloudInitCnfActionName = "add-tink-cloud-init-config"
const UbuntuDisableTinkCloudInitNetworkCapActionName = "disable-cloud-init-network-capabilities"
const UbuntuAddTinkCloudInitDsCnfActionName = "add-tink-cloud-init-ds-config"
const UbuntuKexecImageActionName = "kexec-image"

const EquinixBottleRocketStreamImageActionName = "stream-image"
const EquinixBottleRocketWriteBootCnfActionName = "write-bootconfig"
const EquinixBottleRocketWriteNetPlanActionName = "write-netplan"
const EquinixBottleRocketWriteUserDataActionName = "write-user-data"
const EquinixBottleRocketRebootImageActionName = "reboot-image"

const ImageUrlFieldName = "IMG_URL"

func generateDefaultUbuntuTmpl() {
	tmpJson, err := getTmplJsonFromTmplStr(DefaultUbuntuTmpl)
	if err != nil {
		fmt.Println("fail 1")
		return
	}

	actions := gjson.Get(string(tmpJson), ActionsPathFromRoot)
	if !actions.Exists() || !actions.IsArray() {
		fmt.Println("fail 2")
		return
	}

	var StreamImageActionIdx, WriteNetPlanActionIdx, AddTinkCloudInitCnfActionIdx, DisableTinkCloudInitNetworkCapActionIdx, AddTinkCloudInitDsCnfActionIdx, KexecImageActionIdx int
	for i, a := range actions.Array() {
		name := a.Get(ActionFieldName)
		if !name.Exists() {
			continue
		}

		if name.String() == UbuntuStreamImageActionName {
			StreamImageActionIdx = i
			continue
		}

		if name.String() == UbuntuWriteNetPlanActionName {
			WriteNetPlanActionIdx = i
			continue
		}

		if name.String() == UbuntuAddTinkCloudInitCnfActionName {
			AddTinkCloudInitCnfActionIdx = i
			continue
		}

		if name.String() == UbuntuDisableTinkCloudInitNetworkCapActionName {
			DisableTinkCloudInitNetworkCapActionIdx = i
			continue
		}

		if name.String() == UbuntuAddTinkCloudInitDsCnfActionName {
			AddTinkCloudInitDsCnfActionIdx = i
			continue
		}

		if name.String() == UbuntuKexecImageActionName {
			KexecImageActionIdx = i
		}
	}

	value, err := sjson.Set(string(tmpJson), fmt.Sprintf("%s.%d.%s.%s", ActionsPathFromRoot, StreamImageActionIdx, EnvFieldName, ImageUrlFieldName), "img_url")
	if err != nil {
		fmt.Println("fail 3")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, StreamImageActionIdx, ImageFieldName), "StreamImageAction")
	if err != nil {
		fmt.Println("fail 4")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, AddTinkCloudInitDsCnfActionIdx, ImageFieldName), "AddTinkCloudInitDsCnfAction")
	if err != nil {
		fmt.Println("fail 5")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, KexecImageActionIdx, ImageFieldName), "KexecImageAction")
	if err != nil {
		fmt.Println("fail 6")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, AddTinkCloudInitCnfActionIdx, ImageFieldName), "AddTinkCloudInitCnfAction")
	if err != nil {
		fmt.Println("fail 7")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteNetPlanActionIdx, ImageFieldName), "WriteNetPlanAction")
	if err != nil {
		fmt.Println("fail 8")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, DisableTinkCloudInitNetworkCapActionIdx, ImageFieldName), "DisableTinkCloudInitNetworkCapAction")
	if err != nil {
		fmt.Println("fail 9")
		return
	}

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println(string(b))
}

func generateDefaultBottleRocketTmpl() {
	tmpJson, err := getTmplJsonFromTmplStr(DefaultBottleRocketTmpl)
	if err != nil {
		return
	}

	actions := gjson.Get(string(tmpJson), ActionsPathFromRoot)
	if !actions.Exists() || !actions.IsArray() {
		return
	}

	var StreamImageActionIdx, WriteBootCnfActionIdx, WriteNetCnfActionIdx, WriteUserDataActionIdx, RebootActionIdx int
	for i, a := range actions.Array() {
		name := a.Get(ActionFieldName)
		if !name.Exists() {
			continue
		}

		if name.String() == BottleRocketStreamImageActionName {
			StreamImageActionIdx = i
			continue
		}

		if name.String() == BottleRocketRebootImageActionName {
			RebootActionIdx = i
			continue
		}

		if name.String() == BottleRocketWriteBootCnfActionName {
			WriteBootCnfActionIdx = i
			continue
		}

		if name.String() == BottleRocketWriteNetCnfActionName {
			WriteNetCnfActionIdx = i
			continue
		}

		if name.String() == BottleRocketWriteUserDataActionName {
			WriteUserDataActionIdx = i
			continue
		}
	}

	value, err := sjson.Set(string(tmpJson), fmt.Sprintf("%s.%d.%s.%s", ActionsPathFromRoot, StreamImageActionIdx, EnvFieldName, ImageUrlFieldName), "img_url")
	if err != nil {
		fmt.Println("fail 3")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, StreamImageActionIdx, ImageFieldName), "StreamImageAction")
	if err != nil {
		fmt.Println("fail 4")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, RebootActionIdx, ImageFieldName), "RebootAction")
	if err != nil {
		fmt.Println("fail 5")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteNetCnfActionIdx, ImageFieldName), "WriteNetCnfAction")
	if err != nil {
		fmt.Println("fail 7")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteBootCnfActionIdx, ImageFieldName), "WriteBootCnfAction")
	if err != nil {
		fmt.Println("fail 8")
		return
	}

	consoles := []string{"tty0", "ttyS0,115200n8"}
	consoles = func(cs []string) []string {
		modifiedConsoles := make([]string, 0)
		for _, c := range consoles {
			modifiedConsoles = append(modifiedConsoles, fmt.Sprintf("\"%s\"", c))
		}
		return modifiedConsoles
	}(consoles)
	allConsoles := strings.Join(consoles, ",")

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s.%s", ActionsPathFromRoot, WriteBootCnfActionIdx, EnvFieldName, "BOOTCONFIG_CONTENTS"), fmt.Sprintf("kernel { console = %s }", allConsoles))
	if err != nil {
		fmt.Println("fail kernel")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteUserDataActionIdx, ImageFieldName), "WriteUserDataAction")
	if err != nil {
		fmt.Println("fail 9")
		return
	}

	value = strings.Replace(value, "NIC_NAME", "eno1000", -1)

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println(string(b))
}

func generateEquinixBottleRocketTmpl() {
	tmpJson, err := getTmplJsonFromTmplStr(EquinixBottleRocketTmpl)
	if err != nil {
		return
	}

	actions := gjson.Get(string(tmpJson), ActionsPathFromRoot)
	if !actions.Exists() || !actions.IsArray() {
		return
	}

	var StreamImageActionIdx, WriteBootCnfActionIdx, WriteNetPlanActionIdx, WriteUserDataActionIdx, RebootActionIdx int
	for i, a := range actions.Array() {
		name := a.Get(ActionFieldName)
		if !name.Exists() {
			continue
		}

		if name.String() == EquinixBottleRocketStreamImageActionName {
			StreamImageActionIdx = i
			continue
		}

		if name.String() == EquinixBottleRocketRebootImageActionName {
			RebootActionIdx = i
			continue
		}

		if name.String() == EquinixBottleRocketWriteBootCnfActionName {
			WriteBootCnfActionIdx = i
			continue
		}

		if name.String() == EquinixBottleRocketWriteNetPlanActionName {
			WriteNetPlanActionIdx = i
			continue
		}

		if name.String() == EquinixBottleRocketWriteUserDataActionName {
			WriteUserDataActionIdx = i
			continue
		}
	}

	value, err := sjson.Set(string(tmpJson), fmt.Sprintf("%s.%d.%s.%s", ActionsPathFromRoot, StreamImageActionIdx, EnvFieldName, ImageUrlFieldName), "img_url")
	if err != nil {
		fmt.Println("fail 3")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, StreamImageActionIdx, ImageFieldName), "StreamImageAction")
	if err != nil {
		fmt.Println("fail 4")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, RebootActionIdx, ImageFieldName), "RebootAction")
	if err != nil {
		fmt.Println("fail 5")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteNetPlanActionIdx, ImageFieldName), "WriteNetPlanAction")
	if err != nil {
		fmt.Println("fail 7")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteBootCnfActionIdx, ImageFieldName), "WriteBootCnfAction")
	if err != nil {
		fmt.Println("fail 8")
		return
	}

	consoles := []string{"tty0", "ttyS0,115200n8"}
	consoles = func(cs []string) []string {
		modifiedConsoles := make([]string, 0)
		for _, c := range consoles {
			modifiedConsoles = append(modifiedConsoles, fmt.Sprintf("\"%s\"", c))
		}
		return modifiedConsoles
	}(consoles)
	allConsoles := strings.Join(consoles, ",")

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s.%s", ActionsPathFromRoot, WriteBootCnfActionIdx, EnvFieldName, "BOOTCONFIG_CONTENTS"), fmt.Sprintf("kernel { console = %s }", allConsoles))
	if err != nil {
		fmt.Println("fail kernel")
		return
	}

	value, err = sjson.Set(value, fmt.Sprintf("%s.%d.%s", ActionsPathFromRoot, WriteUserDataActionIdx, ImageFieldName), "WriteUserDataAction")
	if err != nil {
		fmt.Println("fail 9")
		return
	}

	value = strings.Replace(value, "NIC_NAME", "eno1000", -1)

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println(string(b))
}

func updateTemplateName(tmplStr string) {
	tmpJson, err := getTmplJsonFromTmplStr(tmplStr)
	if err != nil {
		return
	}

	tmplMdNamePath := "metadata.name"
	tmplNamePath := "spec.template.name"
	tmplTaskName := "spec.template.tasks.0.name"

	value, err := sjson.Set(string(tmpJson), tmplMdNamePath, "apna-cluster-name")
	if err != nil {
		return
	}

	value, err = sjson.Set(value, tmplNamePath, "apna-cluster-name")
	if err != nil {
		return
	}

	value, err = sjson.Set(value, tmplTaskName, "apna-cluster-name")
	if err != nil {
		return
	}

	b, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return
	}
	println(string(b))
}

func test() {
	tmplStr := "template:\n  global_timeout: 6000\n  id: \"\"\n  name: cp-hiteshp-eksa-cluster-m3-small-x86\n  tasks:\n  - actions:\n    - environment:\n        COMPRESSED: \"true\"\n        DEST_DISK: /dev/sda\n        IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/17/artifacts/raw/1-23/bottlerocket-v1.23.9-eks-d-1-23-5-eks-a-17-amd64.img.gz\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n      name: stream-image\n      timeout: 600\n    - environment:\n        CONTENTS: |\n          # Version is required, it will change as we support\n          # additional settings\n          version = 1\n\n          # \"eno1\" is the interface name\n          # Users may turn on dhcp4 and dhcp6 via boolean\n          [enp1s0f0np0]\n          dhcp4 = true\n          dhcp6 = false\n          # Define this interface as the \"primary\" interface\n          # for the system.  This IP is what kubelet will use\n          # as the node IP.  If none of the interfaces has\n          # \"primary\" set, we choose the first interface in\n          # the file\n          primary = true\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /net.toml\n        DIRMODE: \"0755\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        MODE: \"0644\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n      name: write-netplan\n      pid: host\n      timeout: 90\n    - environment:\n        BOOTCONFIG_CONTENTS: kernel { console = \"ttyS1,115200n8\" }\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /bootconfig.data\n        DIRMODE: \"0700\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        MODE: \"0644\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n      name: write-bootconfig\n      pid: host\n      timeout: 90\n    - environment:\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /user-data.toml\n        DIRMODE: \"0700\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        HEGEL_URLS: http://147.75.88.178:50061,http://147.75.88.189:50061\n        MODE: \"0644\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n      name: write-user-data\n      pid: host\n      timeout: 90\n    - image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\n      name: reboot-image\n      pid: host\n      timeout: 90\n      volumes:\n      - /worker:/worker\n    name: cp-hiteshp-eksa-cluster-m3-small-x86\n    volumes:\n    - /dev:/dev\n    - /dev/console:/dev/console\n    - /lib/firmware:/lib/firmware:ro\n    worker: '{{.device_1}}'\n  version: \"0.1\"\n"
	fmt.Println(tmplStr)
}

func test2() {
	js := "{\"apiVersion\":\"anywhere.eks.amazonaws.com/v1alpha1\",\"kind\":\"TinkerbellTemplateConfig\",\"metadata\":{\"name\":\"cp-hiteshp-eksa-cluster-m3-small-x86\"},\"spec\":{\"template\":{\"global_timeout\":6000,\"id\":\"\",\"name\":\"cp-hiteshp-eksa-cluster-m3-small-x86\",\"tasks\":[{\"actions\":[{\"environment\":{\"COMPRESSED\":\"true\",\"DEST_DISK\":\"/dev/sda\",\"IMG_URL\":\"https://anywhere-assets.eks.amazonaws.com/releases/bundles/17/artifacts/raw/1-23/bottlerocket-v1.23.9-eks-d-1-23-5-eks-a-17-amd64.img.gz\"},\"image\":\"public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\",\"name\":\"stream-image\",\"timeout\":600},{\"environment\":{\"CONTENTS\":\"# Version is required, it will change as we support\\n# additional settings\\nversion = 1\\n\\n# \\\"eno1\\\" is the interface name\\n# Users may turn on dhcp4 and dhcp6 via boolean\\n[enp1s0f0np0]\\ndhcp4 = true\\ndhcp6 = false\\n# Define this interface as the \\\"primary\\\" interface\\n# for the system.  This IP is what kubelet will use\\n# as the node IP.  If none of the interfaces has\\n# \\\"primary\\\" set, we choose the first interface in\\n# the file\\nprimary = true\\n\",\"DEST_DISK\":\"/dev/sda12\",\"DEST_PATH\":\"/net.toml\",\"DIRMODE\":\"0755\",\"FS_TYPE\":\"ext4\",\"GID\":\"0\",\"MODE\":\"0644\",\"UID\":\"0\"},\"image\":\"public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\",\"name\":\"write-netplan\",\"pid\":\"host\",\"timeout\":90},{\"environment\":{\"BOOTCONFIG_CONTENTS\":\"kernel { console = \\\"ttyS1,115200n8\\\" }\",\"DEST_DISK\":\"/dev/sda12\",\"DEST_PATH\":\"/bootconfig.data\",\"DIRMODE\":\"0700\",\"FS_TYPE\":\"ext4\",\"GID\":\"0\",\"MODE\":\"0644\",\"UID\":\"0\"},\"image\":\"public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\",\"name\":\"write-bootconfig\",\"pid\":\"host\",\"timeout\":90},{\"environment\":{\"DEST_DISK\":\"/dev/sda12\",\"DEST_PATH\":\"/user-data.toml\",\"DIRMODE\":\"0700\",\"FS_TYPE\":\"ext4\",\"GID\":\"0\",\"HEGEL_URLS\":\"http://147.75.88.178:50061,http://147.75.88.189:50061\",\"MODE\":\"0644\",\"UID\":\"0\"},\"image\":\"public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\",\"name\":\"write-user-data\",\"pid\":\"host\",\"timeout\":90},{\"image\":\"public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-17\",\"name\":\"reboot-image\",\"pid\":\"host\",\"timeout\":90,\"volumes\":[\"/worker:/worker\"]}],\"name\":\"cp-hiteshp-eksa-cluster-m3-small-x86\",\"volumes\":[\"/dev:/dev\",\"/dev/console:/dev/console\",\"/lib/firmware:/lib/firmware:ro\"],\"worker\":\"{{.device_1}}\"}],\"version\":\"0.1\"}}}"

	m := make(map[string]interface{})

	_ = json.Unmarshal([]byte(js), &m)

	for k, _ := range m {
		fmt.Println(k)
	}
}

func validateSubdomain(subdomain string) bool {
	// Regular expression to match lowercase alphanumeric characters, '-' or '.'
	r, _ := regexp.Compile("^[a-z0-9]([a-z0-9-.]*[a-z0-9])?$")
	return r.MatchString(subdomain)
}

func isColumnNonMandatory(field string) bool {
	for _, f := range NonMandatoryFields {
		if f == field {
			return true
		}
	}
	return false
}

func encodeCSV(columns []string, rows []map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(columns); err != nil {
		return nil, err
	}
	r := make([]string, len(columns))
	for _, row := range rows {
		for i, column := range columns {
			r[i] = row[column]
		}
		if err := w.Write(r); err != nil {
			return nil, err
		}
	}

	w.Flush()

	return buf.Bytes(), nil
}

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

const GenericTmpl = `
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: TinkerbellTemplateConfig
metadata:
  name: my-cluster-name
spec:
  template:
    global_timeout: 6000
    id: ""
    name: my-cluster-name
    tasks:
    - actions:
      - environment:
          COMPRESSED: "true"
          DEST_DISK: /dev/sda
          IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/11/artifacts/raw/1-22/bottlerocket-v1.22.10-eks-d-1-22-8-eks-a-11-amd64.img.gz
        image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
        name: stream-image
        timeout: 360
      name: my-cluster-name
      volumes:
      - /dev:/dev
      - /dev/console:/dev/console
      - /lib/firmware:/lib/firmware:ro
      worker: '{{.device_1}}'
    version: "0.1"
`

const GenericWriteBootConfigActionTmpl = `
environment:
  BOOTCONFIG_CONTENTS: |
    kernel {}
  DEST_DISK: /dev/sda12
  DEST_PATH: /bootconfig.data
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0644"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: write-bootconfig
timeout: 90
`

const GenericWriteNetConfigActionTmpl = `
environment:
  CONTENTS: |
    # Version is required, it will change as we support
    # additional settings
    version = 1
    # "eno1" is the interface name
    # Users may turn on dhcp4 and dhcp6 via boolean
    [NIC_NAME]
    dhcp4 = true
    # Define this interface as the "primary" interface
    # for the system.  This IP is what kubelet will use
    # as the node IP.  If none of the interfaces has
    # "primary" set, we choose the first interface in
    # the file
    primary = true            
  DEST_DISK: /dev/sda12
  DEST_PATH: /net.toml
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0644"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: write-netconfig
timeout: 90
`

const GenericWriteUserDataActionTmpl = `
environment:
  HEGEL_URL: http://<hegel-ip>:50061
  DEST_DISK: /dev/sda12
  DEST_PATH: /user-data.toml
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0644"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: write-user-data
timeout: 90
`

const GenericRebootActionTmpl = `
name: "reboot"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
timeout: 90
volumes:
  - /worker:/worker
`

const GenericWriteNetPlanTmpl = `
environment:
  DEST_DISK: /dev/sda2
  DEST_PATH: /etc/netplan/config.yaml
  STATIC_NETPLAN: true
  DIRMODE: "0755"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0644"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: write-netplan
timeout: 90
`

const GenericAddTinkCloudInitConfigTmpl = `
environment:
  CONTENTS: |
    datasource:
      Ec2:
        metadata_urls: [<admin-machine-ip>, <tinkerbell-ip-from-cluster-config>]
        strict_id: false
        manage_etc_hosts: localhost
        warnings:
        dsid_missing_source: off            
  DEST_DISK: /dev/sda2
  DEST_PATH: /etc/cloud/cloud.cfg.d/10_tinkerbell.cfg
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0600"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: add-tink-cloud-init-config
timeout: 90
`

const GenericDisableCloudInitNetworkCapabilitiesTmpl = `
environment:
  CONTENTS: |
    network:
      config: disabled            
  DEST_DISK: /dev/sda2
  DEST_PATH: /etc/cloud/cloud.cfg.d/99-disable-network-config.cfg
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0600"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: disable-cloud-init-network-capabilities
`

const GenericAddTinkCloudDsConfigTmpl = `
environment:
  CONTENTS: |
    datasource: Ec2
  DEST_DISK: /dev/sda2
  DEST_PATH: /etc/cloud/ds-identify.cfg
  DIRMODE: "0700"
  FS_TYPE: ext4
  GID: "0"
  MODE: "0600"
  UID: "0"
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: add-tink-cloud-init-ds-config
timeout: 90
`

const GenericKexecTmpl = `
environment:
  BLOCK_DEVICE: /dev/sda2
  FS_TYPE: ext4
image: public.ecr.aws/eks-anywhere/tinkerbell/hub/kexec:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-15
name: kexec-image
pid: host
timeout: 90
`

const tmpl = "template:\n  global_timeout: 6000\n  id: \"\"\n  name: tmpl-1\n  tasks:\n  - actions:\n    - environment:\n        COMPRESSED: \"true\"\n        DEST_DISK: /dev/sda\n        IMG_URL: https://anywhere-assets.eks.amazonaws.com/releases/bundles/27/artifacts/raw/1-24/bottlerocket-v1.24.9-eks-d-1-24-7-eks-a-27-amd64.img.gz\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/image2disk:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-27\n      name: stream-image\n      timeout: 600\n    - environment:\n        BOOTCONFIG_CONTENTS: |\n          kernel {}\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /bootconfig.data\n        DIRMODE: \"0700\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        MODE: \"0644\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-27\n      name: write-bootconfig\n      pid: host\n      timeout: 90\n    - environment:\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /user-data.toml\n        DIRMODE: \"0700\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        HEGEL_URLS: http://147.75.88.61:50061,http://2.3.4.5:50061\n        MODE: \"0644\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-27\n      name: write-user-data\n      pid: host\n      timeout: 90\n    - environment:\n        DEST_DISK: /dev/sda12\n        DEST_PATH: /net.toml\n        DIRMODE: \"0755\"\n        FS_TYPE: ext4\n        GID: \"0\"\n        IFNAME: enp0s3\n        MODE: \"0644\"\n        STATIC_BOTTLEROCKET: \"true\"\n        UID: \"0\"\n      image: public.ecr.aws/eks-anywhere/tinkerbell/hub/writefile:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-27\n      name: write-netplan\n      pid: host\n      timeout: 90\n    - image: public.ecr.aws/eks-anywhere/tinkerbell/hub/reboot:6c0f0d437bde2c836d90b000312c8b25fa1b65e1-eks-a-27\n      name: reboot-image\n      pid: host\n      timeout: 90\n      volumes:\n      - /worker:/worker\n    name: tmpl-1\n    volumes:\n    - /dev:/dev\n    - /dev/console:/dev/console\n    - /lib/firmware:/lib/firmware:ro\n    worker: '{{.device_1}}'\n  version: \"0.1\"\n"

const LowerCaseRFC1123SubDomainRegex = "^[a-z0-9]([-\\.]?[a-z0-9])+[a-z0-9]$"

const vg = "Vagrant.configure(2) do |config|\n  config.vm.box = 'bento/ubuntu-20.04'\n  config.vm.network :forwarded_port, guest: 22, host: 2322, id: \"ssh\"\n  config.vm.hostname = '${EKSA_ADMIN_VM_NAME}'\n  config.vm.box_check_update = false\n  config.disksize.size = '100GB'\n  config.vm.boot_timeout = 300\n  config.persistent_storage.enabled = true\n  config.persistent_storage.location = \"virtualdrive.vdi\"\n  config.persistent_storage.size = 102400\n  config.persistent_storage.diskdevice = '/dev/sdc'\n  config.persistent_storage.partition = false\n  config.persistent_storage.use_lvm = false\n  config.vm.provider 'virtualbox' do |vb|\n    vb.cpus = ${EKSA_ADMIN_VM_CPUS}\n    vb.memory = ${EKSA_ADMIN_VM_MEM}\n    vb.name = '${EKSA_ADMIN_VM_NAME}'\n    vb.customize ['modifyvm', :id, '--macaddress1', '${EKSA_ADMIN_MAC_CONCISE}']\n  end\nend\n\n"

var mapDemo map[string]interface{} = map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{"creationTimestamp": "2023-04-13T11:26:50Z", "labels": map[string]interface{}{"run": "nginx-101"}, "managedFields": []interface{}{map[string]interface{}{"apiVersion": "v1", "fieldsType": "FieldsV1", "fieldsV1": map[string]interface{}{"f:metadata": map[string]interface{}{"f:labels": map[string]interface{}{".": map[string]interface{}{}, "f:run": map[string]interface{}{}}}, "f:spec": map[string]interface{}{"f:containers": map[string]interface{}{"k:{\"name\":\"nginx-101\"}": map[string]interface{}{".": map[string]interface{}{}, "f:image": map[string]interface{}{}, "f:imagePullPolicy": map[string]interface{}{}, "f:name": map[string]interface{}{}, "f:resources": map[string]interface{}{}, "f:terminationMessagePath": map[string]interface{}{}, "f:terminationMessagePolicy": map[string]interface{}{}}}, "f:dnsPolicy": map[string]interface{}{}, "f:enableServiceLinks": map[string]interface{}{}, "f:restartPolicy": map[string]interface{}{}, "f:schedulerName": map[string]interface{}{}, "f:securityContext": map[string]interface{}{}, "f:terminationGracePeriodSeconds": map[string]interface{}{}}}, "manager": "kubectl-run", "operation": "Update", "time": "2023-04-13T11:26:50Z"}, map[string]interface{}{"apiVersion": "v1", "fieldsType": "FieldsV1", "fieldsV1": map[string]interface{}{"f:status": map[string]interface{}{"f:conditions": map[string]interface{}{"k:{\"type\":\"ContainersReady\"}": map[string]interface{}{".": map[string]interface{}{}, "f:lastProbeTime": map[string]interface{}{}, "f:lastTransitionTime": map[string]interface{}{}, "f:status": map[string]interface{}{}, "f:type": map[string]interface{}{}}, "k:{\"type\":\"Initialized\"}": map[string]interface{}{".": map[string]interface{}{}, "f:lastProbeTime": map[string]interface{}{}, "f:lastTransitionTime": map[string]interface{}{}, "f:status": map[string]interface{}{}, "f:type": map[string]interface{}{}}, "k:{\"type\":\"Ready\"}": map[string]interface{}{".": map[string]interface{}{}, "f:lastProbeTime": map[string]interface{}{}, "f:lastTransitionTime": map[string]interface{}{}, "f:status": map[string]interface{}{}, "f:type": map[string]interface{}{}}}, "f:containerStatuses": map[string]interface{}{}, "f:hostIP": map[string]interface{}{}, "f:phase": map[string]interface{}{}, "f:podIP": map[string]interface{}{}, "f:podIPs": map[string]interface{}{".": map[string]interface{}{}, "k:{\"ip\":\"172.17.0.3\"}": map[string]interface{}{".": map[string]interface{}{}, "f:ip": map[string]interface{}{}}}, "f:startTime": map[string]interface{}{}}}, "manager": "kubelet", "operation": "Update", "subresource": "status", "time": "2023-04-13T11:26:55Z"}}, "name": "nginx-101", "namespace": "default", "resourceVersion": "95513", "uid": "88f35f48-d87b-4736-a61a-6404e17e30e9"}, "spec": map[string]interface{}{"containers": []interface{}{map[string]interface{}{"image": "nginx", "imagePullPolicy": "Always", "name": "nginx-101", "resources": map[string]interface{}{}, "terminationMessagePath": "/dev/termination-log", "terminationMessagePolicy": "File", "volumeMounts": []interface{}{map[string]interface{}{"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount", "name": "kube-api-access-49xtd", "readOnly": true}}}}, "dnsPolicy": "ClusterFirst", "enableServiceLinks": true, "nodeName": "minikube", "preemptionPolicy": "PreemptLowerPriority", "priority": 0, "restartPolicy": "Always", "schedulerName": "default-scheduler", "securityContext": map[string]interface{}{}, "serviceAccount": "default", "serviceAccountName": "default", "terminationGracePeriodSeconds": 30, "tolerations": []interface{}{map[string]interface{}{"effect": "NoExecute", "key": "node.kubernetes.io/not-ready", "operator": "Exists", "tolerationSeconds": 300}, map[string]interface{}{"effect": "NoExecute", "key": "node.kubernetes.io/unreachable", "operator": "Exists", "tolerationSeconds": 300}}, "volumes": []interface{}{map[string]interface{}{"name": "kube-api-access-49xtd", "projected": map[string]interface{}{"defaultMode": 420, "sources": []interface{}{map[string]interface{}{"serviceAccountToken": map[string]interface{}{"expirationSeconds": 3607, "path": "token"}}, map[string]interface{}{"configMap": map[string]interface{}{"items": []interface{}{map[string]interface{}{"key": "ca.crt", "path": "ca.crt"}}, "name": "kube-root-ca.crt"}}, map[string]interface{}{"downwardAPI": map[string]interface{}{"items": []interface{}{map[string]interface{}{"fieldRef": map[string]interface{}{"apiVersion": "v1", "fieldPath": "metadata.namespace"}, "path": "namespace"}}}}}}}}}, "status": map[string]interface{}{"conditions": []interface{}{map[string]interface{}{"lastProbeTime": interface{}(nil), "lastTransitionTime": "2023-04-13T11:26:50Z", "status": "True", "type": "Initialized"}, map[string]interface{}{"lastProbeTime": interface{}(nil), "lastTransitionTime": "2023-04-13T11:26:55Z", "status": "True", "type": "Ready"}, map[string]interface{}{"lastProbeTime": interface{}(nil), "lastTransitionTime": "2023-04-13T11:26:55Z", "status": "True", "type": "ContainersReady"}, map[string]interface{}{"lastProbeTime": interface{}(nil), "lastTransitionTime": "2023-04-13T11:26:50Z", "status": "True", "type": "PodScheduled"}}, "containerStatuses": []interface{}{map[string]interface{}{"containerID": "docker://48475e7eb971398b2e5bf83ca86b149863fe9495019854af9db11eb0eb9ee318", "image": "nginx:latest", "imageID": "docker-pullable://nginx@sha256:63b44e8ddb83d5dd8020327c1f40436e37a6fffd3ef2498a6204df23be6e7e94", "lastState": map[string]interface{}{}, "name": "nginx-101", "ready": true, "restartCount": 0, "started": true, "state": map[string]interface{}{"running": map[string]interface{}{"startedAt": "2023-04-13T11:26:54Z"}}}}, "hostIP": "192.168.49.2", "phase": "Running", "podIP": "172.17.0.3", "podIPs": []interface{}{map[string]interface{}{"ip": "172.17.0.3"}}, "qosClass": "BestEffort", "startTime": "2023-04-13T11:26:50Z"}}

type TestStruct struct {
	Type string `json:"type"`
}

func main() {

	// jsonStr, err := json.Marshal(mapDemo)
	// if err != nil {
	// 	fmt.Printf("Error: %s", err.Error())
	// } else {
	// 	fmt.Println(string(jsonStr))
	// }

	// GetAllServices()

	// Run()

	// data := []byte("")
	// data2 := []byte(`{"type": "string"}`)

	// ans := TestStruct{}
	// json.Unmarshal(data2, &ans)

	// fmt.Println(ans)

	// ans2 := TestStruct{}
	// json.Unmarshal(data, &ans2)

	// fmt.Println(ans2)
	// action := &Action{}
	// err := yaml.Unmarshal([]byte(GenericWriteNetConfigActionTmpl), action)
	// if err != nil {
	// 	fmt.Println("error while unmarshalling GenericWriteNetConfigActionTmpl. Error:", err)
	// 	return
	// }

	// if action == nil {
	// 	fmt.Println("actioncould not be unmarshalled")
	// 	return
	// }

	// encodedAction, err := encodeToBase64(action)
	// if err != nil {
	// 	fmt.Println("failed to encode action. Error: ", err)
	// 	return
	// }

	// fmt.Println(encodedAction)

	// var data bytes.Buffer
	// wr := csv.NewWriter(&data)
	// wr.Write([]string{"test1", "test2", "test3"})
	// wr.Flush()
	// fmt.Println(string(data.Bytes()))

	// fmt.Println(fmt.Sprintf("mv %s/eksa-scale-cluster.log %s/eksa-scale-cluster-%s.log", "abc", "abc", time.Now().Format("2006-01-02-15-04-05")))

	// cnf := config{}

	// err := json.Unmarshal([]byte(jsonClusterConfig), &cnf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// err = testCert(cnf.CertBundles[0].Data)
	// fmt.Println("err", err)

	// var clusterYaml interface{}

	// err := yaml.Unmarshal([]byte(yamlData), clusterYaml)
	// if err != nil {
	// 	fmt.Println("error while unmarshalling  spec", err)

	// }

	// fmt.Println(string(yamlData))

	// rootCmd.AddCommand(randomCmd)

	// randomCmd.PersistentFlags().String("term", "", "A search term for a dad joke.")

	// fileContent, err := ioutil.ReadFile("dist/file.txt")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	// fmt.Println(encodedContent)

	// sr := ShortenRequest{}

	// fmt.Println(sr.URL)

	// legalURLChars := `^(http://|https://)[A-Za-z0-9\-\.]+\.[A-Za-z]{2,4}(:[0-9]+)?(/.*)?$`

	// matched, _ := regexp.Match(legalURLChars, []byte(sr.URL))

	// fmt.Println(matched)

	// url1 := "https://github.com/paralus/paralus/issues/234"
	// url2 := "https://github.com/paralus/paralus/issues/234"
	// url3 := "https://anywhere.eks.amazonaws.com/docs/getting-started/docker/"

	// shortUrl1, err := shorten(url1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("shorturl1: ", shortUrl1)
	// shortUrl2, err := shorten(url2)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("shorturl2: ", shortUrl2)
	// shortUrl3, err := shorten(url3)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("shorturl3: ", shortUrl3)

	// testCollisions_ankit_jain()

	testCollisions_subham_sahare()

}

func testCollisions_subham_sahare() {
	baseURL := "https://example.com"
	count := 5000

	similarURLs := generateSimilarURLs(baseURL, count)

	shortURLs := make(map[string]string)

	// Print the generated similar URLs
	for i, url := range similarURLs {
		fmt.Println(i)
		shortURL := shorten_subham_sahare(url)
		if longUrl, exists := shortURLs[shortURL]; exists && (url != longUrl) {
			// Collision detected
			panic(fmt.Sprintf("Collision detected for URL: %s", url))
		}
		shortURLs[shortURL] = url
	}
	fmt.Println(len(shortURLs))
}

func generateCRC32Encoding(input string) uint32 {
	crc32Hash := crc32.NewIEEE()
	crc32Hash.Write([]byte(input))
	checksum := crc32Hash.Sum32()
	return checksum
}

func shorten_subham_sahare(sourceUrl string) string {
	crc32Encoding := generateCRC32Encoding(sourceUrl)
	shortUrl := base62.Encode(int(crc32Encoding))
	return shortUrl
}

func testCollisions_ankit_jain() {
	// Example usage:
	baseURL := "https://example.com"
	count := 5000

	similarURLs := generateSimilarURLs(baseURL, count)
	fmt.Println(len(similarURLs))

	shortURLs := make(map[string]string)

	// Print the generated similar URLs
	for _, url := range similarURLs {
		// fmt.Println(url)
		shortURL, _ := shorten_ankit_jain(url)
		if longUrl, exists := shortURLs[shortURL]; exists && (url != longUrl) {
			// Collision detected
			panic(fmt.Sprintf("Collision detected for URL: %s", url))
		}
		shortURLs[shortURL] = url
	}
}

func generateSimilarURLs(baseURL string, count int) []string {
	// Initialize the random number generator with a seed
	rand.Seed(time.Now().UnixNano())

	// Create a slice to store the generated URLs
	urls := make([]string, 0, count)

	for i := 0; i < count; i++ {
		// Generate a random variation of the baseURL
		variation := fmt.Sprintf("-%d", rand.Intn(count))

		// Concatenate the variation with the baseURL to create a similar URL
		similarURL := baseURL + variation

		// Append the similar URL to the slice
		urls = append(urls, similarURL)
	}

	return urls
}

type hashAlgo struct {
	hash hash.Hash64
}

func shorten_ankit_jain(url string) (string, error) {
	hashAlgo := fnv.New64a()
	hashAlgo.Write([]byte(url))
	hashed := fmt.Sprintf("%x", hashAlgo.Sum(nil))
	hashed = strings.ReplaceAll(hashed, "/", "")
	hashed = strings.ReplaceAll(hashed, "+", "")

	shortURL := hashed[:6]
	return shortURL, nil
}

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random dad joke",
	Long:  `This command fetches a random dadjoke from the icanhazdadjoke api`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

func getRandomJokeWithTerm(jokeTerm string) {
	fmt.Printf("You searched for a joke with the term: %v", jokeTerm)
}

func getRandomJoke() {
	fmt.Println("You searched for a joke")
}

type config struct {
	CertBundles []CertBundle `json:"certBundles"`
}

type CertBundle struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var jsonClusterConfig = `
{
	"certBundles": [
		{
			"name": "bundle_1",
			"data": "-----BEGIN CERTIFICATE-----\nMIICDDCCAXUCFFmejalQeOFMOCZ9fo/jRoIxuGSHMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl\ncm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMjMwNjI3MTEwMDIwWhcNMjQwNjI2MTEw\nMDIwWjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UE\nCgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GN\nADCBiQKBgQDQFJJAlnqahO6Rmx/pfJRF2mKV4EHExGNWWeiPD5/sx0Yu6GqNNyJb\n5DkTlP0wpCPfvkewjxZcrAondz64nLrHbTCiouzOSllnJK4GIqjwp17NoqGX5i8K\nK7r6U8pYMznME/GLCGAsDKz2yzgoCZH6mpK+iACE/gS2F/g17q44rwIDAQABMA0G\nCSqGSIb3DQEBCwUAA4GBADEeDEEp/wJXuKVGCdAsqgDxoMhIwD6TyGC91CV2AsmH\ne4fuG0Ld47mfaFamjB6K+bkeJlQyZmSX2dI3t9DnqgLy6PNHplibJkN56dQ4rKy8\nqyGU1QxiBK70lni4QuoeNvhFmlzp5mI654tmXSnoCgriUXtil3uQPmihs+KTtAdT\n-----END CERTIFICATE-----"
		}
	]
}
`

var yamlData = `\x61706956657273696f6e3a20696e6672612e6b38736d676d742e696f2f76330a6b696e643a20436c75737465720a6d657461646174613a0a20206c6162656c733a0a2020202072616661792e6465762f636c75737465724e616d653a206a756c3132636c
75737465720a2020202072616661792e6465762f636c7573746572547970653a20656b73615f626d0a20206e616d653a206a756c3132636c75737465720a202070726f6a6563743a2074657374696e670a737065633a0a2020626c75657072696e743a0a202
020206e616d653a206d696e696d616c0a2020202076657273696f6e3a206c61746573740a2020636f6e6669673a0a20202020656b7361436c7573746572436f6e6669673a0a20202020202061706956657273696f6e3a20616e7977686572652e656b732e61
6d617a6f6e6177732e636f6d2f7631616c706861310a2020202020206b696e643a20436c75737465720a2020202020206d657461646174613a0a20202020202020206e616d653a206a756c3132636c75737465720a202020202020737065633a0a202020202
0202020636c75737465724e6574776f726b3a0a20202020202020202020636e69436f6e6669673a0a20202020202020202020202063696c69756d3a0a2020202020202020202020202020706f6c696379456e666f7263656d656e744d6f64653a2064656661
756c740a20202020202020202020706f64733a0a20202020202020202020202063696472426c6f636b733a0a2020202020202020202020202d203139322e3136382e302e302f31360a2020202020202020202073657276696365733a0a20202020202020202
020202063696472426c6f636b733a0a2020202020202020202020202d2031302e39362e302e302f31320a2020202020202020636f6e74726f6c506c616e65436f6e66696775726174696f6e3a0a20202020202020202020636f756e743a20310a2020202020
2020202020656e64706f696e743a0a202020202020202020202020686f73743a203139322e3136382e31302e3235320a202020202020202020206d616368696e6547726f75705265663a0a2020202020202020202020206b696e643a2054696e6b657262656
c6c4d616368696e65436f6e6669670a2020202020202020202020206e616d653a206d63310a20202020202020206461746163656e7465725265663a0a202020202020202020206b696e643a2054696e6b657262656c6c4461746163656e746572436f6e6669
670a202020202020202020206e616d653a206a756c3132636c75737465720a20202020202020206b756265726e6574657356657273696f6e3a2022312e3236220a20202020202020206d616e6167656d656e74436c75737465723a0a2020202020202020202
06e616d653a206a756c3132636c75737465720a2020202020202020776f726b65724e6f646547726f7570436f6e66696775726174696f6e733a0a20202020202020202d20636f756e743a20310a202020202020202020206d616368696e6547726f75705265
663a0a2020202020202020202020206b696e643a2054696e6b657262656c6c4d616368696e65436f6e6669670a2020202020202020202020206e616d653a206d63320a202020202020202020206e616d653a206d642d300a2020202074696e6b657262656c6
c4461746163656e746572436f6e6669673a0a20202020202061706956657273696f6e3a20616e7977686572652e656b732e616d617a6f6e6177732e636f6d2f7631616c706861310a2020202020206b696e643a2054696e6b657262656c6c4461746163656e
746572436f6e6669670a2020202020206d657461646174613a0a20202020202020206e616d653a206a756c3132636c75737465720a202020202020737065633a0a202020202020202074696e6b657262656c6c49503a203139322e3136382e31302e3231360
a2020202074696e6b657262656c6c4861726477617265436f6e6669673a0a202020202d206469736b3a202f6465762f7364610a202020202020676174657761793a203139322e3136382e31302e310a202020202020686f73746e616d653a20656b7361626d
6870312d63702d6e2d310a20202020202069705f616464726573733a203139322e3136382e31302e3131300a2020202020206c6162656c733a20747970653d63700a2020202020206d61633a2030383a30303a32373a39393a38423a36350a2020202020206
e616d65736572766572733a20382e382e382e380a2020202020206e65746d61736b3a203235352e3235352e3235352e300a202020202d206469736b3a202f6465762f7364610a202020202020676174657761793a203139322e3136382e31302e310a202020
202020686f73746e616d653a20656b7361626d6870312d64702d6e2d310a20202020202069705f616464726573733a203139322e3136382e31302e34360a2020202020206c6162656c733a20747970653d64700a2020202020206d61633a2030383a30303a3
2373a44303a32303a38300a2020202020206e616d65736572766572733a20382e382e382e380a2020202020206e65746d61736b3a203235352e3235352e3235352e300a2020202074696e6b657262656c6c4d616368696e65436f6e6669673a0a202020202d
2061706956657273696f6e3a20616e7977686572652e656b732e616d617a6f6e6177732e636f6d2f7631616c706861310a2020202020206b696e643a2054696e6b657262656c6c4d616368696e65436f6e6669670a2020202020206d657461646174613a0a2
0202020202020206e616d653a206d63310a202020202020737065633a0a2020202020202020686172647761726553656c6563746f723a0a20202020202020202020747970653a2063700a20202020202020206f7346616d696c793a20626f74746c65726f63
6b65740a202020202020202074656d706c6174655265663a0a202020202020202020206b696e643a2054696e6b657262656c6c54656d706c617465436f6e6669670a202020202020202075736572733a0a20202020202020202d206e616d653a206563322d7
57365720a20202020202020202020737368417574686f72697a65644b6579733a0a202020202020202020202d207373682d727361204141414142334e7a614331796332454141414144415141424141414267514442743930344e4c7436674a697130634a45
4e72433077387130656a515047667279342f44422b676c456c6172736670366f78474a4e6f654f543353682f2f546271592f6a534c33714f63724e515459353342306d5675726253304e33536b61764278425a306f3850724d7079483773472b384b65586c4
f34306634632f774638355637435a5638326a613179796330596b555a3871386c58416133507639534b6f51333447512f79533467384d6e36514a317a586a754439683865496c4d74614847542f325a6b584d48774135332f436d48525959594f5251573465
7a4f45394f39716a474a5158567662386b55676f4a55743774744d303739684248385a394c61424979372f54686264616b78774f556336424a68453966685945315a6f43306848565376494445336f72732b6d6b6955577655545a425535687038494342537
235623668726f45445a644d6b4d5963394e62414237462f42684b42685a75323050336e776b773775506770613739334a5a4637597558685038544e342f4a726175716a2b777769394665373674673549777a70486a43647773556e5434654433484b747155
46714b53617749372b7769586739796b794567724938556a7a5946486c55465533344e3361612b414f6b32307258634d7459412f357454486b4c477a745276306a385141434f7649666b344f7273337136713679786b48634d3d0a202020202020202020202
020726f6f7440726f626269652d656b73612d696e73740a202020202d2061706956657273696f6e3a20616e7977686572652e656b732e616d617a6f6e6177732e636f6d2f7631616c706861310a2020202020206b696e643a2054696e6b657262656c6c4d61
6368696e65436f6e6669670a2020202020206d657461646174613a0a20202020202020206e616d653a206d63320a202020202020737065633a0a2020202020202020686172647761726553656c6563746f723a0a20202020202020202020747970653a20647
00a20202020202020206f7346616d696c793a20626f74746c65726f636b65740a202020202020202074656d706c6174655265663a0a202020202020202020206b696e643a2054696e6b657262656c6c54656d706c617465436f6e6669670a20202020202020
2075736572733a0a20202020202020202d206e616d653a206563322d757365720a20202020202020202020737368417574686f72697a65644b6579733a0a202020202020202020202d207373682d727361204141414142334e7a61433179633245414141414
4415141424141414267514442743930344e4c7436674a697130634a454e72433077387130656a515047667279342f44422b676c456c6172736670366f78474a4e6f654f543353682f2f546271592f6a534c33714f63724e515459353342306d567572625330
4e33536b61764278425a306f3850724d7079483773472b384b65586c4f34306634632f774638355637435a5638326a613179796330596b555a3871386c58416133507639534b6f51333447512f79533467384d6e36514a317a586a754439683865496c4d746
14847542f325a6b584d48774135332f436d48525959594f52515734657a4f45394f39716a474a5158567662386b55676f4a55743774744d303739684248385a394c61424979372f54686264616b78774f556336424a68453966685945315a6f433068485653
76494445336f72732b6d6b6955577655545a425535687038494342537235623668726f45445a644d6b4d5963394e62414237462f42684b42685a75323050336e776b773775506770613739334a5a4637597558685038544e342f4a726175716a2b777769394
665373674673549777a70486a43647773556e5434654433484b74715546714b53617749372b7769586739796b794567724938556a7a5946486c55465533344e3361612b414f6b32307258634d7459412f357454486b4c477a745276306a385141434f764966
6b344f7273337136713679786b48634d3d0a202020202020202020202020726f6f7440726f626269652d656b73612d696e73740a2020202074696e6b657262656c6c54656d706c617465436f6e6669673a205b5d0a2020747970653a20456b73615f626d0a`

func buildGenericUbuntuRocketTmpl() (string, error) {
	updatedTmpl := ""
	tmpJson, err := getTmplJsonFromTmplStr(GenericTmpl)
	if err != nil {
		return "", err
	}

	actionPath := "spec.template.tasks.0.actions"

	currLen, exists := getCurrentLengthOfActionsInGivenTemplate(string(tmpJson), actionPath)
	if !exists {
		return "", errors.New("action path is missing")
	}

	updatedTmpl, err = addActionInTmpl(GenericWriteNetPlanTmpl, actionPath, string(tmpJson), currLen)
	if err != nil {
		return "", err
	}

	currLen, exists = getCurrentLengthOfActionsInGivenTemplate(updatedTmpl, actionPath)
	if !exists {
		return "", errors.New("action path is missing")
	}

	updatedTmpl, err = addActionInTmpl(GenericAddTinkCloudInitConfigTmpl, actionPath, updatedTmpl, currLen)
	if err != nil {
		return "", err
	}

	currLen, exists = getCurrentLengthOfActionsInGivenTemplate(updatedTmpl, actionPath)
	if !exists {
		return "", errors.New("action path is missing")
	}

	updatedTmpl, err = addActionInTmpl(GenericDisableCloudInitNetworkCapabilitiesTmpl, actionPath, updatedTmpl, currLen)
	if err != nil {
		return "", err
	}

	currLen, exists = getCurrentLengthOfActionsInGivenTemplate(updatedTmpl, actionPath)
	if !exists {
		return "", errors.New("action path is missing")
	}

	updatedTmpl, err = addActionInTmpl(GenericAddTinkCloudDsConfigTmpl, actionPath, updatedTmpl, currLen)
	if err != nil {
		return "", err
	}

	currLen, exists = getCurrentLengthOfActionsInGivenTemplate(updatedTmpl, actionPath)
	if !exists {
		return "", errors.New("action path is missing")
	}

	updatedTmpl, err = addActionInTmpl(GenericKexecTmpl, actionPath, updatedTmpl, currLen)
	if err != nil {
		return "", err
	}

	return updatedTmpl, nil
}

func addActionInTmpl(actionStr string, actionPath string, updatedTmpl string, actionIdx int) (string, error) {
	if actionIdx > 0 {
		actionJson, err := getTmplJsonFromTmplStr(actionStr)
		if err != nil {
			return "", err
		}

		updatedTmpl, err = sjson.Set(updatedTmpl, fmt.Sprintf("%s.%d", actionPath, actionIdx), actionJson)
		if err != nil {
			return "", err
		}
	}
	return updatedTmpl, nil
}

func getCurrentLengthOfActionsInGivenTemplate(tmpJson string, actionPath string) (int, bool) {
	res := gjson.Get(tmpJson, actionPath)
	if !res.Exists() {
		return 0, false
	}

	currLen := 0
	if res.IsArray() {
		currLen = len(res.Array())
	}
	return currLen, true
}

type Action struct {
	Environment map[string]string `json:"environment,omitempty"`
	Name        string            `json:"name,omitempty"`
	Image       string            `json:"image,omitempty"`
	Timeout     int               `json:"timeout,omitempty"`
	Pid         string            `json:"pid,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
}

type Task struct {
	Name    string    `json:"name"`
	Actions []*Action `json:"actions"`
	Volumes []string  `json:"volumes"`
	Worker  string    `json:"worker"`
}

type TinkerbellTemplate struct {
	GlobalTimeout int     `json:"global_timeout"`
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Tasks         []*Task `json:"tasks"`
	Version       string  `json:"version"`
}

type TinkerbellTemplateSpec struct {
	Template *TinkerbellTemplate `json:"template"`
}
