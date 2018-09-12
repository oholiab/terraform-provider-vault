package vault

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/vault/api"
)

func genericPostResource() *schema.Resource {
	return &schema.Resource{
		Create: genericPostWrite,
		Update: genericPostWrite,
		Delete: genericPostDelete,
		Read:   genericPostRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the genericPost",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Path for resource",
			},
			"data_json": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The document to post in JSON format",
			},
		},
	}
}

func genericPostWrite(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	name := d.Get("name").(string)
	path := "/v1/" + d.Get("path").(string)
	genericPost := d.Get("data_json").(string)

	log.Printf("[DEBUG] Writing genericPost %s to Vault", name)
	r := client.NewRequest("POST", path)
	r.Body = strings.NewReader(genericPost)

	resp, err := client.RawRequest(r)
	if err != nil {
		return fmt.Errorf("error writing to Vault: %s", err)
	}
	defer resp.Body.Close()

	d.SetId(name)

	return genericPostRead(d, meta)
}

func genericPostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	name := d.Id()
	path := "/v1/" + d.Get("path").(string)

	log.Printf("[DEBUG] Deleting genericPost %s from Vault", name)

	_, err := client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("error deleting from Vault: %s", err)
	}

	return nil
}

func genericPostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	name := d.Id()
	path := "/v1/" + d.Get("path").(string)

	genericPost, err := client.Logical().Read(path)

	if err != nil {
		return fmt.Errorf("error reading from Vault: %s", err)
	}

	d.Set("genericPost", genericPost)
	d.Set("name", name)

	return nil
}
