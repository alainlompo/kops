/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kops/util/pkg/vfs"
)

func (c *openstackCloud) ListKeypair(name string) (*keypairs.KeyPair, error) {
	var k *keypairs.KeyPair
	done, err := vfs.RetryWithBackoff(readBackoff, func() (bool, error) {
		rs, err := keypairs.Get(c.novaClient, name).Extract()
		if err != nil {
			if err.Error() == ErrNotFound {
				return true, nil
			}
			return false, fmt.Errorf("error listing keypair: %v", err)
		}
		k = rs
		return true, nil
	})
	if err != nil {
		return k, err
	} else if done {
		return k, nil
	} else {
		return k, wait.ErrWaitTimeout
	}
}

func (c *openstackCloud) CreateKeypair(opt keypairs.CreateOptsBuilder) (*keypairs.KeyPair, error) {
	var k *keypairs.KeyPair

	done, err := vfs.RetryWithBackoff(writeBackoff, func() (bool, error) {
		v, err := keypairs.Create(c.novaClient, opt).Extract()
		if err != nil {
			return false, fmt.Errorf("error creating keypair: %v", err)
		}
		k = v
		return true, nil
	})
	if err != nil {
		return k, err
	} else if done {
		return k, nil
	} else {
		return k, wait.ErrWaitTimeout
	}
}
