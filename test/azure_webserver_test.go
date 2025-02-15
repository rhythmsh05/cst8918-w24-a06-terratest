package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Subscription ID for Azure
var subscriptionID string = "06926fb4-9ce6-411d-9ff8-1bdec3ad355c"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "shar0855",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}

func TestNetworkInterfaceExists(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	terraform.InitAndApply(t, terraformOptions)

	nicID := terraform.Output(t, terraformOptions, "nic_id")
	vmID := terraform.Output(t, terraformOptions, "vm_id")

	// Test that the NIC exists
	nicExists := azure.NetworkInterfaceExists(t, nicID, subscriptionID)
	assert.True(t, nicExists, "Network Interface should exist")

	// Check if NIC is connected to the VM
	vmNicConnection := azure.CheckNICConnection(t, vmID, nicID, subscriptionID)
	assert.True(t, vmNicConnection, "NIC should be connected to the VM")

	terraform.Destroy(t, terraformOptions)
}

func TestUbuntuVersion(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	terraform.InitAndApply(t, terraformOptions)

	publicIP := terraform.Output(t, terraformOptions, "public_ip")

	sshUser := "azureuser"
	sshKeyPath := "/path/to/your/ssh/private/key"
	versionCheckCmd := "lsb_release -a"
	output := terraform.RunSSHCommand(t, publicIP, sshUser, sshKeyPath, versionCheckCmd)

	assert.Contains(t, output, "Ubuntu 20.04", "VM should be running Ubuntu 20.04")

	terraform.Destroy(t, terraformOptions)
}
