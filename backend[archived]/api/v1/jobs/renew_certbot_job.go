package cronjobs

import (
	"log"
	"os/exec"
)

type RenewCertbotJob struct {
}

func NewRenewCertbotJob() *RenewCertbotJob {
	return &RenewCertbotJob{}
}

func (j *RenewCertbotJob) Run() {
	cmd := exec.Command("docker-compose", "run", "certbot-renew")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error renewing certificate: %v", err)
	}
	log.Printf("Certbot output: %s", string(cmdOutput))
}

func (j *RenewCertbotJob) Schedule() string {
	return "0 0 0 * * *"
}
