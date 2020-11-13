package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	tmjson "github.com/vbhp/supermint/libs/json"
	"github.com/vbhp/supermint/privval"
	"github.com/vbhp/supermint/types"
)

// GenValidatorCmd allows the generation of a keypair for a
// validator.
var GenValidatorCmd = &cobra.Command{
	Use:   "gen_validator",
	Short: "Generate new validator keypair",
	Run:   genValidator,
}

func init() {
	GenValidatorCmd.Flags().StringVar(&keyType, "key", types.ABCIPubKeyTypeEd25519,
		"Key type to generate privval file with. Options: ed25519, secp256k1")
}

func genValidator(cmd *cobra.Command, args []string) {
	pv, err := privval.GenFilePV("", "", keyType)
	if err != nil {
		panic(err)
	}
	jsbz, err := tmjson.Marshal(pv)
	if err != nil {
		panic(err)
	}
	fmt.Printf(`%v
`, string(jsbz))
}
