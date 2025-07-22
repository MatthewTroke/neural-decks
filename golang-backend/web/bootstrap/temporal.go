package bootstrap

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"

	"go.temporal.io/sdk/client"
)

func DialTemporal(env *Env) client.Client {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Current working directory:", wd)
	// Assuming the working directory is the project root
	clientKeyPath := filepath.Join(wd, "certs/temporal", "cardbeyondhumanity.o83iu.key")
	clientCertPath := filepath.Join(wd, "certs/temporal", "cardbeyondhumanity.o83iu.crt")

	// Specify the host and port of your Temporal Cloud Namespace
	// Host and port format: namespace.unique_id.tmprl.cloud:port
	hostPort := env.TemporalHostPort
	namespace := env.TemporalNamespace
	// Use the crypto/tls package to create a cert object
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		err = fmt.Errorf("unable to load cert and key pair: %w", err)
		panic(err)
	}
	// Add the cert to the tls certificates in the ConnectionOptions of the Client
	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	})

	if err != nil {
		err = fmt.Errorf("unable to connect to Temporal Cloud: %w", err)

		panic(err)
	}

	return c
}
