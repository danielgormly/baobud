package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

var serverCmd *exec.Cmd
var vaultDevClient *VaultDevClient

const testAddr = "http://localhost:8200"
const testToken = "dev"

var baobudConfig = &BaobudConfig{
	BaoAddress: testAddr,
	BaoToken:   testToken,
}

type VaultDevClient struct {
	client *api.Client
}

func NewVaultDevClient() *VaultDevClient {
	config := api.DefaultConfig()
	config.Address = testAddr
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}
	client.SetToken(testToken)
	return &VaultDevClient{client: client}
}

func (v *VaultDevClient) CreateKV2Database(name string) error {
	mountInput := &api.MountInput{
		Type:        "kv",
		Description: "baobud-test-suite",
		Options: map[string]string{
			"version": "2",
		},
	}
	err := v.client.Sys().Mount(name, mountInput)
	if err != nil {
		return fmt.Errorf("failed to create KV2 '%s' database: %w", name, err)
	}
	return nil
}

func (c *VaultDevClient) WriteKV2(path string, data map[string]interface{}) error {
	kvData := map[string]interface{}{
		"data": data,
	}
	secretPath := path
	_, err := c.client.Logical().Write(secretPath, kvData)
	return err
}

func (c *VaultDevClient) ReadKV2(path string) (map[string]interface{}, error) {
	secretPath := path
	secret, err := c.client.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}
	// KV2 secrets store the actual data in a nested "data" field
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse secret data")
	}
	return data, nil
}

func TestMain(m *testing.M) {
	// Start the server before running tests
	serverCmd = exec.Command("bao", "server", "-dev", "-dev-root-token-id=dev")

	// Optionally capture server output
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	err := serverCmd.Start()
	if err != nil {
		panic(err)
	}

	// Ensure cleanup even if tests panic
	defer func() {
		if serverCmd.Process != nil {
			fmt.Println("ok we try exit")
			// Send an interrupt signal to the process
			if err := serverCmd.Process.Kill(); err != nil {
				fmt.Printf("Failed to kill process: %v\n", err)
			}
			// Wait for the process to exit
			serverCmd.Wait()
		}
	}()

	// Give the server a moment to start up
	time.Sleep(2 * time.Second)

	vaultDevClient = NewVaultDevClient()
	vaultDevClient.CreateKV2Database("kv2")
	fmt.Println("We created a database, allegedly")

	// Run the tests
	code := m.Run()

	if serverCmd.Process != nil {
		serverCmd.Process.Kill()
	}

	os.Exit(code)
	return
}

func TestAnalyze(t *testing.T) {
	fmt.Println("TestAnalyze")
	vaultDevClient.WriteKV2("kv2/data/test", map[string]interface{}{
		"foo": "bar",
		"bar": "foo",
	})
	file := ReadFile("../test/template.ctmpl")
	evaluated, err := Analyze(string(file), *baobudConfig)
	if err != nil {
		log.Fatalf("Failed to evaluate %v", err)
	}
	if evaluated[0] != "kv2/test" {
		t.Errorf("Expected to receive dependency kv2/test, received %v instead", evaluated[0])
	}
	policy := CreateVaultPolicy(evaluated)
	expected := `path "kv2/test" {
    capabilities = ["read"]
}`
	if strings.TrimSpace(policy) != strings.TrimSpace(expected) {
		log.Fatalf("Policy not expected")
	}
}
