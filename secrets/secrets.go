package secrets

import (
	"git-platform.dbc.dk/platform/morph/nix"
	"git-platform.dbc.dk/platform/morph/ssh"
	"git-platform.dbc.dk/platform/morph/utils"
	"os"
)

func GetSecretSize(secret nix.Secret, deploymentWD string) (size int64, err error) {
	fh, err := os.Open(utils.GetAbsPathRelativeTo(secret.Source, deploymentWD))
	if err != nil {
		return size, err
	}

	fStats, err := fh.Stat()
	if err != nil {
		return size, err
	}

	return fStats.Size(), nil
}

func UploadSecret(host nix.Host, sudoPasswd string, secret nix.Secret, deploymentWD string) (err error) {
	tempPath, err := ssh.MakeTempFile(host)
	if err != nil {
		return err
	}

	err = ssh.UploadFile(host, utils.GetAbsPathRelativeTo(secret.Source, deploymentWD), tempPath)
	if err != nil {
		return err
	}

	err = ssh.MoveFile(host, sudoPasswd, tempPath, secret.Destination)
	if err != nil {
		return err
	}

	err = ssh.SetOwner(host, sudoPasswd, secret.Destination, secret.Owner.User, secret.Owner.Group)
	if err != nil {
		return err
	}

	err = ssh.SetPermissions(host, sudoPasswd, secret.Destination, secret.Permissions)
	if err != nil {
		return nil
	}

	return nil
}