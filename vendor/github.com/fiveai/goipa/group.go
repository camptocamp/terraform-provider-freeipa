package ipa

import (
	"encoding/json"
	"errors"
)

type GroupRecord struct {
	Dn           string    `json:"dn"`
	Description  IpaString `json:"description"`
	Gid          IpaString `json:"cn"`
	GidNumber    IpaString `json:"gidnumber"`
	MepManagedBy IpaString `json:"mepmanagedby"`
	IpaUniqueId  IpaString `json:"ipauniqueid"`
	Users        []string  `json:"member_user"`
	HbacRules    []string  `json:"memberof_hbacrule"`
}

// Fetch user details by calling the FreeIPA group-show method
func (c *Client) GetGroup(gid string) (*GroupRecord, error) {
	options := map[string]interface{}{
		"no_members": false,
		"all":        true}

	res, err := c.rpc("group_show", []string{gid}, options)

	if err != nil {
		return nil, err
	}

	var groupRec GroupRecord
	err = json.Unmarshal(res.Result.Data, &groupRec)
	if err != nil {
		return nil, err
	}
	return &groupRec, nil
}

// This doesn't work for primary groups - this appears to be a deficiency
// in FreeIPA as it also doesn't work from the ipa CLI
func (c *Client) GetGroupByGidNumber(gidNumber string) (*GroupRecord, error) {

	options := map[string]interface{}{
		"no_members": false,
		"all":        true,
		"gidnumber":  gidNumber,
		"version":    "2.228"}

	res, err := c.rpc("group_find", []string{}, options)

	if err != nil {
		return nil, err
	}

	var groupRec []GroupRecord
	err = json.Unmarshal(res.Result.Data, &groupRec)
	if err != nil {
		return nil, err
	}

	if len(groupRec) != 1 {
		return nil, errors.New("too many records returned")
	}

	return &groupRec[0], nil
}

func (c *Client) GroupExists(uid string) (bool, error) {
	_, err := c.GetGroup(uid)

	if err == nil {
		return true, nil
	}

	if err.(*IpaError).Code == 4001 {
		return false, nil
	}

	return false, err
}

func (c *Client) GroupAddMember(gid string, memberId string, memberType string) (error) {
	options := map[string]interface{}{
		"raw":      true,
		"all":      true,
		memberType: memberId}

	_, err := c.rpc("group_add_member", []string{gid}, options)

	return err
}

func (c *Client) GroupRemoveMembers(gid string, members []string, memberType string) (error) {
	options := map[string]interface{}{
		"all":        false,
		"no_members": false,
		"raw":        false,
		memberType:   members,
		"version":    "2.164"}

	_, err := c.rpc("group_remove_member", []string{gid}, options)

	return err
}

func (c *Client) GroupMod(gid string, key string, value string) error {
	options := map[string]interface{}{
		key:       value,
		"version": "2.228"}

	_, err := c.rpc("group_mod", []string{gid}, options)

	return err
}

func (c *Client) GroupUpdateGid(oldGid string, newGid string) error {
	return c.GroupMod(oldGid, "rename", newGid)
}

func (c *Client) GroupUpdateGidNumber(gid string, gidNumber string) error {
	return c.GroupMod(gid, "gidnumber", gidNumber)
}

func (c *Client) GroupUpdateDescription(gid string, description string) error {
	return c.GroupMod(gid, "description", description)
}

func (c *Client) GroupRemoveMember(gid string, member string, memberType string) (error) {
	return c.GroupRemoveMembers(gid, []string{member}, memberType)
}

func (c *Client) GroupRemoveUser(gid string, uid string) (error) {
	return c.GroupRemoveMember(gid, uid, "user")
}

func (c *Client) GroupRemoveUsers(gid string, uids []string) (error) {
	return c.GroupRemoveMembers(gid, uids, "user")
}

func (c *Client) GroupAddUser(gid string, uid string) (error) {
	return c.GroupAddMember(gid, uid, "user")
}

func (c *Client) CreateGroup(gid string, description string, options map[string]interface{}) (*GroupRecord, error) {
	options["all"] = true
	options["description"] = description
	options["version"] = "2.228"
	options["no_members"] = false

	res, err := c.rpc("group_add", []string{gid}, options)

	if err != nil {
		return nil, err
	}

	var groupRec GroupRecord
	err = json.Unmarshal(res.Result.Data, &groupRec)
	if err != nil {
		return nil, err
	}
	return &groupRec, nil
}

func (c *Client) DeleteGroup(gid string) (error) {

	options := map[string]interface{}{
		"version": "2.164"}

	_, err := c.rpc("group_del", []string{gid}, options)

	return err
}