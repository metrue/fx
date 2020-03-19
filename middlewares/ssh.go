package middlewares

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/go-ssh-client"
)

// SSH create a ssh client
func SSH(ctx context.Contexter) error {
	host := ctx.Get("host").(string)
	user := ctx.Get("user").(string)
	port := ctx.Get("ssh_port").(string)
	keyfile := ctx.Get("ssh_key").(string)
	sshClient := ssh.New(host).WithUser(user).WithPort(port).WithKey(keyfile)
	ctx.Set("ssh", sshClient)
	return nil
}
