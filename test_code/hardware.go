package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Labels map[string]string
type Specs map[string]interface{}

type LabelsAndSpecs struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Specs  interface{}       `json:"specs"`
}

func GetHardwareLabelsAndSpecs() (Labels, Specs) {

	labels := make(map[string]string)
	specs := make(map[string]interface{})

	ls := []LabelsAndSpecs{}

	// output, err := ioutil.ReadFile("dist/hardware-config.json")
	// if err != nil {
	// 	log.Fatal("Error when opening file: ", err)
	// }

	err := json.Unmarshal([]byte(hardwareData), &ls)
	if err != nil {
		log.Fatal("failed to unmarshal hardware labels and specs: ", err)
		return labels, specs
	}

	for _, l := range ls {
		labels[l.Name] = labelMapToString(l.Labels)
		specs[l.Name] = l.Specs
	}

	return labels, specs
}

func labelMapToString(mp map[string]string) string {
	str := ""

	for k, v := range mp {
		str = fmt.Sprintf("%s,%s", str, fmt.Sprintf("%s=%s", k, v))
	}

	return str[1:]
}

func Run() {
	labels, specs := GetHardwareLabelsAndSpecs()

	for k, v := range labels {
		b, _ := json.Marshal(specs[k])
		// label := labels[k]

		fmt.Println("hostname", k)
		fmt.Println("label", v)
		fmt.Println("spec", string(b))
	}
}

const hardwareData = `[
    {
        "Name": "eksabmhp1-cp-n-1",
        "Labels": {
            "type": "cp",
            "v1alpha1.tinkerbell.org/ownerName": "mgmt-control-plane-template-1685978142418-nhg4r",
            "v1alpha1.tinkerbell.org/ownerNamespace": "eksa-system"
        },
        "Specs": {
            "disks": [
                {
                    "device": "/dev/sda"
                }
            ],
            "interfaces": [
                {
                    "dhcp": {
                        "arch": "x86_64",
                        "hostname": "eksabmhp1-cp-n-1",
                        "ip": {
                            "address": "192.168.10.11",
                            "family": 4,
                            "gateway": "192.168.10.1",
                            "netmask": "255.255.255.0"
                        },
                        "lease_time": 4294967294,
                        "mac": "08:00:27:63:c6:1b",
                        "name_servers": [
                            "8.8.8.8"
                        ],
                        "uefi": true
                    },
                    "netboot": {
                        "allowPXE": true,
                        "allowWorkflow": true
                    }
                }
            ],
            "metadata": {
                "facility": {
                    "facility_code": "onprem",
                    "plan_slug": "c2.medium.x86"
                },
                "instance": {
                    "allow_pxe": true,
                    "always_pxe": true,
                    "hostname": "eksabmhp1-cp-n-1",
                    "id": "08:00:27:63:c6:1b",
                    "ips": [
                        {
                            "address": "192.168.10.11",
                            "family": 4,
                            "gateway": "192.168.10.1",
                            "netmask": "255.255.255.0",
                            "public": true
                        }
                    ],
                    "operating_system": {},
                    "state": "provisioned"
                },
                "state": "in_use"
            },
            "userData": "\n[settings.host-containers.admin]\nenabled = true\nsuperpowered = true\nuser-data = \"CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FESktDaGllZEhqb2UwRDhZclAyUGFFQVR3bllYZ3JsanpnQ2gvSGIrcXNETXdYOU8vUE9IL0dod0V5WE1SSlZkdmo2eEJubXRyS1RRKzJneS9BcGxjelZvVTZkcnR0UzJ1Vk11aW1LZkNnZUFQOGthTHhwWUlNcWF6QkIvanZjK3paNUlScDJIcEMyamJ4ZndBSWdRVzNYaVdQcG9HbkFyYlpvLys2YVhMVmswRi9SclN2bUxEYlZGQTZKR0h1Qi9KY3pKaDBCZHB3RmxXR1JaWGRabEwzZUl1dmRyMXZjeWZPSStCU1RqbzRyRXc3TFEwN2drUFpjTStjVktWTE9YcEVWUThrTXduOHNHSnQxWXR3Qkw4UFgvWTU0TTRGL3JWTlVyWThkS2lscTJsSWN6UFdxbEtMTVhhdTZ1Wkg3TXp2ZzdNZE9mNlh1UGxFd2Y2ejF0K1JzLy85RnJ1WCtQZWRFY1FuM0xZaXhYWW91V0RKRWJ4Uk9zdlRZVW5NWlZBNXN3bXRDL3ZJZXk2RHFLSVJnWlRDSHdrTEE0aEZpMVBISzFjN0lsSVY4Wm9ibUttWUZWbmswdk40OGlyMC9Hd0RoSDdYSzFUWDdqakp3b08wVTJTZUlmcTFBVVJQU1ZYdEZOSUFMSFNIM2JPT2dDRUhvaHVCSXl2b21hcUJQd1U9Il0KCX0KfQ==\"\n[settings.host-containers.kubeadm-bootstrap]\nenabled = true\nsuperpowered = true\nsource = \"public.ecr.aws/eks-anywhere/bottlerocket-bootstrap:v1-25-10-eks-a-36\"\nuser-data = \"IyMgdGVtcGxhdGU6IGppbmphCiNjbG91ZC1jb25maWcKCndyaXRlX2ZpbGVzOgotICAgcGF0aDogL3Zhci9saWIva3ViZWFkbS9wa2kvY2EuY3J0CiAgICBvd25lcjogcm9vdDpyb290CiAgICBwZXJtaXNzaW9uczogJzA2NDAnCiAgICBjb250ZW50OiB8CiAgICAgIC0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQogICAgICBNSUlDNmpDQ0FkS2dBd0lCQWdJQkFEQU5CZ2txaGtpRzl3MEJBUXNGQURBVk1STXdFUVlEVlFRREV3cHJkV0psCiAgICAgIGNtNWxkR1Z6TUI0WERUSXpNRFl3TlRFMU1UQTBORm9YRFRNek1EWXdNakUxTVRVME5Gb3dGVEVUTUJFR0ExVUUKICAgICAgQXhNS2EzVmlaWEp1WlhSbGN6Q0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUpuNgogICAgICBRcysxc2tlVVdFUXdmb2pCanpCcUI2WFcyNTcraGhXakg4NnN0L0hPT2hDYmE4dnJTbFdjRU1jVE9OeXBYVmVoCiAgICAgIHdRNjBoZkRiaG5CVTZybkxuLzJnRWNmQ1d3dlhUOVA3TUIxc2Y3a3ZDeTN4WUtyY3pSVjVORXpHeUxkdzFMR3kKICAgICAgdmlYRk44V1B3WDh5Wkt3S1YyaW5nWWdvT3UxUHpwaFNXeVFVRzIxV24rMzR2L3MyT1hvTFYvNDkyVHJRUTV3aAogICAgICBCK2JJYkJaUjMxcVEvYXgvZkErWUxEeU9kVUJCTGI0MHBiakVqN0p3VXhNU1VQVVpEYWpHVFJJaGJuV3Y3aXNrCiAgICAgIHR2Mk0yWTRndHV0bGsrTzVUOEtaV0l6ejdwYldOYW5QcndnV3NxWTBURitybGRsejdVRlRaZzdNdElSVElqS2wKICAgICAgVjgwRG8yRVV4aUMzUmhWMVZLOENBd0VBQWFORk1FTXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01CSUdBMVVkRXdFQgogICAgICAvd1FJTUFZQkFmOENBUUF3SFFZRFZSME9CQllFRk5xSkJ1YVNkK0NtNjJJakZIY0lpYWlKUUlxaU1BMEdDU3FHCiAgICAgIFNJYjNEUUVCQ3dVQUE0SUJBUUFPUUEzeUdmcU5sWTlJc3FPT2EyMmhYU3hEc0JBY0hYRTZwbFRHckdJMk5HY2UKICAgICAgSDlnK2FCYUZUYnc5eFZTSG9LYzRjamlYSGJaYXByaWkzRVloR2NTeEtLcG11QUFacStuM0tscWJGaUpsYndwUQogICAgICA2NlBWWFd2VStpUGpCeUtCNGR3eDBOK2xlcTVIemVCeUpPUERSUlVQZnJ2ZXV1UE01VHYwVi9IWmNhS0ZxN0xuCiAgICAgIE1aTUtXbEdHbVBXM0FTNUoyZ3JNejYxaEhJbXZ6WE5LUEd5QVpKWFpnNjdFQzRuTmI5RVJpMGdaKzVZNHpJSGIKICAgICAgbGlYMXBSMFErMFRZS3BjQmZkcFRlY0UyY3lVVEpXY3Q2dSt6VVV4L3ByMlhIZVZpbWhzUThpT1BzYmdJV2R0QwogICAgICBDeVRkMnVkY3E2RDhpUnhNZENMZjloK1VvNUtIN29vdDZZSVhMZDlHCiAgICAgIC0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KICAgICAgCi0gICBwYXRoOiAvdmFyL2xpYi9rdWJlYWRtL3BraS9jYS5rZXkKICAgIG93bmVyOiByb290OnJvb3QKICAgIHBlcm1pc3Npb25zOiAnMDYwMCcKICAgIGNvbnRlbnQ6IHwKICAgICAgLS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQogICAgICBNSUlFcEFJQkFBS0NBUUVBbWZwQ3o3V3lSNVJZUkRCK2lNR1BNR29IcGRiYm52NkdGYU1menF5MzhjNDZFSnRyCiAgICAgIHkrdEtWWndReHhNNDNLbGRWNkhCRHJTRjhOdUdjRlRxdWN1Zi9hQVJ4OEpiQzlkUDAvc3dIV3gvdVM4TExmRmcKICAgICAgcXR6TkZYazBUTWJJdDNEVXNiSytKY1UzeFkvQmZ6SmtyQXBYYUtlQmlDZzY3VS9PbUZKYkpCUWJiVmFmN2ZpLwogICAgICArelk1ZWd0WC9qM1pPdEJEbkNFSDVzaHNGbEhmV3BEOXJIOThENWdzUEk1MVFFRXR2alNsdU1TUHNuQlRFeEpRCiAgICAgIDlSa05xTVpORWlGdWRhL3VLeVMyL1l6WmppQzI2MldUNDdsUHdwbFlqUFB1bHRZMXFjK3ZDQmF5cGpSTVg2dVYKICAgICAgMlhQdFFWTm1Ec3kwaEZNaU1xVlh6UU9qWVJUR0lMZEdGWFZVcndJREFRQUJBb0lCQUd2bEhKbjF4UFk1Y2R4ZQogICAgICBhdHdqWGRYY0JuM2ozOE02c1BSS2VRaFdxUktNb3MxdVN3djZWUDFHUGM5eG5iK3FBaGxjNUM0TXUycDZPV3RQCiAgICAgIFphZEprMU9pcmJMZmN2TUtMZ3Jwa2Q2Y1g2ZUEyb0dZQ1ZmUkh2ZzBGVHpUY21Sd2NPb1B5UVFvZlJzM3o2djAKICAgICAgd094NDFxVWo4elZZazJvbEtTdzlVVlluck14Zzk2ZWIrN1JjSGU1NURySWhJQ1cyb2JWbDlwNUkyOHVwMTNRRgogICAgICAxUm5Band0V2F6M3pzUjU1MXhnam00bk9HMm42WEhxamFmQ3RueFFvcTFLV21Vczd6bkFiN0VvS283TWYwaUZkCiAgICAgIDhGdjE0YVJsc3pDc0diaGttZmdGSGxnWGkwNURNem9RQm11MzUxemk4NC9CSysySzJ5QmlyUVhXcmZmVEJ6TGcKICAgICAgakJXR3dKa0NnWUVBeWlxNnBOZkZjbWNaUG5XenhlZlZnUnlOSWVseTNheXBTK3VqRVdwRHJqamRVM3dyV2daRwogICAgICBBYnFNMllBbHQyVlRHalVUbTlCd2drT250YWY3Q2l1L1BSQ3ZNc2NjaTNlc3MvQys5c0p0OU9xeDQ1b09sQnMyCiAgICAgIFpCRDRqcGt2aUQ1VnYrNHdrckl0YUJwRU1FQk9uNFMvekpQU1VoSUpyZDdrUjZaSGRrQ3MwZ1VDZ1lFQXd2cVMKICAgICAgYjUvWDJYTXpsVUQ4Umk2QVBBR3g4Zk9BN1FQMnlyQmlWMUNPZUYxa1ErQWozN1RDdy9XNThHb1hURzRFTWdoagogICAgICAxbnZyVXNsajYyQ0hXeW9LeFhrbGRVakJ6bTZPVnJ4VHRldXV2WGpwY3dHOVE4SWZzbVg5SXpQZEV0QzVTUnpICiAgICAgIGlBMVYxVHphcVhnUGhLaDhUV0UvdXJqU0dpRnFVQ3BhcmM4dWhpTUNnWUVBeFBLTmMydDB3YmVvZ3cyZFBjNVYKICAgICAgVTN0aURraGppNHJhUHNqbXlsOXdZYmlwL051NVMvRlNuL3FCbnAzVm9HMUlZUDZXQkxReDl1VTc2NThpMDh5OAogICAgICBlQnZaNGFqUnFSakVHV1FPVlV3aVhIZUxKd1I4OFZIMVVkU3FvQmloa3FQUFc3UUtnODZxcDREM0x0NW0rY1lVCiAgICAgIEo1TldVSGVjRUZOVXBteFpyOXpmdjVVQ2dZRUF1MWE0aG9vRmdnZy9Zc1A2NEkvalpFU1lyZ215TVlraWdlTngKICAgICAgeWNVNzdvaUZRdlpFWWJnemZzZEdMYW94MHB6T2FTaElqUmVwcG5TY0RkZEVscUpSa1NWeWlUc3NBK1dUMitDOQogICAgICBhY2tXcnpSUzBjNjFCRHFyNitRMGtiTk9VYnE4bkhRTGZ6eVk0UGJFZmhvK2hzN0FDRFZOWDJJZmRUM3dBVENBCiAgICAgIHlnbU1BUk1DZ1lBOEl0a1hpbWl4M2l6MkkvSytzeklIWXNrRDNsblhQRDVML3JOcmFhS3crUzcyWVZuN1lDYkkKICAgICAgR2hSRG1aTC9OTDZJZnAxeG5NODA0ZXFldExBSTlaVmpHcE9CZlVwaGN4SjMrRTJNbUVDTmVTcW1zMExMQ2phVAogICAgICBjU2NScXltbVVFWlAwamJaZFN4cGg5Y2ppMVFwSHVZWEN1aHVvNlo2ZTNtRmgxS1dsYXFLTGc9PQogICAgICAtLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQogICAgICAKLSAgIHBhdGg6IC92YXIvbGliL2t1YmVhZG0vcGtpL2V0Y2QvY2EuY3J0CiAgICBvd25lcjogcm9vdDpyb290CiAgICBwZXJtaXNzaW9uczogJzA2NDAnCiAgICBjb250ZW50OiB8CiAgICAgIC0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQogICAgICBNSUlDNmpDQ0FkS2dBd0lCQWdJQkFEQU5CZ2txaGtpRzl3MEJBUXNGQURBVk1STXdFUVlEVlFRREV3cHJkV0psCiAgICAgIGNtNWxkR1Z6TUI0WERUSXpNRFl3TlRFMU1UQTBOVm9YRFRNek1EWXdNakUxTVRVME5Wb3dGVEVUTUJFR0ExVUUKICAgICAgQXhNS2EzVmlaWEp1WlhSbGN6Q0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUxQKwogICAgICA0Rmo4WjVBcElFWE4zTDNqbThTZmpYTmxCbWd4LzRtOTJsSU1kdEFsSGVsNVhkbURRcm1IUVZiN09lZmlicU9pCiAgICAgIHJnRnB3T0ZmU21vR0ptUXJZK0Q1cUE5c3RLd3NITG5XT2N3M0VleWRLYlptYXJmaXExQ2RwU216V1A5aS9ZTkMKICAgICAgVm9VZEhHVjEzcXRXUU4wUzRnT1hQMnJjR3hwYm03THB2T3dHWFFEdmh4WFN1RHJQVC94ZWc1NWVOSmhPUlU5VAogICAgICBTYTJuMEpoTU9BcXdtTXVybElidmtseXdJWDhaWnBsUFVwK3ZKYWhTb09LVWJQOUlUd1gxS0t0cWZ4dStiZkMxCiAgICAgIDh5dXNmQ3pIaDhwQTRVcUtKZXI3d2JRMDRWblI5Z3VCTG83SDhZYWlyazFGQm5OSlZkNnZuUysyVERBYnZpRXMKICAgICAgMzU4VERrbjZVa1B5UWFMR2h1OENBd0VBQWFORk1FTXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01CSUdBMVVkRXdFQgogICAgICAvd1FJTUFZQkFmOENBUUF3SFFZRFZSME9CQllFRkJCZWZsNGFRdXFvVnF3TFdPa0FlZFhDZ0phMk1BMEdDU3FHCiAgICAgIFNJYjNEUUVCQ3dVQUE0SUJBUUFFb2RtRzNtZVE2RVFQMEdFbnhnYlBhQmRjaFBjTzBHMGZRYncyZ0I5QzlmY1MKICAgICAgK2FVTnBYZVRHNjVPRjVuWmZQbUhyZUhZdGY0V253TWhlKzFUelhLVlFXL0MxRkhJRFBSZHdmbk83TXkveHNudgogICAgICBIMDNNNDk1MlM1K3BNQXpuN0k2NzVXVktXa1NTT1B2bjBEd2pUMS9veXJNL2g0QkQ5Q1dsa3lIUWVnRlV3Nk9WCiAgICAgIFEyWXZ1TGw4TnVjVVo5QzVPeG9HN3FEQWdoc2h5NG14bTdjbkxpWUxPSUtqM056eGo5bVlZZDJIcmxqWWRwbzYKICAgICAgSEwrVkhYbVhaWG5ndDFEaFRiQ2QxdDFxYzhpSFBGb3R4MHlVRkhBWGZSQWJRQzc1ZGxCcGdzV1Jtb3Arb3dHVQogICAgICBYYWF1MGtVQ0ZvczRseE9zbkQrb2VqRTN2eFRRdkNLR0U4emRoMTkwCiAgICAgIC0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KICAgICAgCi0gICBwYXRoOiAvdmFyL2xpYi9rdWJlYWRtL3BraS9ldGNkL2NhLmtleQogICAgb3duZXI6IHJvb3Q6cm9vdAogICAgcGVybWlzc2lvbnM6ICcwNjAwJwogICAgY29udGVudDogfAogICAgICAtLS0tLUJFR0lOIFJTQSBQUklWQVRFIEtFWS0tLS0tCiAgICAgIE1JSUVvZ0lCQUFLQ0FRRUFzLzdnV1B4bmtDa2dSYzNjdmVPYnhKK05jMlVHYURIL2liM2FVZ3gyMENVZDZYbGQKICAgICAgMllOQ3VZZEJWdnM1NStKdW82S3VBV25BNFY5S2FnWW1aQ3RqNFBtb0QyeTByQ3djdWRZNXpEY1I3SjBwdG1acQogICAgICB0K0tyVUoybEtiTlkvMkw5ZzBKV2hSMGNaWFhlcTFaQTNSTGlBNWMvYXR3YkdsdWJzdW04N0FaZEFPK0hGZEs0CiAgICAgIE9zOVAvRjZEbmw0MG1FNUZUMU5KcmFmUW1FdzRDckNZeTZ1VWh1K1NYTEFoZnhsbW1VOVNuNjhscUZLZzRwUnMKICAgICAgLzBoUEJmVW9xMnAvRzc1dDhMWHpLNng4TE1lSHlrRGhTb29sNnZ2QnREVGhXZEgyQzRFdWpzZnhocUt1VFVVRwogICAgICBjMGxWM3ErZEw3Wk1NQnUrSVN6Zm54TU9TZnBTUS9KQm9zYUc3d0lEQVFBQkFvSUJBQXdBaU5FbU0zbW5aR0dECiAgICAgIE1wN3UyR29xQVhLYVgreit1SDFrelhQL0pNVUlWSkJHNlgwaGhZZDlmMXZmY0tNSHREenhFZzFDRWZ4MU1DOFYKICAgICAgcmVXc0J3THZZc2QwaGkrYzJXV0M1V082b0ZWWXZYbE9KdnVQLzhCbVlxYi9NSVlIQWNTSTNlQU9hdUJSRmNIegogICAgICBCR3c1WUJ6Qkx4ZGZDZWVlQ0NuTDhDOEUxWnJFUVdCMENCZmpnbzNmY2hCS2l1NVBTL2JSVnJ3TmJ2eFduYXZYCiAgICAgIGJIRVRaVkYyenpFdzJqeXZHZTZWVW5JT05pNHJCUDhjWlFLN3JwVmtCVFVpSVNTVnJIV3BjVzNyNnVuVVdtRXUKICAgICAgV2lWV0FFcUpqUnlBc0k0cXhBUmlCMTdYd1dpTEx1UnJuTnVvQkF5VEtIZnRlOS8xWWlFU2lCdk8rS25lRVNOSQogICAgICBFZWNETUhFQ2dZRUF5Mm0zbFVxMWcrUDkvaTk2bFZ4a1V5eVluVVZGM3UxUHBkUllrSnRha2pQbnFUUGVKNXhuCiAgICAgIG5kWXk2RkdrTldLMC9OMFRFNWtLSTJnSGJwd000TnFYWXBDOERXditLeTEzTGU1RWdzaGpWRjk0dndjZkxtSnEKICAgICAgczBzRHhaZUdqQXZSS2JHOEVRVTJ3ZHl6a1l0V1RlLzFZTnZoYkNwNGN2aE9hbkh5U2dNa1RrMENnWUVBNG9kYQogICAgICBCVXpxckFkamREbkhibkhMbGJBUlptRlh6cVpTelZvL3BVWThQNHpJMXNQUzFHdXlOMy9KU1dzdGZOZzduK25GCiAgICAgIC85STU0Qm4zNG9mWStWMTFQMjhZaDdSYXV3dE9zdEFjQ3Q5bmc1bXltVmdKVmNXbUtveVZqMWdWTmlEcmo2VGcKICAgICAgMzNJOVBoZlN2a0JTTGZmbmNXVmJ5TmpVN1F5ZXVYWk8zWXpTNENzQ2dZQjFZRDV3SmxrUmp5a01XT2RhY3FMdAogICAgICBuOGs4enpGZlR2N3J5Tm1HTUM4V09GOVFNcjdaaXBYNzdSTVpIYXNzcHhXYTZCTE85enR6Yk44RkE0VW01dHYrCiAgICAgIEkxaHdRa1c3TXBRWDYrcWFzUGtvUWFNU1VCQzlHa3RKeEZxYjFURHRkUkF3Q0FCbXJlU0gvMHViQzVVMGllZkYKICAgICAgQ1h4TmgwQlR6MWFvYzdJRTVVVTQ3UUtCZ0hTZGpVQXhTcFhvNzlBRGRxRnF2NDE5cUZkMlFVZkc0OVdIWWtCcAogICAgICBHZGIxV09jR3hHQktXT0t0VENnWm5yOG9hZmwyMVZGUEhqQTU3aHlXSnFLbzlCVUYwakQ2TGNNZ25SRDhoWk1yCiAgICAgIFV4U3lhUGo4RTBJdWo4NVR0U0tvQzdOajJ5Q0ZscVl4SDBuNTVhbS9YdzcvRWd5VVMxM29FaVUrVysvSjhldW4KICAgICAgOVY4bEFvR0FLclpkUE1Lb01uYTVGWWRoeFNXaTRrMnJ0WVVtRER4MUM3QkRSTGRqWTYrZ2lWa2doZ29HU0RLSAogICAgICB5cUtZNjkzaWFOSStIOTQ1RTE3eng4bG9ySnFHOXlRTWZYM0tvNFU4eWFhenpFempzbDY2Ti8rdTNGajZOc2JhCiAgICAgIGZlZGYyVXdkUC9rK2E3WnFxd1ZPZERUN3FTck9sOHZZeXJmTXZIRS85OThYQ01zaGdjST0KICAgICAgLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0KICAgICAgCi0gICBwYXRoOiAvdmFyL2xpYi9rdWJlYWRtL3BraS9mcm9udC1wcm94eS1jYS5jcnQKICAgIG93bmVyOiByb290OnJvb3QKICAgIHBlcm1pc3Npb25zOiAnMDY0MCcKICAgIGNvbnRlbnQ6IHwKICAgICAgLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCiAgICAgIE1JSUM2akNDQWRLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKICAgICAgY201bGRHVnpNQjRYRFRJek1EWXdOVEUxTVRBME5Gb1hEVE16TURZd01qRTFNVFUwTkZvd0ZURVRNQkVHQTFVRQogICAgICBBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTlE4CiAgICAgIGtYRHREZEFwS2JvTjdoVXgvSGJTU0pXU09oUU9tZFJUU2xweFV4M1JvZG0valk4eEhCbG4wbzMwMFNKQWZreUUKICAgICAgRy9RdCtVcjB5YjB3RVRURHJIZnFsWTRuTFBoTjd0NlMzdllncUhQbEtVQVVYTHNGOWVPRnhxSXZRZm03SXloQgogICAgICBoNkNBZEhsM3hacCtSNmF3TTQ2L3g0WklDVk1CTytQRVRhaEYvVWFKV0VIMW85MUpvL2JuRnlUWkE2dnFrMFB6CiAgICAgIGpHR0RYckUwODNsT2QzMDZIMXJwOFVVQ3ZnN1JtNkFWWnVhZW0zeUJIY2xsRjZLWnlFSWE1LzF4OFV5ekhtaGUKICAgICAgT2UwRzVxSzdMSDRuMStRUGY1c0dxVUVsZ0NYeGhPR0F0a3NRaS92WWN5ZnFyTVhpOThnVXRJZE1FSWQ5ZXUyQQogICAgICBOczlFZ25FaXF0TjRYcWpJZkFVQ0F3RUFBYU5GTUVNd0RnWURWUjBQQVFIL0JBUURBZ0trTUJJR0ExVWRFd0VCCiAgICAgIC93UUlNQVlCQWY4Q0FRQXdIUVlEVlIwT0JCWUVGQnE5a2c1NVpRaHNZblZ2QVlPZ0o1Qk01anppTUEwR0NTcUcKICAgICAgU0liM0RRRUJDd1VBQTRJQkFRRExoUEJ5Y0dDYkRpWVNnM2dsSmJlZFFoZ3p3cTNhRjYzd2xBUTU1WS9iYmZrZgogICAgICBheDdVQmVNdlZZeFI5Ynh0RnUrQllMNkVBODdrV1NDT1dDVlhSSi9CZVIrNnFmN2lwUHh6T1hPVmNIUHdERC9NCiAgICAgIHArUkx4aTRmcnBRUDc3Um8xN21yc25NVklXNTlNSGlCMWJGdG83emlNWWVoaHFOTGkzcmhOQ3FMcDRtbDljSnUKICAgICAgbFN0TnJSeTdwY0ZzaFJad20vWk0xeFhUSFpyODVuY09yWlc1K1h2bWZLQW0rZldtU3MvVlZ4VkFmakx0N2hNTAogICAgICA1M0VWaFRHUHlIUWNMRkd2cERrUEt3aFdSTVNNbngwWnZkM3Ixb096cmNFNFcwdGNURXlhS3ZCTHBMczFPYWtUCiAgICAgIG55djhpN1ZnOHRFZWhjQkRaaGZWOThma0MxdElabHJzREpjZEU0QS8KICAgICAgLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQogICAgICAKLSAgIHBhdGg6IC92YXIvbGliL2t1YmVhZG0vcGtpL2Zyb250LXByb3h5LWNhLmtleQogICAgb3duZXI6IHJvb3Q6cm9vdAogICAgcGVybWlzc2lvbnM6ICcwNjAwJwogICAgY29udGVudDogfAogICAgICAtLS0tLUJFR0lOIFJTQSBQUklWQVRFIEtFWS0tLS0tCiAgICAgIE1JSUVwUUlCQUFLQ0FRRUExRHlSY08wTjBDa3B1ZzN1RlRIOGR0SklsWkk2RkE2WjFGTktXbkZUSGRHaDJiK04KICAgICAganpFY0dXZlNqZlRSSWtCK1RJUWI5QzM1U3ZUSnZUQVJOTU9zZCtxVmppY3MrRTN1M3BMZTlpQ29jK1VwUUJSYwogICAgICB1d1gxNDRYR29pOUIrYnNqS0VHSG9JQjBlWGZGbW41SHByQXpqci9IaGtnSlV3RTc0OFJOcUVYOVJvbFlRZldqCiAgICAgIDNVbWo5dWNYSk5rRHErcVRRL09NWVlOZXNUVHplVTUzZlRvZld1bnhSUUsrRHRHYm9CVm01cDZiZklFZHlXVVgKICAgICAgb3BuSVFocm4vWEh4VExNZWFGNDU3UWJtb3Jzc2ZpZlg1QTkvbXdhcFFTV0FKZkdFNFlDMlN4Q0wrOWh6SitxcwogICAgICB4ZUwzeUJTMGgwd1FoMzE2N1lBMnowU0NjU0txMDNoZXFNaDhCUUlEQVFBQkFvSUJBQlNDU2Rlc3dLSDE2RmRYCiAgICAgIFUwTnZFaU4xVEpoUGloYzRGZnRqbFJMS1VxcldBQnJubk1sU2lKR1g2OHZWOVlPbTBjMlpFbzUxQnRzTWJwWSsKICAgICAgbDlzT2NaTWc4eFRLaWxqd1J2M3hHV1NWWVZIWnVqTzhBLzM2cEhrNUN5blBBVVFkcGxjVWhnT25oaG5hemhpUQogICAgICB6V0c1TXpJL0xBdTQyRlhTRDdTQjJyTkJHQ21CblBFaldGckxmMWUxOTlnOVBkWC9GRm1sWEo4RDFHL2I3R1ZwCiAgICAgIEU0VGFLUlJTaWFTSEs0UzRqRm9DNFY4ZlFVVEsxbUZIcnh4T1RBSDYrQTVTbjRzME9YRWZoVys2aFM4UU1QaU8KICAgICAgWFJXNC9SNFlIbDJtQTV0enRrQU00aXMwbnNUdndCUlc3akp2Y3RFVG1zUmZPUmc5M2xFZXBLMk5UUTVoekNVcAogICAgICBVakxBVzJFQ2dZRUEvcGl3UTRnWXo3dHg0UUEzekczY3BVajhna3RnRDd0dVFWazQzZ3dvdnZlUE5tQlloTUtLCiAgICAgIEtzOXVkYytNMWJ4cURFVW1JampaenFzRURTb0FBVWNtU3NCZTdXa1VkVW03NkFiekFMcjVRNFBsMlJ5Y3BHUGgKICAgICAgU1BBMld6K2phZXlmdXRmUUNOemJjMkZYZEJtemNrM3NEWlU4RXRFQVp0bk9SMlZseXFnT3laMENnWUVBMVdnWQogICAgICA1QkQyWHg0dHpUeU9pWWpNMyt5V0ExSVJPSlFFc2RRdzQwdzNBMlV3WGpOQlV1UG5XaE1aTTk1OStESlNYTDFICiAgICAgIHk2Z0dpb2xZcVluT2w4anN1Wkc1djVjUnB5cndvZ1BkbE5LVUNjS1lBSjZZNVdaOTlHY3pFQ01Ldm5tRGRqZWwKICAgICAgeTZ3Qk0zUHd6Q3dIdVBkYlBNd0tMcE9DWUJTd0ZsZzVweDY4dzRrQ2dZRUFwNEh3WHM0NWZOdVlKbkNOUmN4MgogICAgICBvcXp0cmhCSHFMSXA3WWIxZW1ySG1EV3JIUnl3d05CNk5ZWjY1N3Boci9LaVYyWmJtNktKODRialNJSDh0TnFLCiAgICAgIElCNkhsbTVQam9ldndRNXBiVzZYTjh1ZE80YXVyUjRtQ0dZN2JUZm1uWGVZOUVhdjBsVDFjZWwycjZXRlFreHcKICAgICAgWmROdFRmZ0M2cWlkSnE1WkZjZ1N5cVVDZ1lFQXRTQ05QOFZGMXFWK3FsdGpmMGdrMjBtcWFWY1dWcmNLNVFQOAogICAgICBHbTl0b2V3WWlWdG5ianNRK1ZxTVlZSE4yUUtjOVNtUjdrREdqSDdXU2M4MUVZN3ZuUEVhZm9weDZUaUExSUlECiAgICAgIFozVHpRUFZ3bmRYK1gzWUdJWklWdlBTQkFmbVFvcDNJa24yQzUzRFlSL0oyKzM2MmFYdWtpTE9hVElKQ2tqUi8KICAgICAgZW1DVVJrRUNnWUVBaHRUSVNEejBsWk9zNmRSSWxoZWlaTzB2citiM0FOOFYrUXMzR1NLNWdiQ1RGRWFrNzFEUgogICAgICBIQ1ZqSW9ScDJEZXdsK3lrS2I5a0RpdHFPd29halBzOWlyUUZVa2FDV0xEZXBWVkJLRmtJWXo5U01SWExSMVNVCiAgICAgIFNiTWNHYS9QaXZvc2p1d1RNaGRoUk9TZ3dmTnZuUWZUeWNzWW1vamNvQmNwdndBdW91YW1JRms9CiAgICAgIC0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCiAgICAgIAotICAgcGF0aDogL3Zhci9saWIva3ViZWFkbS9wa2kvc2EucHViCiAgICBvd25lcjogcm9vdDpyb290CiAgICBwZXJtaXNzaW9uczogJzA2NDAnCiAgICBjb250ZW50OiB8CiAgICAgIC0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tCiAgICAgIE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNS9VNTFYakU2RCtJaEZuNGt5NGoKICAgICAgL0ZacjdBdUx5cDl0WkFXQXhHUmZrYnFJbDFNcGxva2JPWDZicHl3ZnRTSnJ6Q3V0dDM4Kzhsc1pPelFKSUxxLwogICAgICBrdzFhM0pWYVVqdW5wcWhoVUFoTkEyaFZ6Y2lodkRoUytFWUhySmtzdkcxR01iN3k2VnNXcnNvUG41VjRkQmx6CiAgICAgIDFMZlRMQjhuL0VVQlMyNjY4U1d5L2w1YkQ0MHdlRGgvamREUCtwbitEUDkrZnV1QmZMaHh3ZkxyU1ZRRWYzU00KICAgICAgZVo5NU1JZDE5VWZWMlptc2FJTHhrVmRlRmM1dHI4anhYODJIcDA2MVJwd05BaTRQNWtSY1YvZUNGUHVFUnJQYgogICAgICB1Ymw1VFN6Z3lUdG84cXIrREg0RkxNUUNEZFBSbVBPR2pKQWRSS2FFTGpGeitpVU5JcVo0TWpaRm9yNHYyenZJCiAgICAgIEFRSURBUUFCCiAgICAgIC0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQogICAgICAKLSAgIHBhdGg6IC92YXIvbGliL2t1YmVhZG0vcGtpL3NhLmtleQogICAgb3duZXI6IHJvb3Q6cm9vdAogICAgcGVybWlzc2lvbnM6ICcwNjAwJwogICAgY29udGVudDogfAogICAgICAtLS0tLUJFR0lOIFJTQSBQUklWQVRFIEtFWS0tLS0tCiAgICAgIE1JSUVwZ0lCQUFLQ0FRRUE1L1U1MVhqRTZEK0loRm40a3k0ai9GWnI3QXVMeXA5dFpBV0F4R1Jma2JxSWwxTXAKICAgICAgbG9rYk9YNmJweXdmdFNKcnpDdXR0MzgrOGxzWk96UUpJTHEva3cxYTNKVmFVanVucHFoaFVBaE5BMmhWemNpaAogICAgICB2RGhTK0VZSHJKa3N2RzFHTWI3eTZWc1dyc29QbjVWNGRCbHoxTGZUTEI4bi9FVUJTMjY2OFNXeS9sNWJENDB3CiAgICAgIGVEaC9qZERQK3BuK0RQOStmdXVCZkxoeHdmTHJTVlFFZjNTTWVaOTVNSWQxOVVmVjJabXNhSUx4a1ZkZUZjNXQKICAgICAgcjhqeFg4MkhwMDYxUnB3TkFpNFA1a1JjVi9lQ0ZQdUVSclBidWJsNVRTemd5VHRvOHFyK0RINEZMTVFDRGRQUgogICAgICBtUE9HakpBZFJLYUVMakZ6K2lVTklxWjRNalpGb3I0djJ6dklBUUlEQVFBQkFvSUJBUUM3bE9YalM1bzVrMytNCiAgICAgIFFOSXovQ0ZmNUdlOGFRM3dtNE0wV3ZycVY1MnQxU0syOWFyeE1RbVNNbUFnRGgvS05QN21Dd0NlSDBwQlpnaCsKICAgICAgaHpOR2c1OS9oVkpRaG51WGV1UzJjdjdYWVE4ZXpWWnVaMnpjTU5Sd01QbnR5NldRNy9IUE12TndZWmh6VzdiTQogICAgICB3R3k3dndXY0pkaWhtc1NVVHgyZjZmbEdJTTJpeEJNMWR2bUs1bytYRG1QOVl2MW9jamh5NHFkdUdEWDltUVdXCiAgICAgIGVNOWdabUZ5NnRoUDhheXFqbEs0TnZxcG9xRnN1RTF0Ukk2Y2pIYkVFNGxyVWlPdS96RXIyaXRwVzdzTVU2b04KICAgICAgbnZjUmx1V3BNRlRqN0gyVkJrSUhoVkJCd3VUT2NvQUQ5WCthWXBxUDg1T3ZtL1UwN3NRdkJLVXlKUUhGdEFkaQogICAgICBKbVVzZUp6UkFvR0JBUGxYYkJicC9QdGFZdkM5MndSZE92Q3ZLWDRjWUxwSlJyUjk4aU5YRm54bERNOHo2Qk1FCiAgICAgIHdmT3NCWVBwQk05UVd0NXZCdDRHdHFNbjhjZUYza0V6ZjY1NmxLRWg5SDZQRE04MjhvaU53VEYyRlVWc2JXYXYKICAgICAgMlQ2NldudDZ5RldHb25vbFNLWVJrMy9zRXNiT2tWK05pZGpzUzdaSmNQbjFDNEFvaFUwMzFyVVZBb0dCQU80bQogICAgICA5c2xhcnVqUUNWOTh0QzQvWWlxZkNnYThEcEVITWVRaHFlak5rYXBaTU1ETkhlcmtPR1pRQlNUVU1MUkJ6Sk9rCiAgICAgIEtWNitZZURINEdaUkxKbzZkcldNQmJGRnZSZHNuY29GdGY3cG4zQXJoMWllNFdOTmRhQTRUdGw0YXVwZE94MVoKICAgICAgcGFDcDVXZVhrN3BUanJaaStOL2tObm5GdkxQbEpCbXhSMS9rbEpvOUFvR0JBTVZ0aDFTTmFaYk1kdE1RUVQ4dgogICAgICBZdC80a2U0ZElpbmVvM0YrMkI5TGNhNkZoS0w4QXFJc2ZqWW0yeWNiZG9lQXBMTERUcmkyc0I1NEhtVlJoaTR5CiAgICAgIFRNTW1wRkVCeGNvQUVyQndYWkhxVERLUndUMzdJSlRTWUQzZTZJNGxKa015RzZ5RG9RWjluRUVKOThRYkE5aVoKICAgICAgQmJFUlNOSEpUUDllSEFFYUZKS1R5Qm9SQW9HQkFMNE1GQWF0T2tXSnR1RWZkLzRzRUorb21PeTA1Lzd1T2U4dQogICAgICB1aE9ROEx4N1BuK3RjRUdCYkV5aGNPbHA5NC94cmxybnR5Zm5UOTU4UXVRRHhVOHlkb2I4UFpLdzcyd2cvbTQ0CiAgICAgIFRuc2xYbG02TXVFU3NSUjR2UFJsMnU2S3ZPOVlCUk92OVkrWDVQemRKa09iNkpnOXRST2VYNmFmbUs0S250dHQKICAgICAgOEdKaTIvK1ZBb0dCQUkvZEVlS0RrZVdITnREMmVsOFlUTXJtbFNCREZDMVBuQWNnRkJ6cytmdVMybkljN2ZsWQogICAgICBaSGE1NVB0RGYzMUx3SjhsSFFWaGpGc0lHVmE4WG5UZHZMMGhqUThLTXhET3lGNDFGcmhKMERLaGs0aTc4dEpkCiAgICAgIE9uNS95SU5MVlJXeDFuNlNrd1hyV3FjRzhhOGp3cSt5YXV0SXJROXprSGVHcHJNTFUrUW9kb29MCiAgICAgIC0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCiAgICAgIAotICAgcGF0aDogL2V0Yy9rdWJlcm5ldGVzL21hbmlmZXN0cy9rdWJlLXZpcC55YW1sCiAgICBvd25lcjogcm9vdDpyb290CiAgICBjb250ZW50OiB8CiAgICAgIGFwaVZlcnNpb246IHYxCiAgICAgIGtpbmQ6IFBvZAogICAgICBtZXRhZGF0YToKICAgICAgICBjcmVhdGlvblRpbWVzdGFtcDogbnVsbAogICAgICAgIG5hbWU6IGt1YmUtdmlwCiAgICAgICAgbmFtZXNwYWNlOiBrdWJlLXN5c3RlbQogICAgICBzcGVjOgogICAgICAgIGNvbnRhaW5lcnM6CiAgICAgICAgLSBhcmdzOgogICAgICAgICAgLSBtYW5hZ2VyCiAgICAgICAgICBlbnY6CiAgICAgICAgICAtIG5hbWU6IHZpcF9hcnAKICAgICAgICAgICAgdmFsdWU6ICJ0cnVlIgogICAgICAgICAgLSBuYW1lOiBwb3J0CiAgICAgICAgICAgIHZhbHVlOiAiNjQ0MyIKICAgICAgICAgIC0gbmFtZTogdmlwX2NpZHIKICAgICAgICAgICAgdmFsdWU6ICIzMiIKICAgICAgICAgIC0gbmFtZTogY3BfZW5hYmxlCiAgICAgICAgICAgIHZhbHVlOiAidHJ1ZSIKICAgICAgICAgIC0gbmFtZTogY3BfbmFtZXNwYWNlCiAgICAgICAgICAgIHZhbHVlOiBrdWJlLXN5c3RlbQogICAgICAgICAgLSBuYW1lOiB2aXBfZGRucwogICAgICAgICAgICB2YWx1ZTogImZhbHNlIgogICAgICAgICAgLSBuYW1lOiB2aXBfbGVhZGVyZWxlY3Rpb24KICAgICAgICAgICAgdmFsdWU6ICJ0cnVlIgogICAgICAgICAgLSBuYW1lOiB2aXBfbGVhc2VkdXJhdGlvbgogICAgICAgICAgICB2YWx1ZTogIjE1IgogICAgICAgICAgLSBuYW1lOiB2aXBfcmVuZXdkZWFkbGluZQogICAgICAgICAgICB2YWx1ZTogIjEwIgogICAgICAgICAgLSBuYW1lOiB2aXBfcmV0cnlwZXJpb2QKICAgICAgICAgICAgdmFsdWU6ICIyIgogICAgICAgICAgLSBuYW1lOiBhZGRyZXNzCiAgICAgICAgICAgIHZhbHVlOiAxOTIuMTY4LjEwLjg2CiAgICAgICAgICBpbWFnZTogcHVibGljLmVjci5hd3MvZWtzLWFueXdoZXJlL2t1YmUtdmlwL2t1YmUtdmlwOnYwLjUuNS1la3MtYS0zNgogICAgICAgICAgaW1hZ2VQdWxsUG9saWN5OiBJZk5vdFByZXNlbnQKICAgICAgICAgIG5hbWU6IGt1YmUtdmlwCiAgICAgICAgICByZXNvdXJjZXM6IHt9CiAgICAgICAgICBzZWN1cml0eUNvbnRleHQ6CiAgICAgICAgICAgIGNhcGFiaWxpdGllczoKICAgICAgICAgICAgICBhZGQ6CiAgICAgICAgICAgICAgLSBORVRfQURNSU4KICAgICAgICAgICAgICAtIE5FVF9SQVcKICAgICAgICAgIHZvbHVtZU1vdW50czoKICAgICAgICAgIC0gbW91bnRQYXRoOiAvZXRjL2t1YmVybmV0ZXMvYWRtaW4uY29uZgogICAgICAgICAgICBuYW1lOiBrdWJlY29uZmlnCiAgICAgICAgaG9zdE5ldHdvcms6IHRydWUKICAgICAgICB2b2x1bWVzOgogICAgICAgIC0gaG9zdFBhdGg6CiAgICAgICAgICAgIHBhdGg6IC92YXIvbGliL2t1YmVhZG0vYWRtaW4uY29uZgogICAgICAgICAgICB0eXBlOiBGaWxlCiAgICAgICAgICBuYW1lOiBrdWJlY29uZmlnCiAgICAgIHN0YXR1czoge30KICAgICAgCi0gICBwYXRoOiAvdG1wL2t1YmVhZG0ueWFtbAogICAgb3duZXI6IHJvb3Q6cm9vdAogICAgcGVybWlzc2lvbnM6ICcwNjQwJwogICAgY29udGVudDogfAogICAgICAtLS0KICAgICAgYXBpU2VydmVyOiB7fQogICAgICBhcGlWZXJzaW9uOiBrdWJlYWRtLms4cy5pby92MWJldGEzCiAgICAgIGJvdHRsZXJvY2tldEJvb3RzdHJhcDoKICAgICAgICBpbWFnZVJlcG9zaXRvcnk6IHB1YmxpYy5lY3IuYXdzL2Vrcy1hbnl3aGVyZS9ib3R0bGVyb2NrZXQtYm9vdHN0cmFwCiAgICAgICAgaW1hZ2VUYWc6IHYxLTI1LTEwLWVrcy1hLTM2CiAgICAgIGJvdHRsZXJvY2tldENvbnRyb2w6IHt9CiAgICAgIGNlcnRpZmljYXRlc0RpcjogL3Zhci9saWIva3ViZWFkbS9wa2kKICAgICAgY2x1c3Rlck5hbWU6IG1nbXQKICAgICAgY29udHJvbFBsYW5lRW5kcG9pbnQ6IDE5Mi4xNjguMTAuODY6NjQ0MwogICAgICBjb250cm9sbGVyTWFuYWdlcjoKICAgICAgICBleHRyYVZvbHVtZXM6CiAgICAgICAgLSBob3N0UGF0aDogL3Zhci9saWIva3ViZWFkbS9jb250cm9sbGVyLW1hbmFnZXIuY29uZgogICAgICAgICAgbW91bnRQYXRoOiAvZXRjL2t1YmVybmV0ZXMvY29udHJvbGxlci1tYW5hZ2VyLmNvbmYKICAgICAgICAgIG5hbWU6IGt1YmVjb25maWcKICAgICAgICAgIHBhdGhUeXBlOiBGaWxlCiAgICAgICAgICByZWFkT25seTogdHJ1ZQogICAgICBkbnM6CiAgICAgICAgaW1hZ2VSZXBvc2l0b3J5OiBwdWJsaWMuZWNyLmF3cy9la3MtZGlzdHJvL2NvcmVkbnMKICAgICAgICBpbWFnZVRhZzogdjEuOS4zLWVrcy0xLTI1LTEwCiAgICAgIGV0Y2Q6CiAgICAgICAgbG9jYWw6CiAgICAgICAgICBkYXRhRGlyOiAiIgogICAgICAgICAgaW1hZ2VSZXBvc2l0b3J5OiBwdWJsaWMuZWNyLmF3cy9la3MtZGlzdHJvL2V0Y2QtaW8KICAgICAgICAgIGltYWdlVGFnOiB2My41LjYtZWtzLTEtMjUtMTAKICAgICAgaW1hZ2VSZXBvc2l0b3J5OiBwdWJsaWMuZWNyLmF3cy9la3MtZGlzdHJvL2t1YmVybmV0ZXMKICAgICAga2luZDogQ2x1c3RlckNvbmZpZ3VyYXRpb24KICAgICAga3ViZXJuZXRlc1ZlcnNpb246IHYxLjI1LjgtZWtzLTEtMjUtMTAKICAgICAgbmV0d29ya2luZzoKICAgICAgICBwb2RTdWJuZXQ6IDE5Mi4xNjguMC4wLzE2CiAgICAgICAgc2VydmljZVN1Ym5ldDogMTAuOTYuMC4wLzEyCiAgICAgIHBhdXNlOgogICAgICAgIGltYWdlUmVwb3NpdG9yeTogcHVibGljLmVjci5hd3MvZWtzLWRpc3Ryby9rdWJlcm5ldGVzL3BhdXNlCiAgICAgICAgaW1hZ2VUYWc6IHYxLjI1LjgtZWtzLTEtMjUtMTAKICAgICAgcHJveHk6IHt9CiAgICAgIHJlZ2lzdHJ5TWlycm9yOiB7fQogICAgICBzY2hlZHVsZXI6CiAgICAgICAgZXh0cmFWb2x1bWVzOgogICAgICAgIC0gaG9zdFBhdGg6IC92YXIvbGliL2t1YmVhZG0vc2NoZWR1bGVyLmNvbmYKICAgICAgICAgIG1vdW50UGF0aDogL2V0Yy9rdWJlcm5ldGVzL3NjaGVkdWxlci5jb25mCiAgICAgICAgICBuYW1lOiBrdWJlY29uZmlnCiAgICAgICAgICBwYXRoVHlwZTogRmlsZQogICAgICAgICAgcmVhZE9ubHk6IHRydWUKICAgICAgCiAgICAgIC0tLQogICAgICBhcGlWZXJzaW9uOiBrdWJlYWRtLms4cy5pby92MWJldGEzCiAgICAgIGtpbmQ6IEluaXRDb25maWd1cmF0aW9uCiAgICAgIGxvY2FsQVBJRW5kcG9pbnQ6IHt9CiAgICAgIG5vZGVSZWdpc3RyYXRpb246CiAgICAgICAga3ViZWxldEV4dHJhQXJnczoKICAgICAgICAgIGFub255bW91cy1hdXRoOiAiZmFsc2UiCiAgICAgICAgICBwcm92aWRlci1pZDogUFJPVklERVJfSUQKICAgICAgICAgIHJlYWQtb25seS1wb3J0OiAiMCIKICAgICAgICAgIHRscy1jaXBoZXItc3VpdGVzOiBUTFNfRUNESEVfUlNBX1dJVEhfQUVTXzEyOF9HQ01fU0hBMjU2CiAgICAgICAgdGFpbnRzOiBudWxsCiAgICAgIApydW5jbWQ6ICJDb250cm9sUGxhbmVJbml0Igo=\"\n\n[settings.kubernetes]\ncluster-domain = \"cluster.local\"\nstandalone-mode = true\nauthentication-mode = \"tls\"\nserver-tls-bootstrap = false\npod-infra-container-image = \"public.ecr.aws/eks-distro/kubernetes/pause:v1.25.8-eks-1-25-10\"\nprovider-id = \"tinkerbell://eksa-system/eksabmhp1-cp-n-1\"\n\n[settings.network]\nhostname = \"mgmt-hsbj2\""
        }
    },
    {
        "Name": "eksabmhp1-dp-n-1",
        "Labels": {
            "type": "dp",
            "v1alpha1.tinkerbell.org/ownerName": "mgmt-md-0-1685978142429-ppwdp",
            "v1alpha1.tinkerbell.org/ownerNamespace": "eksa-system"
        },
        "Specs": {
            "disks": [
                {
                    "device": "/dev/sda"
                }
            ],
            "interfaces": [
                {
                    "dhcp": {
                        "arch": "x86_64",
                        "hostname": "eksabmhp1-dp-n-1",
                        "ip": {
                            "address": "192.168.10.97",
                            "family": 4,
                            "gateway": "192.168.10.1",
                            "netmask": "255.255.255.0"
                        },
                        "lease_time": 4294967294,
                        "mac": "08:00:27:13:27:db",
                        "name_servers": [
                            "8.8.8.8"
                        ],
                        "uefi": true
                    },
                    "netboot": {
                        "allowPXE": true,
                        "allowWorkflow": true
                    }
                }
            ],
            "metadata": {
                "facility": {
                    "facility_code": "onprem",
                    "plan_slug": "c2.medium.x86"
                },
                "instance": {
                    "allow_pxe": true,
                    "always_pxe": true,
                    "hostname": "eksabmhp1-dp-n-1",
                    "id": "08:00:27:13:27:db",
                    "ips": [
                        {
                            "address": "192.168.10.97",
                            "family": 4,
                            "gateway": "192.168.10.1",
                            "netmask": "255.255.255.0",
                            "public": true
                        }
                    ],
                    "operating_system": {},
                    "state": "provisioned"
                },
                "state": "in_use"
            },
            "userData": "\n[settings.host-containers.admin]\nenabled = true\nsuperpowered = true\nuser-data = \"CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FESktDaGllZEhqb2UwRDhZclAyUGFFQVR3bllYZ3JsanpnQ2gvSGIrcXNETXdYOU8vUE9IL0dod0V5WE1SSlZkdmo2eEJubXRyS1RRKzJneS9BcGxjelZvVTZkcnR0UzJ1Vk11aW1LZkNnZUFQOGthTHhwWUlNcWF6QkIvanZjK3paNUlScDJIcEMyamJ4ZndBSWdRVzNYaVdQcG9HbkFyYlpvLys2YVhMVmswRi9SclN2bUxEYlZGQTZKR0h1Qi9KY3pKaDBCZHB3RmxXR1JaWGRabEwzZUl1dmRyMXZjeWZPSStCU1RqbzRyRXc3TFEwN2drUFpjTStjVktWTE9YcEVWUThrTXduOHNHSnQxWXR3Qkw4UFgvWTU0TTRGL3JWTlVyWThkS2lscTJsSWN6UFdxbEtMTVhhdTZ1Wkg3TXp2ZzdNZE9mNlh1UGxFd2Y2ejF0K1JzLy85RnJ1WCtQZWRFY1FuM0xZaXhYWW91V0RKRWJ4Uk9zdlRZVW5NWlZBNXN3bXRDL3ZJZXk2RHFLSVJnWlRDSHdrTEE0aEZpMVBISzFjN0lsSVY4Wm9ibUttWUZWbmswdk40OGlyMC9Hd0RoSDdYSzFUWDdqakp3b08wVTJTZUlmcTFBVVJQU1ZYdEZOSUFMSFNIM2JPT2dDRUhvaHVCSXl2b21hcUJQd1U9Il0KCX0KfQ==\"\n[settings.host-containers.kubeadm-bootstrap]\nenabled = true\nsuperpowered = true\nsource = \"public.ecr.aws/eks-anywhere/bottlerocket-bootstrap:v1-25-10-eks-a-36\"\nuser-data = \"d3JpdGVfZmlsZXM6Ci0gICBwYXRoOiAvdG1wL2t1YmVhZG0tam9pbi1jb25maWcueWFtbAogICAgb3duZXI6IHJvb3Q6cm9vdAogICAgcGVybWlzc2lvbnM6ICcwNjQwJwogICAgY29udGVudDogfAogICAgICAtLS0KICAgICAgYXBpVmVyc2lvbjoga3ViZWFkbS5rOHMuaW8vdjFiZXRhMwogICAgICBib3R0bGVyb2NrZXRCb290c3RyYXA6CiAgICAgICAgaW1hZ2VSZXBvc2l0b3J5OiBwdWJsaWMuZWNyLmF3cy9la3MtYW55d2hlcmUvYm90dGxlcm9ja2V0LWJvb3RzdHJhcAogICAgICAgIGltYWdlVGFnOiB2MS0yNS0xMC1la3MtYS0zNgogICAgICBib3R0bGVyb2NrZXRDb250cm9sOiB7fQogICAgICBkaXNjb3Zlcnk6CiAgICAgICAgYm9vdHN0cmFwVG9rZW46CiAgICAgICAgICBhcGlTZXJ2ZXJFbmRwb2ludDogMTkyLjE2OC4xMC44Njo2NDQzCiAgICAgICAgICBjYUNlcnRIYXNoZXM6CiAgICAgICAgICAtIHNoYTI1NjpjOWQ2ZjFhYzY4YWIzNDBiYTYwNDlhNmRkMzkyMTk5NmMwMTExOWU4YzA2MTg4YWY1NTFhOTllNjhiNTI2OGNmCiAgICAgICAgICB0b2tlbjogenh4aXRuLnN3c3RkZXE2ejN2YnFkc2cKICAgICAga2luZDogSm9pbkNvbmZpZ3VyYXRpb24KICAgICAgbm9kZVJlZ2lzdHJhdGlvbjoKICAgICAgICBrdWJlbGV0RXh0cmFBcmdzOgogICAgICAgICAgYW5vbnltb3VzLWF1dGg6ICJmYWxzZSIKICAgICAgICAgIHByb3ZpZGVyLWlkOiBQUk9WSURFUl9JRAogICAgICAgICAgcmVhZC1vbmx5LXBvcnQ6ICIwIgogICAgICAgICAgdGxzLWNpcGhlci1zdWl0ZXM6IFRMU19FQ0RIRV9SU0FfV0lUSF9BRVNfMTI4X0dDTV9TSEEyNTYKICAgICAgICB0YWludHM6IG51bGwKICAgICAgcGF1c2U6CiAgICAgICAgaW1hZ2VSZXBvc2l0b3J5OiBwdWJsaWMuZWNyLmF3cy9la3MtZGlzdHJvL2t1YmVybmV0ZXMvcGF1c2UKICAgICAgICBpbWFnZVRhZzogdjEuMjUuOC1la3MtMS0yNS0xMAogICAgICBwcm94eToge30KICAgICAgcmVnaXN0cnlNaXJyb3I6IHt9CiAgICAgIApydW5jbWQ6ICJXb3JrZXJKb2luIgo=\"\n\n[settings.kubernetes]\ncluster-domain = \"cluster.local\"\nstandalone-mode = true\nauthentication-mode = \"tls\"\nserver-tls-bootstrap = false\npod-infra-container-image = \"public.ecr.aws/eks-distro/kubernetes/pause:v1.25.8-eks-1-25-10\"\nprovider-id = \"tinkerbell://eksa-system/eksabmhp1-dp-n-1\"\n\n[settings.network]\nhostname = \"mgmt-md-0-5cd48579fd-cxztz\""
        }
    }
]`