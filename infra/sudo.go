package infra

// Sudo append sudo when user is not root
func Sudo(cmd string, user string) string {
	if user == "root" {
		return cmd
	}
	return "sudo " + cmd
}
