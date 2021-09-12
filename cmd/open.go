/*********************************************************************
 * Copyright (c) 2020 Red Hat, Inc.
 *
 * This program and the accompanying materials are made
 * available under the terms of the Eclipse Public License 2.0
 * which is available at https://www.eclipse.org/legal/epl-2.0/
 *
 * SPDX-License-Identifier: EPL-2.0
 **********************************************************************/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"golang.org/x/net/websocket"
)

func NewOpenCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "open",
		Short:        "Open a file in Eclipse Che",
		Long:         `Open the given file in the Eclipse Che editor`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("The filename argument is required")
			}

			filename, err := filepath.Abs(args[0])

			if err != nil {
				return errors.New("File " + args[0] + " is wrong\n")
			}

			if _, err := os.Stat(filename); os.IsNotExist(err) {
				return errors.New("File " + filename + " does not exist\n")
			}

			// grab URL for theia endpoint
			workspaceId, defined := os.LookupEnv("CHE_WORKSPACE_ID")
			if !defined {
				return errors.New("CHE_WORKSPACE_ID is not defined as environment variable")
			}

			cheApi, defined := os.LookupEnv("CHE_API")
			if !defined {
				return errors.New("CHE_API is not defined as environment variable")
			}

			cheApiToken, defined := os.LookupEnv("CHE_MACHINE_TOKEN")
			if !defined {
				return errors.New("CHE_MACHINE_TOKEN is not defined as environment variable")
			}

			// ${CHE_API}/workspace/${CHE_WORKSPACE_ID}?token=${CHE_MACHINE_TOKEN}
			cheWorkpaceDetailsUrl := cheApi + "/workspace/" + workspaceId + "?token=" + cheApiToken

			resp, err := http.Get(cheWorkpaceDetailsUrl)
			if err != nil {
				return errors.New("Unable to get workspace details:" + err.Error())
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			statusCode := resp.StatusCode
			if statusCode != 200 {
				return errors.New(fmt.Sprintf("Unable to get workspace details: status %d", statusCode))
			}
			// Decoding json string into map
			m, ok := gjson.Parse(string(body)).Value().(map[string]interface{})
			if !ok {
				return errors.New("Unable to get workspace details Parsing invalid")
			}
			runtimeMap := m["runtime"].(map[string]interface{})
			machinesMap := runtimeMap["machines"].(map[string]interface{})
			var theiaUrl string
			for _, machineElement := range machinesMap {
				machineDetails := machineElement.(map[string]interface{})
				serverAttributes := machineDetails["servers"]
				if serverAttributes != nil {
					theiaElement := serverAttributes.(map[string]interface{})["theia"]
					if theiaElement != nil {
						url := theiaElement.(map[string]interface{})["url"]
						if url != nil {
							theiaUrl = url.(string)
						}
					}
				}
			}

			origin := theiaUrl
			url := "ws://127.0.0.1:3100/services"
			ws, err := websocket.Dial(url, "", origin)
			if err != nil {
				return errors.New("Eclipse Che IDE is not running:" + err.Error())
			}

			defer ws.Close()

			openMessage := map[string]interface{}{
				"path": "/services/cli-endpoint",
				"kind": "open",
				"id":   0,
			}
			openMessageJson, err := json.Marshal(openMessage)
			if err != nil {
				return errors.New("Unable to marshal JSON:" + err.Error())
			}
			ws.Write(openMessageJson)

			contentJson := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      0,
				"method":  "openFile",
				"params":  filename,
			}
			contentEncoded, err := json.Marshal(contentJson)
			if err != nil {
				return errors.New("Unable to marshal JSON:" + err.Error())
			}
			content := string(contentEncoded)

			dataMessage := map[string]interface{}{
				"path":    "/services/cli-endpoint",
				"kind":    "data",
				"id":      0,
				"content": content,
			}

			dataMessageJson, err := json.Marshal(dataMessage)
			if err != nil {
				return errors.New("Unable to marshal JSON:" + err.Error())
			}
			_, err = ws.Write(dataMessageJson)
			if err != nil {
				return err
			}

			return nil

		},
	}
}

func init() {
	openCmd := NewOpenCmd()
	rootCmd.AddCommand(openCmd)
}
