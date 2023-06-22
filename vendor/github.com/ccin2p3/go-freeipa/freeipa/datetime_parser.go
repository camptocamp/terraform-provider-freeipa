// Copyright © 2022 IN2P3 Computing Centre, IN2P3, CNRS
// Copyright © 2018 Philippe Voinov
//
// Contributor(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2021
//
// This software is governed by the CeCILL license under French law and
// abiding by the rules of distribution of free software.  You can  use,
// modify and/ or redistribute the software under the terms of the CeCILL
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info".
//
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability.
//
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or
// data to be ensured and,  more generally, to use and operate it in the
// same conditions as regards security.
//
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL license and that you accept its terms.

package freeipa

import (
	"fmt"
	"time"
)

const (
	// 	RFC3339     = "2006-01-02T15:04:05Z07:00"
	LDAPGeneralizedTimeFormat = "20060102150405Z"
)

// parse LDAP generalized time format as present
// in some "__datetime__" responses
// See https://github.com/freeipa/freeipa/blob/master/ipalib/constants.py#L280
// LDAP_GENERALIZED_TIME_FORMAT = "%Y%m%d%H%M%SZ"
func parseFreeIPADateTimeStr(str string) (time.Time, error) {
	return time.Parse(LDAPGeneralizedTimeFormat, str)
}

// tryParseFreeIPADatetimeMap tries to solve https://github.com/ccin2p3/go-freeipa/issues/1
// Krbprincipalexpiration is returned as a []interface {}
// that is [map[__datetime__:20220428000000Z]]
func tryParseFreeIPADatetimeMap(m map[string]interface{}) (time.Time, error) {
	var tt time.Time
	dV, ok := m["__datetime__"]
	if !ok {
		return tt, fmt.Errorf("no __datetime__ key")
	}

	dsV, ok := dV.(string)
	if !ok {
		return tt, fmt.Errorf("__datetime__ key not a string")
	}

	return parseFreeIPADateTimeStr(dsV)
}
