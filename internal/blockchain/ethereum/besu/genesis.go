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

package besu

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Storage struct {
	Field1 string `json:"0x0000000000000000000000000000000000000000000000000000000000000000"`
	Field2 string `json:"0x0000000000000000000000000000000000000000000000000000000000000001"`
	Field3 string `json:"0x0000000000000000000000000000000000000000000000000000000000000004"`
}

type Genesis struct {
	Config     *GenesisConfig    `json:"config"`
	Nonce      string            `json:"nonce"`
	Timestamp  string            `json:"timestamp"`
	ExtraData  string            `json:"extraData"`
	GasLimit   string            `json:"gasLimit"`
	Difficulty string            `json:"difficulty"`
	MixHash    string            `json:"mixHash"`
	Coinbase   string            `json:"coinbase"`
	Alloc      map[string]*Alloc `json:"alloc"`
	Number     string            `json:"number"`
	GasUsed    string            `json:"gasUsed"`
	ParentHash string            `json:"parentHash"`
}

type GenesisConfig struct {
	ChainId                int           `json:"chainId"`
	ConstantinopleFixBlock int           `json:"constantinoplefixblock"`
	Clique                 *CliqueConfig `json:"clique"`
}

type CliqueConfig struct {
	EpochLength        int `json:"epochlength"`
	Blockperiodseconds int `json:"blockperiodseconds"`
}

type Alloc struct {
	Balance string   `json:"balance"`
	Code    string   `json:"code,omitempty"`
	Storage *Storage `json:"storage,omitempty"`
}

func (g *Genesis) WriteGenesisJson(filename string) error {
	genesisJsonBytes, _ := json.MarshalIndent(g, "", " ")
	if err := ioutil.WriteFile(filepath.Join(filename), genesisJsonBytes, 0755); err != nil {
		return err
	}
	return nil
}

func CreateGenesis(addresses []string) *Genesis {
	alloc := make(map[string]*Alloc)
	for _, address := range addresses {
		alloc[address] = &Alloc{
			Balance: "0x200000000000000000000000000000000000000000000000000000000000000",
		}

	}
	return &Genesis{
		Config: &GenesisConfig{
			ChainId:                1337,
			ConstantinopleFixBlock: 0,
			Clique: &CliqueConfig{
				Blockperiodseconds: 5,
				EpochLength:        30000,
			},
		},
		Coinbase:   "0x0000000000000000000000000000000000000000",
		Difficulty: "0x1",
		ExtraData:  "0x00000000000000000000000000000000000000000000000000000000000000004592c8e45706cc08b8f44b11e43cba0cfc5892cb0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		GasLimit:   "0xffffffff",
		MixHash:    "0x0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:      "0x0",
		Timestamp:  "0x5c51a607",
		Alloc:      alloc,
		Number:     "0x0",
		GasUsed:    "0x0",
		ParentHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
}
