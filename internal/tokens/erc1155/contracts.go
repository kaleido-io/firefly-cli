// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package erc1155

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hyperledger/firefly-cli/internal/blockchain/ethereum"
	"github.com/hyperledger/firefly-cli/internal/blockchain/ethereum/ethconnect"
	"github.com/hyperledger/firefly-cli/internal/constants"
	"github.com/hyperledger/firefly-cli/internal/log"
	"github.com/hyperledger/firefly-cli/pkg/types"
)

const TOKEN_URI_PATTERN = "firefly://token/{id}"

func DeployContracts(s *types.Stack, log log.Logger, verbose bool, tokenIndex int) error {
	var containerName string
	for _, member := range s.Members {
		if !member.External {
			containerName = fmt.Sprintf("%s_tokens_%s_%d", s.Name, member.ID, tokenIndex)
			break
		}
	}
	if containerName == "" {
		return errors.New("unable to extract contracts from container - no valid tokens containers found in stack")
	}
	log.Info("extracting smart contracts")

	if err := ethereum.ExtractContracts(s.Name, containerName, "/root/contracts", verbose); err != nil {
		return err
	}

	tokenContract, err := ethereum.ReadCompiledContract(filepath.Join(constants.StacksDir, s.Name, "contracts", "ERC1155MixedFungible.json"))
	if err != nil {
		return err
	}

	var tokenContractAddress string
	for _, member := range s.Members {
		if tokenContractAddress == "" {
			// TODO: version the registered name
			log.Info(fmt.Sprintf("deploying ERC1155 contract on '%s'", member.ID))
			tokenContractAddress, err = ethconnect.DeployContract(member, tokenContract, "erc1155", map[string]string{"uri": TOKEN_URI_PATTERN})
			if err != nil {
				return err
			}
		} else {
			log.Info(fmt.Sprintf("registering ERC1155 contract on '%s'", member.ID))
			err = ethconnect.RegisterContract(member, tokenContract, tokenContractAddress, "erc1155", map[string]string{"uri": TOKEN_URI_PATTERN})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
