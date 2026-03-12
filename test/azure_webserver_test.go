package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var subscriptionID string = "9341defd-4cd3-457c-bdf1-12be7fd2bb93"

func TestAzureLinuxVMCreation(t *testing.T) {

	os.Setenv("PATH", os.Getenv("PATH")+";C:\\Program Files\\Microsoft SDKs\\Azure\\CLI2\\wbin")

	terraformOptions := &terraform.Options{
		TerraformDir: "../",

		Vars: map[string]interface{}{
			"labelPrefix": "shap0011",
		},

		NoColor: true,

		EnvVars: map[string]string{
			"TF_CLI_ARGS_apply":   "-parallelism=1",
			"TF_CLI_ARGS_destroy": "-parallelism=1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Test 1: VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Test 2: NIC exists and is connected to VM
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))

	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)

	assert.NotNil(t, vm.NetworkProfile)
	assert.Len(t, *vm.NetworkProfile.NetworkInterfaces, 1)
	assert.Contains(t, *(*vm.NetworkProfile.NetworkInterfaces)[0].ID, nicName)

	// Test 3: VM running correct Ubuntu image
	assert.Equal(t, "Canonical", *vm.StorageProfile.ImageReference.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", *vm.StorageProfile.ImageReference.Offer)
	assert.Equal(t, "22_04-lts-gen2", *vm.StorageProfile.ImageReference.Sku)
}
