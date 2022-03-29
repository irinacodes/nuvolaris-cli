// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
)

func embedSecrets(cmd *SecretsCmd) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	originalFilePath := path.Join(workingDir, cmd.OriginalFileName)
	buf := bytes.NewBuffer(nil)
	in, err := os.Open(originalFilePath)
	if err != nil {
		return err
	}
	io.Copy(buf, in)
	in.Close()

	originalContent := string(buf.Bytes())
	wskPropsMap, err := readWskPropsAsMap()
	if err != nil {
		return err
	}
	replacedContent := originalContent
	for k, v := range wskPropsMap {
		replacedContent = strings.ReplaceAll(replacedContent, "@"+k+"@", v)
	}
	replacedFilePath := path.Join(workingDir, cmd.ProcessedFileName)
	out, err := os.Create(replacedFilePath)
	if err != nil {
		return err
	}
	_, err = out.WriteString(replacedContent)
	if err != nil {
		return err
	}
	return out.Close()
}
