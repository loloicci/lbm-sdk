module github.com/cosmos/cosmos-sdk

require (
	github.com/99designs/keyring v1.1.4
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/cosmos/ledger-cosmos-go v0.11.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1-0.20190508161146-9fa652df1129
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pelletier/go-toml v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/afero v1.2.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.13.2
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.8
)

go 1.13

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
