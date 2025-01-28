package mail

import (
	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfigFrom("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
	<h1>Wassup world</h1>
	<p>This is a test message from Tom Cruise :)</p>
	`
	to := []string{"musor.akkount7@gmail.com"}
	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
