package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingFolderSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingFolderSinkCreate,
		Read:   resourceLoggingFolderSinkRead,
		Delete: resourceLoggingFolderSinkDelete,
		Update: resourceLoggingFolderSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("folder"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["folder"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The folder to be exported to the sink. Note that either [FOLDER_ID] or "folders/[FOLDER_ID]" is accepted.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "folders/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: `Whether or not to include children folders in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided folder are included.`,
	}
	schm.Schema["intercept_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: `Whether or not to intercept logs from child projects. If true, matching logs will not match with sinks in child resources, except _Required sinks. This sink will be visible to child resources when listing sinks.`,
	}
	schm.Schema["custom_writer_identity"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: `A service account provided by the caller that will be used to write the log entries. The format must be serviceAccount:some@email. This field can only be specified if you are routing logs to a destination outside this sink's folder. If not specified, a Logging service account will automatically be generated.`,
	}

	return schm
}

func resourceLoggingFolderSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	folder := resourcemanager.ParseFolderId(d.Get("folder"))
	id, sink := expandResourceLoggingSink(d, "folders", folder)
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.InterceptChildren = d.Get("intercept_children").(bool)
	customWriterIdentity := d.Get("custom_writer_identity").(string)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	folderSinkCreateRequest := NewClient(config, userAgent).Folders.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true)
	if customWriterIdentity != "" {
		folderSinkCreateRequest = folderSinkCreateRequest.CustomWriterIdentity(customWriterIdentity)
	}

	_, err = folderSinkCreateRequest.Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, err := NewClient(config, userAgent).Folders.Sinks.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Folder Logging Sink %s", d.Get("name").(string)))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	if err := d.Set("include_children", sink.IncludeChildren); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	if err := d.Set("intercept_children", sink.InterceptChildren); err != nil {
		return fmt.Errorf("Error setting intercept_children: %s", err)
	}

	return nil
}

func resourceLoggingFolderSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	customWriterIdentity := d.Get("custom_writer_identity").(string)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	folderSinkUpdateRequest := NewClient(config, userAgent).Folders.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true)
	if customWriterIdentity != "" {
		folderSinkUpdateRequest = folderSinkUpdateRequest.CustomWriterIdentity(customWriterIdentity)
	}

	_, err = folderSinkUpdateRequest.Do()
	if err != nil {
		return err
	}

	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = NewClient(config, userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_logging_folder_sink",
		ProductName: "logging",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceLoggingFolderSink(),
	}.Register()
}
