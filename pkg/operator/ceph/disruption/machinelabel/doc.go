/*
Copyright 2019 The Rook Authors. All rights reserved.

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
Package machinelabel implements the controller for ensuring that machines are labeled in correct manner for fencing.
The design and purpose for machine disruption management is found at:
https://github.com/koor-tech/koor/blob/master/design/ceph/ceph-openshift-fencing-mitigation.md
*/

package machinelabel
