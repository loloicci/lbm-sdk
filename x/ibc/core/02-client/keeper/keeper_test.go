package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	ostbytes "github.com/line/ostracon/libs/bytes"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	octypes "github.com/line/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/keeper"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	commitmenttypes "github.com/line/lbm-sdk/x/ibc/core/23-commitment/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	ibctmtypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	localhoctypes "github.com/line/lbm-sdk/x/ibc/light-clients/09-localhost/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	ibctestingmock "github.com/line/lbm-sdk/x/ibc/testing/mock"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

const (
	testChainID          = "gaiahub-0"
	testChainIDRevision1 = "gaiahub-1"

	testClientID  = "tendermint-0"
	testClientID2 = "tendermint-1"
	testClientID3 = "tendermint-2"

	height = 5

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testClientHeight          = types.NewHeight(0, 5)
	testClientHeightRevision1 = types.NewHeight(1, 5)
	newClientHeight           = types.NewHeight(1, 1)
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	cdc            codec.Marshaler
	ctx            sdk.Context
	keeper         *keeper.Keeper
	consensusState *ibctmtypes.ConsensusState
	header         *ibctmtypes.Header
	valSet         *octypes.ValidatorSet
	valSetHash     ostbytes.HexBytes
	privVal        octypes.PrivValidator
	now            time.Time
	past           time.Time

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	isCheckTx := false
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now2 := suite.now.Add(time.Hour)
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, ocproto.Header{Height: height, ChainID: testClientID, Time: now2})
	suite.keeper = &app.IBCKeeper.ClientKeeper
	suite.privVal = ibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	testClientHeightMinus1 := types.NewHeight(0, height-1)

	validator := ibctesting.NewTestValidator(pubKey, 1)
	suite.valSet = octypes.NewValidatorSet([]*octypes.Validator{validator})
	suite.valSetHash = suite.valSet.Hash()
	voterSet := octypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	suite.header = suite.chainA.CreateOCClientHeader(testChainID, int64(testClientHeight.RevisionHeight), testClientHeightMinus1, now2, suite.valSet, suite.valSet, voterSet, voterSet, []octypes.PrivValidator{suite.privVal})
	suite.consensusState = ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("hash")), suite.valSetHash)

	var validators stakingtypes.Validators
	for i := 1; i < 11; i++ {
		privVal := ibctestingmock.NewPV()
		tmPk, err := privVal.GetPubKey()
		suite.Require().NoError(err)
		pk, err := cryptocodec.FromOcPubKeyInterface(tmPk)
		suite.Require().NoError(err)
		val, err := stakingtypes.NewValidator(sdk.BytesToValAddress(pk.Address()), pk, stakingtypes.Description{})
		suite.Require().NoError(err)

		val.Status = stakingtypes.Bonded
		val.Tokens = sdk.NewInt(rand.Int63())
		validators = append(validators, val)

		hi := stakingtypes.NewHistoricalInfo(suite.ctx.BlockHeader(), validators)
		app.StakingKeeper.SetHistoricalInfo(suite.ctx, int64(i), &hi)
	}

	// add localhost client
	revision := types.ParseChainID(suite.chainA.ChainID)
	localHostClient := localhoctypes.NewClientState(
		suite.chainA.ChainID, types.NewHeight(revision, uint64(suite.chainA.GetContext().BlockHeight())),
	)
	suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), exported.Localhost, localHostClient)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.IBCKeeper.ClientKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetClientState() {
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	retrievedState, found := suite.keeper.GetClientState(suite.ctx, testClientID)
	suite.Require().True(found, "GetClientState failed")
	suite.Require().Equal(clientState, retrievedState, "Client states are not equal")
}

func (suite *KeeperTestSuite) TestSetClientConsensusState() {
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	retrievedConsState, found := suite.keeper.GetClientConsensusState(suite.ctx, testClientID, testClientHeight)
	suite.Require().True(found, "GetConsensusState failed")

	tmConsState, ok := retrievedConsState.(*ibctmtypes.ConsensusState)
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, tmConsState, "ConsensusState not stored correctly")
}

func (suite *KeeperTestSuite) TestValidateSelfClient() {
	testClientHeight := types.NewHeight(0, uint64(suite.chainA.GetContext().BlockHeight()-1))

	testCases := []struct {
		name        string
		clientState exported.ClientState
		expPass     bool
	}{
		{
			"success",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			true,
		},
		{
			"success with nil UpgradePath",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), nil, false, false),
			true,
		},
		{
			"invalid client type",
			localhoctypes.NewClientState(suite.chainA.ChainID, testClientHeight),
			false,
		},
		{
			"frozen client",
			&ibctmtypes.ClientState{suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false},
			false,
		},
		{
			"incorrect chainID",
			ibctmtypes.NewClientState("gaiatestnet", ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid client height",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, uint64(suite.chainA.GetContext().BlockHeight())), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid client revision",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeightRevision1, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid proof specs",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, nil, ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid trust level",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.Fraction{0, 1}, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid unbonding period",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod+10, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid trusting period",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, ubdPeriod+10, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid upgrade path",
			ibctmtypes.NewClientState(suite.chainA.ChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), []string{"bad", "upgrade", "path"}, false, false),
			false,
		},
	}

	for _, tc := range testCases {
		err := suite.chainA.App.IBCKeeper.ClientKeeper.ValidateSelfClient(suite.chainA.GetContext(), tc.clientState)
		if tc.expPass {
			suite.Require().NoError(err, "expected valid client for case: %s", tc.name)
		} else {
			suite.Require().Error(err, "expected invalid client for case: %s", tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestGetAllGenesisClients() {
	clientIDs := []string{
		testClientID2, testClientID3, testClientID,
	}
	expClients := []exported.ClientState{
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
	}

	expGenClients := make(types.IdentifiedClientStates, len(expClients))

	for i := range expClients {
		suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), clientIDs[i], expClients[i])
		expGenClients[i] = types.NewIdentifiedClientState(clientIDs[i], expClients[i])
	}

	// add localhost client
	localHostClient, found := suite.chainA.App.IBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), exported.Localhost)
	suite.Require().True(found)
	expGenClients = append(expGenClients, types.NewIdentifiedClientState(exported.Localhost, localHostClient))

	genClients := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllGenesisClients(suite.chainA.GetContext())

	suite.Require().Equal(expGenClients.Sort(), genClients)
}

func (suite KeeperTestSuite) TestGetAllGenesisMetadata() {
	expectedGenMetadata := []types.IdentifiedGenesisMetadata{
		types.NewIdentifiedGenesisMetadata(
			"clientA",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 1)), []byte("foo")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 2)), []byte("bar")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 3)), []byte("baz")),
			},
		),
		types.NewIdentifiedGenesisMetadata(
			"clientB",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(1, 100)), []byte("val1")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(2, 300)), []byte("val2")),
			},
		),
	}

	genClients := []types.IdentifiedClientState{
		types.NewIdentifiedClientState("clientA", &ibctmtypes.ClientState{}), types.NewIdentifiedClientState("clientB", &ibctmtypes.ClientState{}),
		types.NewIdentifiedClientState("clientC", &ibctmtypes.ClientState{}), types.NewIdentifiedClientState("clientD", &localhoctypes.ClientState{}),
	}

	suite.chainA.App.IBCKeeper.ClientKeeper.SetAllClientMetadata(suite.chainA.GetContext(), expectedGenMetadata)

	actualGenMetadata, err := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllClientMetadata(suite.chainA.GetContext(), genClients)
	suite.Require().NoError(err, "get client metadata returned error unexpectedly")
	suite.Require().Equal(expectedGenMetadata, actualGenMetadata, "retrieved metadata is unexpected")
}

func (suite KeeperTestSuite) TestGetConsensusState() {
	suite.ctx = suite.ctx.WithBlockHeight(10)
	cases := []struct {
		name    string
		height  types.Height
		expPass bool
	}{
		{"zero height", types.ZeroHeight(), false},
		{"height > latest height", types.NewHeight(0, uint64(suite.ctx.BlockHeight())+1), false},
		{"latest height - 1", types.NewHeight(0, uint64(suite.ctx.BlockHeight())-1), true},
		{"latest height", types.GetSelfHeight(suite.ctx), true},
	}

	for i, tc := range cases {
		tc := tc
		cs, found := suite.keeper.GetSelfConsensusState(suite.ctx, tc.height)
		if tc.expPass {
			suite.Require().True(found, "Case %d should have passed: %s", i, tc.name)
			suite.Require().NotNil(cs, "Case %d should have passed: %s", i, tc.name)
		} else {
			suite.Require().False(found, "Case %d should have failed: %s", i, tc.name)
			suite.Require().Nil(cs, "Case %d should have failed: %s", i, tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestConsensusStateHelpers() {
	// initial setup
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)

	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	nextState := ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("next")), suite.valSetHash)

	testClientHeightPlus5 := types.NewHeight(0, height+5)
	voterSet := octypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	header := suite.chainA.CreateOCClientHeader(testClientID, int64(testClientHeightPlus5.RevisionHeight), testClientHeight, suite.header.Header.Time.Add(time.Minute),
		suite.valSet, suite.valSet, voterSet, voterSet, []octypes.PrivValidator{suite.privVal})

	// mock update functionality
	clientState.LatestHeight = header.GetHeight().(types.Height)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, header.GetHeight(), nextState)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	latest, ok := suite.keeper.GetLatestClientConsensusState(suite.ctx, testClientID)
	suite.Require().True(ok)
	suite.Require().Equal(nextState, latest, "Latest client not returned correctly")
}

// 2 clients in total are created on chainA. The first client is updated so it contains an initial consensus state
// and a consensus state at the update height.
func (suite KeeperTestSuite) TestGetAllConsensusStates() {
	clientA, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Ostracon)

	clientState := suite.chainA.GetClientState(clientA)
	expConsensusHeight0 := clientState.GetLatestHeight()
	consensusState0, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight0)
	suite.Require().True(ok)

	// update client to create a second consensus state
	err := suite.coordinator.UpdateClient(suite.chainA, suite.chainB, clientA, exported.Ostracon)
	suite.Require().NoError(err)

	clientState = suite.chainA.GetClientState(clientA)
	expConsensusHeight1 := clientState.GetLatestHeight()
	suite.Require().True(expConsensusHeight1.GT(expConsensusHeight0))
	consensusState1, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight1)
	suite.Require().True(ok)

	expConsensus := []exported.ConsensusState{
		consensusState0,
		consensusState1,
	}

	// create second client on chainA
	clientA2, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Ostracon)
	clientState = suite.chainA.GetClientState(clientA2)

	expConsensusHeight2 := clientState.GetLatestHeight()
	consensusState2, ok := suite.chainA.GetConsensusState(clientA2, expConsensusHeight2)
	suite.Require().True(ok)

	expConsensus2 := []exported.ConsensusState{consensusState2}

	expConsensusStates := types.ClientsConsensusStates{
		types.NewClientConsensusStates(clientA, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight0.(types.Height), expConsensus[0]),
			types.NewConsensusStateWithHeight(expConsensusHeight1.(types.Height), expConsensus[1]),
		}),
		types.NewClientConsensusStates(clientA2, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight2.(types.Height), expConsensus2[0]),
		}),
	}.Sort()

	consStates := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllConsensusStates(suite.chainA.GetContext())
	suite.Require().Equal(expConsensusStates, consStates, "%s \n\n%s", expConsensusStates, consStates)
}
