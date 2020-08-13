// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceMicrosoftADPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrosoftADPolicyCreate,
		Read:   resourceMicrosoftADPolicyRead,
		Update: resourceMicrosoftADPolicyUpdate,
		Delete: resourceMicrosoftADPolicyDelete,
		Schema: map[string]*schema.Schema{
			"microsoft_ad_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microsoft_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"computer_name_letter_case": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Either Lowercase or Uppercase",
			},
			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func bindMicrosoftADPolicyResource(d *schema.ResourceData, policy *MicrosoftADPolicy) error {
	d.SetId(fmt.Sprintf("%v", policy.ID))

	if err := d.Set("name", policy.Name); err != nil {
		return errors.WithMessage(err, "cannot set name")
	}
	if err := d.Set("description", policy.Description); err != nil {
		return errors.WithMessage(err, "cannot set description")
	}
	if err := d.Set("workspace_url", policy.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "cannot set workspace")
	}
	microsoftEndpointURLSplit := strings.Split(policy.Links.MicrosoftEndpoint.Href, "/")
	microsoftEndpointID := microsoftEndpointURLSplit[len(microsoftEndpointURLSplit)-1]
	if err := d.Set("microsoft_endpoint_id", microsoftEndpointID); err != nil {
		return errors.WithMessage(err, "cannot set microsoft_endpoint")
	}
	if err := d.Set("computer_name_letter_case", policy.ComputerNameLetterCase); err != nil {
		return errors.WithMessage(err, "cannot set computer_name_letter_case")
	}
	if err := d.Set("ou", policy.OU); err != nil {
		return errors.WithMessage(err, "cannot set ou")
	}

	return nil
}

func resourceMicrosoftADPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("calling resourceMicrosoftADPolicyCreate")

	newPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
	}
	config := m.(Config)
	policy, err := config.NewOneFuseApiClient().CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return err
	}

	err = bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftADPolicyRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("microsoft_ad_policy_id").(int)
	policy, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(id)
	bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftADPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("microsoft_endpoint_id") ||
		d.HasChange("computer_name_letter_case") ||
		d.HasChange("workspace_url") ||
		d.HasChange("ou"))
	if !changed {
		return nil
	}

	// Create the desired AD Policy object
	id := d.Get("id").(int)
	desiredPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		OU:                     d.Get("ou").(string),
	}

	// Make the API call to update the policy
	config := m.(Config)
	updatedPolicy, err := config.NewOneFuseApiClient().UpdateMicrosoftADPolicy(id, &desiredPolicy)
	if err != nil {
		return err
	}

	// Update Terraform's state
	d.Set("name", updatedPolicy.Name)
	d.Set("description", updatedPolicy.Description)
	d.Set("microsoft_endpoint", updatedPolicy.MicrosoftEndpointID)
	d.Set("computer_name_letter_case", updatedPolicy.ComputerNameLetterCase)
	d.Set("workspace_url", updatedPolicy.WorkspaceURL)
	d.Set("ou", updatedPolicy.OU)

	return nil
}

func resourceMicrosoftADPolicyDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("microsoft_ad_policy_id").(int)
	return config.NewOneFuseApiClient().DeleteMicrosoftADPolicy(id)
}
