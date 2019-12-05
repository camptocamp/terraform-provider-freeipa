package freeipa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

const (
	userSchemaUid       = "uid"
	userSchemaFirstName = "first_name"
	userSchemaLastName  = "last_name"
	userSchemaUidNumber = "uidnumber"
	userSchemaGidNumber = "gidnumber"
	userSchemaGroups    = "groups"
	userSchemaEmail     = "email"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			userSchemaUid: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			userSchemaFirstName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			userSchemaLastName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			userSchemaUidNumber: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaGidNumber: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaEmail: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaGroups: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceUserCreate - %s", d.Id())

	c := m.(*Connection)

	uid := d.Get(userSchemaUid).(string)
	first := d.Get(userSchemaFirstName).(string)
	last := d.Get(userSchemaLastName).(string)

	options := map[string]interface{}{}

	uidNumber, ok := d.GetOk(userSchemaUidNumber)
	if ok {
		options["uidnumber"] = uidNumber.(string)
	}

	gidNumber, ok := d.GetOk(userSchemaGidNumber)
	if ok {
		options["gidnumber"] = gidNumber.(string)
	}

	email, ok := d.GetOk(userSchemaEmail)
	if ok {
		options["mail"] = email.(string)
	}

	log.Printf("[TRACE] creating user with name - %s, first - %s, last - %s",
		uid, first, last)

	userRec, err := c.client.CreateUser(uid, first, last, options)

	if err != nil {
		log.Printf("[ERROR] Error creating user %s - %s", uid, err)
		return err
	}

	d.SetId(string(userRec.IpaUniqueId))

	groupsInterface, ok := d.GetOk(userSchemaGroups)
	if ok {
		groupsRaw := groupsInterface.(*schema.Set)
		if groupsRaw.Len() > 0 {
			groups := make([]string, groupsRaw.Len())
			for i, d := range groupsRaw.List() {
				groups[i] = d.(string)
			}

			err = c.client.UserSyncGroups(uid, groups)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceUserRead - %s", d.Id())

	c := m.(*Connection)

	uid, err := c.ldapClient.GetUserForUUID(d.Id())

	if err != nil {
		log.Printf("[ERROR] Error reading user %s - %s", uid, err)
		return err
	}

	rec, err := c.client.GetUser(*uid)

	if err != nil {
		log.Printf("[ERROR] Error getting user %s - %s", uid, err)
		return err
	}

	err = d.Set(userSchemaUid, rec.Uid)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaUidNumber, rec.UidNumber)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaLastName, rec.Last)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaFirstName, rec.First)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaEmail, rec.Email)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaGidNumber, rec.GidNumber)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaGroups, rec.Groups)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceUserUpdate - %s", d.Id())

	c := m.(*Connection)

	uid, err := c.ldapClient.GetUserForUUID(d.Id())

	if err != nil {
		log.Printf("[ERROR] Error reading user %s - %s", uid, err)
		return err
	}

	d.Partial(true)

	if d.HasChange(userSchemaUid) {
		_, newValue := d.GetChange(userSchemaUid)
		val := newValue.(string)
		c.client.UserUpdateUid(*uid, newValue.(string))
		d.SetPartial(userSchemaUid)
		uid = &val
	}

	if d.HasChange(userSchemaEmail) {
		_, newValue := d.GetChange(userSchemaEmail)
		c.client.UserUpdateEmail(*uid, newValue.(string))
		d.SetPartial(userSchemaEmail)
	}

	if d.HasChange(userSchemaGidNumber) {
		_, newValue := d.GetChange(userSchemaGidNumber)
		c.client.UserUpdateGidNumber(*uid, newValue.(string))
		d.SetPartial(userSchemaGidNumber)
	}

	if d.HasChange(userSchemaUidNumber) {
		_, newValue := d.GetChange(userSchemaUidNumber)
		c.client.UserUpdateUidNumber(*uid, newValue.(string))
		d.SetPartial(userSchemaUidNumber)
	}

	if d.HasChange(userSchemaFirstName) {
		_, newValue := d.GetChange(userSchemaFirstName)
		c.client.UserUpdateFirstName(*uid, newValue.(string))
		d.SetPartial(userSchemaFirstName)
	}

	if d.HasChange(userSchemaLastName) {
		_, newValue := d.GetChange(userSchemaLastName)
		c.client.UserUpdateLastName(*uid, newValue.(string))
		d.SetPartial(userSchemaLastName)
	}

	if d.HasChange(userSchemaGroups) {
		_, newValueInterface := d.GetChange(userSchemaGroups)

		groupsRaw := newValueInterface.(*schema.Set)
		newValue := make([]string, groupsRaw.Len())
		for i, d := range groupsRaw.List() {
			newValue[i] = d.(string)
		}

		c.client.UserSyncGroups(*uid, newValue)
		d.SetPartial(userSchemaGroups)
	}

	d.Partial(false)

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceUserExists - %s", d.Id())

	c := m.(*Connection)

	uid, err := c.ldapClient.GetUserForUUID(d.Id())

	if err != nil {
		log.Printf("[ERROR] Error reading user %s - %s", uid, err)
		return err
	}

	return c.client.DeleteUser(*uid)
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Printf("[TRACE] resourceUserExists - %s", d.Id())

	c := m.(*Connection)

	id := d.Id()

	exists, err := c.ldapClient.UserExistsForUUID(id)

	if err != nil {
		return false, err
	}

	return exists, nil
}
