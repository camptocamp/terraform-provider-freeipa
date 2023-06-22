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
	"encoding/json"
	"fmt"
)

const (
	FailedReasonNoSuchEntry    = "no such entry"
	FailedReasonAlreadyAMember = "This entry is already a member"
)

type FailedOperations map[string]map[string]failedOperations

type fromRootFailedOperations map[string]failedOperations

func (f fromRootFailedOperations) String() string {
	userFriendlyFailures := make(map[string]string)
	for rootFailName, fOperations := range f {
		for _, fOperation := range fOperations {
			fromRootName := fmt.Sprintf("%s/%s", rootFailName, fOperation.Name)
			userFriendlyFailures[fromRootName] = fOperation.Reason
		}
	}

	return fmt.Sprintf("%+v", userFriendlyFailures)
}

func (f FailedOperations) GetFailures() fromRootFailedOperations {
	failures := make(fromRootFailedOperations)
	for rootFailureCategoryName, v := range f {
		for subFailureCategoryName, failedOp := range v {
			if len(failedOp) > 0 {
				fromRootFailureName := fmt.Sprintf("%s/%s", rootFailureCategoryName, subFailureCategoryName)
				failures[fromRootFailureName] = append(failures[fromRootFailureName], failedOp...)
			}
		}
	}

	return failures
}

type failedOperation struct {
	Name   string
	Reason string
}

type failedOperations []failedOperation

func (f *failedOperations) UnmarshalJSON(b []byte) error {
	var rawFailedStr [][]string
	if err := json.Unmarshal(b, &rawFailedStr); err != nil {
		return err
	}

	*f = failedOperations{}
	for _, failedEntry := range rawFailedStr {
		if len(failedEntry) != 2 {
			return fmt.Errorf("failed entry %v does not have two elements", failedEntry)
		}

		*f = append(*f, failedOperation{
			Name:   failedEntry[0],
			Reason: failedEntry[1],
		})
	}

	return nil
}
