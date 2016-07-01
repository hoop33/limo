//  Copyright (c) 2016 Couchbase, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cellar

import "errors"

var (
	// ErrTxClosed is returned whe operating on a closed cellar/bolt
	ErrTxClosed = errors.New("tx closed")
	// ErrTxNotWritable is returned when attempting a write operation on a
	// read only transaction
	ErrTxNotWritable = errors.New("tx not writable")
	// ErrTxIsManaged is returned when commit/rollback has been performed
	// on a managed transaction (Update/View)
	ErrTxIsManaged = errors.New("managed tx rollback/commit not allowed")
)
