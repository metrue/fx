package middlewares

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	"github.com/metrue/go-ssh-client"
	sshMocks "github.com/metrue/go-ssh-client/mocks"
)

func TestDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := mockCtx.NewMockContexter(ctrl)

	kubeconf, err := ioutil.TempFile("", "*.kubeconf")
	if err != nil {
		t.Fatal(err)
	}
	config := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://kubernetes.docker.internal:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-desktop
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-for-desktop
current-context: docker-desktop
kind: Config
preferences: {}
users:
- name: docker-desktop
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM5RENDQWR5Z0F3SUJBZ0lJZmd6Rml2L0lKVzR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RFeU1UWXdPREU0TWpkYUZ3MHlNVEF6TURrd016RTNNamRhTURZeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sc3dHUVlEVlFRREV4SmtiMk5yWlhJdFptOXlMV1JsCmMydDBiM0F3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzdtQk9tdHArMWU3TEsKYzNnWDVLcHU2aWRvSTd2V3lmeXRpdGxXU3d3cHc0d1pmN0d6QnhncTdXS1l5WERBSjVNR1JIeElkekVQMHcyUQo2WjlpV2d4Yk9YS2NMRCtvc1dPSlR2azB2NzBGUmg5QUNTWmNYTTQrakxsUG1VNXNZY0xBNi9RK24yQitxc3o3Ckhwb0FzZlpieVFmV0MvTG9uUmY5QVVONmVLcjFuZzVSeGdqMDFQNnN2bHBjUy9BTjhLNjZTcFJkczVocGVUU2IKMUYwUXlERFhXT0w0QTNlZUhPZGIzaC9tT3dtOEdkb1dSbGhJQjV5enlKVzF0Ny9pNVVmZm9BLzFZT0pOc2pEUgo2ZElYUTF4djJWQ0IzZnNZc0Z3NWFqb2s0aDNYSlB4N04yUWxxeWlPeXBXVjYzdDc1QmtHR083SHUxWlhUa2RRCndjc0dnN3ZUQWdNQkFBR2pKekFsTUE0R0ExVWREd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUYKQlFjREFqQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFZbEE0OXdwY0p5eUExbmlWK25UMlI2bTBURXNjMkdYTgpad0RXZFVrRWgwRHN2dGp2NUhqU1BCMmZiOUZ2V3VpcTI3YU1aTmVDdEt4WGlNcnpTUmh2YmNqS3pFQ0M2VVBDCkJQdUw3NysrZk42Rlh5eG9haXExbjRVcGhSNDl0azI4eEJORm1DTWdyKzFLOWExTDgrQ2RtaXNFRHAwKzVPQmgKYWhIUytSYk5pUmhBQ1YrazJSZ01tK0swemNvUm41c0pGS012SitCZWhuTEdUN0VLVjRFMnpOZkZiQUI0b0k3eAo0NmxsSVBKKzN1Uy9oUlEzNkR5bzZ0OUc3K204NXBKTmFFdkFJdmVxaXJ5VlRJbjF5T1BURXVFREVTM1FZaWxYCmJOZFFBNUNXcWxKRnhhZkNlNzMwNE5sbVVhTW0xZTROVHB2cmQyT2FzUDRKT1JPY2dJK1dVUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdTVnVHByYWZ0WHV5eW5ONEYrU3FidW9uYUNPNzFzbjhyWXJaVmtzTUtjT01HWCt4CnN3Y1lLdTFpbU1sd3dDZVRCa1I4U0hjeEQ5TU5rT21mWWxvTVd6bHluQ3cvcUxGamlVNzVOTCs5QlVZZlFBa20KWEZ6T1BveTVUNWxPYkdIQ3dPdjBQcDlnZnFyTSt4NmFBTEgyVzhrSDFndnk2SjBYL1FGRGVuaXE5WjRPVWNZSQo5TlQrckw1YVhFdndEZkN1dWtxVVhiT1lhWGswbTlSZEVNZ3cxMWppK0FOM25oem5XOTRmNWpzSnZCbmFGa1pZClNBZWNzOGlWdGJlLzR1VkgzNkFQOVdEaVRiSXcwZW5TRjBOY2I5bFFnZDM3R0xCY09XbzZKT0lkMXlUOGV6ZGsKSmFzb2pzcVZsZXQ3ZStRWkJoanV4N3RXVjA1SFVNSExCb083MHdJREFRQUJBb0lCQVFDZHV3REs3RUxkRldULwpWSmRsZjU3T0k1Tit2SXp6ekdIb2lSYTB0K1ZDT0dsVUIwb2lmWlNVZzRTamNyeWExS3VLV1lzbVl4R2RmSmVyCmdNUENyblExUDloZDk5YU93Smd3bTNadUk4bUs1YXJnN05DVVdIUVJvOEVzYkhyRUptN2FSNHJXSEt2RjFWY0UKem5ZdW4zUEZPUUtkdHU1SEo4OURyQXhRcmFVUlhxRE5OUUtKTGYydFNIWEJqNUtDam5jSXRBRXBwM2V3eXo3WQpnVmtlU0NrY3BRL0lZanJ0QS83c3dGR3IydWIrek9Jd1dKc3c4V01WcGZrbHZDVmR4VEk5YkdBQVRYNWdXMzZGCjhDSGRkUzd4UU5obEprbjNUVUVjeFU5UUw2TWEwNWdxWkcxVXBGaHpKQ0gybHlmQWhqQkZSbFkvcEN4YmJ1QXYKOFFGUjUyb0JBb0dCQU1HSzQ5bGFWUGhiak5Dc3J5N2dxUWVaYWZWbWQ1V2pMUzVrRFpYb0lQZkgxZWFVS253Ywp0cjd2SjJrNHZQRVR4eGdzOUtDMk5HSUFiUUIySlBFK3daYVo1OTk3NFdiUTNIaXVwWTBxMit3aFNseWpyWkxPCjI2MyswUzBLUVFDYUk1SkNseUJxTjFUOUdTeEMwYzlPTE9YMnpMbkMzUEthRHpZeHRyWlZGNDc5QW9HQkFQZ2gKdzhrSWxXYUdZd2RvMHNaeHowTzJoL2lmQ1U5ZmFGaTFpNjBrN2tBR2hxVVB4TCtyN0F5TDc4N2FlcnoxKzN5bApqVzB4YnFPZi9VTFZlaFBNWU40b0pFZkJpWmtCNnZ2Ri9HMllKeWg0MnJvQ2VDUmNLRVBqSjBST3FyQVI3Nks1CkZuMitWL2Nqc29Dd2ZsTjI4VC9LcUl4bjU3TWd1Wll1Ykt4Uk1qY1BBb0dBVWdvUnN4eDdVQnRlZ1VYeHJDbEcKL1JXbXVJTUt4Yjg1YzZTdHJaR01CL3dKUzRnYXlpbFJ2WFdhZXh1MTIyckt4aENvVVVkcXhPL3hSSFRRREFMUwpCSWlRcFViWnNMOXY5U2Z5dlBnaDZPSGpwNGtxRmtUaEVjd2wxclcyQUE5V2JMVVZZb1FqbUQ4QTRLWWlVWUdOCnZwenpBdnI2dFV0Z2oxUmJZc2FIQ2ZFQ2dZRUE3eDd5NDZoKythZW1oWHh5SzBXQVhSdnBteUlBUWRxSzMzcE4KR2RYT09DdFIxSDRHdUVRQkhmSTVieG5EVUppbysrMDdCckN0azhmWnRHK3Z6cWFWNzJHMTNPVFpLbmZic1RpUwpWRGRkL1RYQ2E2RjNrR3F6YndEWVZZNk9GVkdqb3loRlVYWitwUzlrbFhvQXM0U2JaME53L0tZaGR0R2hwK1lqCldraUJZT2NDZ1lCVmZDZGpCWDJvUmZIWkNJWm9NRVpuckNOS04rWWp1dHBnQTg5V1BEMmhPNlJLTG1DeG5GZ0oKbnFBUnUrZDBzMkhzaUErUzBKRXVMTVcxRjdnaml3Zm1zc3lSWkttWlNWaTZJZnFKUlFnSktDR3VzM1lPc2RHZQovSm5hTUlaRG1DRzNyNXpwZVl4VEowM0NIR2pReDl4dDRrallFU1F5MHdlVUtwelRITEhmWEE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
	if err := ioutil.WriteFile(kubeconf.Name(), []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(kubeconf.Name())

	sshClient := sshMocks.NewMockClienter(ctrl)
	sshClient.EXPECT().Connectable(10*time.Second).Return(true, nil).Times(3)
	sshClient.EXPECT().RunCommand("uname -a", gomock.Any()).Return(nil)
	sshClient.EXPECT().RunCommand("docker version", ssh.CommandOptions{
		Timeout: 10 * time.Second,
	}).Return(nil)
	sshClient.EXPECT().RunCommand("docker inspect fx-agent", ssh.CommandOptions{
		Timeout: 10 * time.Second,
	}).Return(nil)

	cntx := context.Background()
	ctx.EXPECT().GetContext().Return(cntx).Times(2)
	ctx.EXPECT().Get("ssh").Return(sshClient)
	ctx.EXPECT().Get("host").Return("1.2.3.4")
	ctx.EXPECT().Get("kubeconf").Return(kubeconf.Name())
	// TODO mock http call
	if err := Driver(ctx); err == nil {
		t.Fatal("should failed on initial docker client dude to /version api not ready")
	}
}
