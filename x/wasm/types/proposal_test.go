package types

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestValidateProposalCommons(t *testing.T) {
	type commonProposal struct {
		Title, Description string
	}

	specs := map[string]struct {
		src    commonProposal
		expErr bool
	}{
		"all good": {src: commonProposal{
			Title:       "Foo",
			Description: "Bar",
		}},
		"prevent empty title": {
			src: commonProposal{
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent white space only title": {
			src: commonProposal{
				Title:       " ",
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent leading white spaces in title": {
			src: commonProposal{
				Title:       " Foo",
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent title exceeds max length ": {
			src: commonProposal{
				Title:       strings.Repeat("a", govtypes.MaxTitleLength+1),
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent empty description": {
			src: commonProposal{
				Title: "Foo",
			},
			expErr: true,
		},
		"prevent leading white spaces in description": {
			src: commonProposal{
				Title:       "Foo",
				Description: " Bar",
			},
			expErr: true,
		},
		"prevent white space only description": {
			src: commonProposal{
				Title:       "Foo",
				Description: " ",
			},
			expErr: true,
		},
		"prevent descr exceeds max length ": {
			src: commonProposal{
				Title:       "Foo",
				Description: strings.Repeat("a", govtypes.MaxDescriptionLength+1),
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := validateProposalCommons(spec.src.Title, spec.src.Description)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateStoreCodeProposal(t *testing.T) {
	var (
		anyAddress     = sdk.BytesToAccAddress(bytes.Repeat([]byte{0x0}, sdk.BytesAddrLen))
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    *StoreCodeProposal
		expErr bool
	}{
		"all good": {
			src: StoreCodeProposalFixture(),
		},
		"with instantiate permission": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				accessConfig := AccessTypeOnlyAddress.With(anyAddress)
				p.InstantiatePermission = &accessConfig
			}),
		},

		"without source": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.Source = ""
			}),
		},
		"base data missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"run_as missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.RunAs = ""
			}),
			expErr: true,
		},
		"run_as invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
		"wasm code missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = nil
			}),
			expErr: true,
		},
		"wasm code invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = bytes.Repeat([]byte{0x0}, MaxWasmSize+1)
			}),
			expErr: true,
		},
		"source invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.Source = "not an url"
			}),
			expErr: true,
		},
		"builder invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.Builder = "not a builder"
			}),
			expErr: true,
		},
		"with invalid instantiate permission": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.InstantiatePermission = &AccessConfig{}
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateInstantiateContractProposal(t *testing.T) {
	var (
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    *InstantiateContractProposal
		expErr bool
	}{
		"all good": {
			src: InstantiateContractProposalFixture(),
		},
		"without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Admin = ""
			}),
		},
		"without init msg": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.InitMsg = nil
			}),
			expErr: true,
		},
		"with invalid init msg": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.InitMsg = []byte("not a json string")
			}),
			expErr: true,
		},
		"without init funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = nil
			}),
		},
		"base data missing": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"run_as missing": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.RunAs = ""
			}),
			expErr: true,
		},
		"run_as invalid": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
		"admin invalid": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Admin = invalidAddress
			}),
			expErr: true,
		},
		"code id empty": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.CodeID = 0
			}),
			expErr: true,
		},
		"label empty": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Label = ""
			}),
			expErr: true,
		},
		"init funds negative": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.Coins{{Denom: "foo", Amount: sdk.NewInt(-1)}}
			}),
			expErr: true,
		},
		"init funds with duplicates": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.Coins{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "foo", Amount: sdk.NewInt(2)}}
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateMigrateContractProposal(t *testing.T) {
	var (
		invalidAddress = "invalid address2"
	)

	specs := map[string]struct {
		src    *MigrateContractProposal
		expErr bool
	}{
		"all good": {
			src: MigrateContractProposalFixture(),
		},
		"without migrate msg": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.MigrateMsg = nil
			}),
			expErr: true,
		},
		"migrate msg with invalid json": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.MigrateMsg = []byte("not a json message")
			}),
			expErr: true,
		},
		"base data missing": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"code id empty": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.CodeID = 0
			}),
			expErr: true,
		},
		"run_as missing": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.RunAs = ""
			}),
			expErr: true,
		},
		"run_as invalid": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateAdminProposal(t *testing.T) {
	var (
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    *UpdateAdminProposal
		expErr bool
	}{
		"all good": {
			src: UpdateAdminProposalFixture(),
		},
		"base data missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"admin missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.NewAdmin = ""
			}),
			expErr: true,
		},
		"admin invalid": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.NewAdmin = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateClearAdminProposal(t *testing.T) {
	var (
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    *ClearAdminProposal
		expErr bool
	}{
		"all good": {
			src: ClearAdminProposalFixture(),
		},
		"base data missing": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateContractStatusProposal(t *testing.T) {
	var (
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    *UpdateContractStatusProposal
		expErr bool
	}{
		"all good": {
			src: UpdateContractStatusProposalFixture(),
		},
		"base data missing": {
			src: UpdateContractStatusProposalFixture(func(p *UpdateContractStatusProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: UpdateContractStatusProposalFixture(func(p *UpdateContractStatusProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: UpdateContractStatusProposalFixture(func(p *UpdateContractStatusProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"status missing": {
			src: UpdateContractStatusProposalFixture(func(p *UpdateContractStatusProposal) {
				p.Status = ContractStatusUnspecified
			}),
			expErr: true,
		},
		"status invalid": {
			src: UpdateContractStatusProposalFixture(func(p *UpdateContractStatusProposal) {
				p.Status = 3
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProposalStrings(t *testing.T) {
	specs := map[string]struct {
		src govtypes.Content
		exp string
	}{
		"store code": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = []byte{01, 02, 03, 04, 05, 06, 07, 0x08, 0x09, 0x0a}
			}),
			exp: `Store Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  WasmCode:    0102030405060708090A
  Source:      https://example.com/code
  Builder:     foo/bar:latest
`,
		},
		"instantiate contract": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.Coins{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "bar", Amount: sdk.NewInt(2)}}
			}),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  Admin:       link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  Code id:     1
  Label:       testing
  InitMsg:     "{\"verifier\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\",\"beneficiary\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\"}"
  Funds:       1foo,2bar
`,
		},
		"instantiate contract without funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Funds = nil }),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  Admin:       link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  Code id:     1
  Label:       testing
  InitMsg:     "{\"verifier\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\",\"beneficiary\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\"}"
  Funds:       
`,
		},
		"instantiate contract without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Admin = "" }),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  Admin:       
  Code id:     1
  Label:       testing
  InitMsg:     "{\"verifier\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\",\"beneficiary\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\"}"
  Funds:       
`,
		},
		"migrate contract": {
			src: MigrateContractProposalFixture(),
			exp: `Migrate Contract Proposal:
  Title:       Foo
  Description: Bar
  Contract:    link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
  Code id:     1
  Run as:      link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
  MigrateMsg   "{\"verifier\":\"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5\"}"
`,
		},
		"update admin": {
			src: UpdateAdminProposalFixture(),
			exp: `Update Contract Admin Proposal:
  Title:       Foo
  Description: Bar
  Contract:    link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
  New Admin:   link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
`,
		},
		"clear admin": {
			src: ClearAdminProposalFixture(),
			exp: `Clear Contract Admin Proposal:
  Title:       Foo
  Description: Bar
  Contract:    link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
`,
		},
		"pin codes": {
			src: &PinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{1, 2, 3},
			},
			exp: `Pin Wasm Codes Proposal:
  Title:       Foo
  Description: Bar
  Codes:       [1 2 3]
`,
		},
		"unpin codes": {
			src: &UnpinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{3, 2, 1},
			},
			exp: `Unpin Wasm Codes Proposal:
  Title:       Foo
  Description: Bar
  Codes:       [3 2 1]
`,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			assert.Equal(t, spec.exp, spec.src.String())
		})
	}
}

func TestProposalYaml(t *testing.T) {
	specs := map[string]struct {
		src govtypes.Content
		exp string
	}{
		"store code": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = []byte{01, 02, 03, 04, 05, 06, 07, 0x08, 0x09, 0x0a}
			}),
			exp: `title: Foo
description: Bar
run_as: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
wasm_byte_code: AQIDBAUGBwgJCg==
source: https://example.com/code
builder: foo/bar:latest
instantiate_permission: null
`,
		},
		"instantiate contract": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.Coins{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "bar", Amount: sdk.NewInt(2)}}
			}),
			exp: `title: Foo
description: Bar
run_as: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
admin: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
code_id: 1
label: testing
init_msg: '{"verifier":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5","beneficiary":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5"}'
funds:
- denom: foo
  amount: "1"
- denom: bar
  amount: "2"
`,
		},
		"instantiate contract without funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Funds = nil }),
			exp: `title: Foo
description: Bar
run_as: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
admin: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
code_id: 1
label: testing
init_msg: '{"verifier":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5","beneficiary":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5"}'
funds: []
`,
		},
		"instantiate contract without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Admin = "" }),
			exp: `title: Foo
description: Bar
run_as: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
admin: ""
code_id: 1
label: testing
init_msg: '{"verifier":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5","beneficiary":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5"}'
funds: []
`,
		},
		"migrate contract": {
			src: MigrateContractProposalFixture(),
			exp: `title: Foo
description: Bar
contract: link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
code_id: 1
msg: '{"verifier":"link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5"}'
run_as: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
`,
		},
		"update admin": {
			src: UpdateAdminProposalFixture(),
			exp: `title: Foo
description: Bar
new_admin: link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5
contract: link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
`,
		},
		"clear admin": {
			src: ClearAdminProposalFixture(),
			exp: `title: Foo
description: Bar
contract: link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8
`,
		},
		"pin codes": {
			src: &PinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{1, 2, 3},
			},
			exp: `title: Foo
description: Bar
code_ids:
- 1
- 2
- 3
`,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			v, err := yaml.Marshal(&spec.src)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, string(v))
		})
	}
}

func TestConvertToProposals(t *testing.T) {
	cases := map[string]struct {
		input     string
		isError   bool
		proposals []ProposalType
	}{
		"one proper item": {
			input:     "UpdateAdmin",
			proposals: []ProposalType{ProposalTypeUpdateAdmin},
		},
		"multiple proper items": {
			input:     "StoreCode,InstantiateContract,MigrateContract",
			proposals: []ProposalType{ProposalTypeStoreCode, ProposalTypeInstantiateContract, ProposalTypeMigrateContract},
		},
		"empty trailing item": {
			input:   "StoreCode,",
			isError: true,
		},
		"invalid item": {
			input:   "StoreCode,InvalidProposalType",
			isError: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			chunks := strings.Split(tc.input, ",")
			proposals, err := ConvertToProposals(chunks)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, proposals, tc.proposals)
			}
		})
	}
}

func TestUnmarshalContentFromJson(t *testing.T) {
	specs := map[string]struct {
		src string
		got govtypes.Content
		exp govtypes.Content
	}{
		"instantiate ": {
			src: `
{
	"title": "foo",
	"description": "bar",
	"admin": "myAdminAddress",
	"code_id": 1,
	"funds": [{"denom": "ALX", "amount": "2"},{"denom": "BLX","amount": "3"}],
	"init_msg": "e30=",
	"label": "testing",
	"run_as": "myRunAsAddress"
}`,
			got: &InstantiateContractProposal{},
			exp: &InstantiateContractProposal{
				Title:       "foo",
				Description: "bar",
				RunAs:       "myRunAsAddress",
				Admin:       "myAdminAddress",
				CodeID:      1,
				Label:       "testing",
				InitMsg:     []byte("{}"),
				Funds:       sdk.NewCoins(sdk.NewCoin("ALX", sdk.NewInt(2)), sdk.NewCoin("BLX", sdk.NewInt(3))),
			},
		},
		"migrate ": {
			src: `
{
	"title": "foo",
	"description": "bar",
	"code_id": 1,
	"contract": "myContractAddr",
	"migrate_msg": "e30=",
	"run_as": "myRunAsAddress"
}`,
			got: &MigrateContractProposal{},
			exp: &MigrateContractProposal{
				Title:       "foo",
				Description: "bar",
				RunAs:       "myRunAsAddress",
				Contract:    "myContractAddr",
				CodeID:      1,
				MigrateMsg:  []byte("{}"),
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(spec.src), spec.got))
			assert.Equal(t, spec.exp, spec.got)
		})
	}

}
