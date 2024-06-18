package freeipa

func utilsGetArry(itemsRaw []interface{}) []string {
	res := make([]string, len(itemsRaw))
	for i, raw := range itemsRaw {
		res[i] = raw.(string)
	}
	return res
}
