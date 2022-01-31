package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/security"
	"github.com/colonyos/colonies/pkg/security/crypto"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	runtimeCmd.AddCommand(registerRuntimeCmd)
	runtimeCmd.AddCommand(lsRuntimesCmd)
	runtimeCmd.AddCommand(approveRuntimeCmd)
	runtimeCmd.AddCommand(rejectRuntimeCmd)
	runtimeCmd.AddCommand(deleteRuntimeCmd)
	rootCmd.AddCommand(runtimeCmd)

	runtimeCmd.PersistentFlags().StringVarP(&ColonyID, "colonyid", "", "", "Colony Id")
	runtimeCmd.PersistentFlags().StringVarP(&ServerHost, "host", "", "localhost", "Server host")
	runtimeCmd.PersistentFlags().IntVarP(&ServerPort, "port", "", 8080, "Server HTTP port")

	registerRuntimeCmd.Flags().StringVarP(&ColonyPrvKey, "colonyprvkey", "", "", "Colony private key")
	registerRuntimeCmd.Flags().StringVarP(&SpecFile, "spec", "", "", "JSON specification of a Colony Runtime")
	registerRuntimeCmd.MarkFlagRequired("spec")

	lsRuntimesCmd.Flags().BoolVarP(&JSON, "json", "", false, "Print JSON instead of tables")
	lsRuntimesCmd.Flags().StringVarP(&RuntimeID, "runtimeid", "", "", "Runtime Id")
	lsRuntimesCmd.Flags().StringVarP(&RuntimePrvKey, "runtimeprvkey", "", "", "Runtime private key")

	approveRuntimeCmd.Flags().StringVarP(&ColonyPrvKey, "colonyprvkey", "", "", "Colony private key")
	approveRuntimeCmd.Flags().StringVarP(&RuntimeID, "runtimeid", "", "", "Colony Runtime Id")
	approveRuntimeCmd.MarkFlagRequired("runtimeid")

	rejectRuntimeCmd.Flags().StringVarP(&ColonyPrvKey, "colonyprvkey", "", "", "Colony private key")
	rejectRuntimeCmd.Flags().StringVarP(&RuntimeID, "runtimeid", "", "", "Colony Runtime Id")
	rejectRuntimeCmd.MarkFlagRequired("runtimeid")

	deleteRuntimeCmd.Flags().StringVarP(&ColonyPrvKey, "colonyprvkey", "", "", "Colony private key")
	deleteRuntimeCmd.Flags().StringVarP(&RuntimeID, "runtimeid", "", "", "Colony Runtime Id")
	deleteRuntimeCmd.MarkFlagRequired("runtimeid")
}

var runtimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "Manage Colony Runtimes",
	Long:  "Manage Colony Runtimes",
}

var registerRuntimeCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Runtime",
	Long:  "Register a new Runtime",
	Run: func(cmd *cobra.Command, args []string) {
		jsonSpecBytes, err := ioutil.ReadFile(SpecFile)
		CheckError(err)

		if ColonyID == "" {
			ColonyID = os.Getenv("COLONYID")
		}
		if ColonyID == "" {
			CheckError(errors.New("Unknown Colony Id"))
		}

		runtime, err := core.ConvertJSONToRuntime(string(jsonSpecBytes))
		CheckError(err)

		keychain, err := security.CreateKeychain(KEYCHAIN_PATH)
		CheckError(err)

		crypto := crypto.CreateCrypto()

		prvKey, err := crypto.GeneratePrivateKey()
		CheckError(err)

		runtimeID, err := crypto.GenerateID(prvKey)
		CheckError(err)
		runtime.SetID(runtimeID)
		runtime.SetColonyID(ColonyID)

		if ColonyPrvKey == "" {
			ColonyPrvKey, err = keychain.GetPrvKey(ColonyID)
			CheckError(err)
		}

		client := client.CreateColoniesClient(ServerHost, ServerPort, true) // XXX: Insecure
		addedRuntime, err := client.AddRuntime(runtime, ColonyPrvKey)
		CheckError(err)

		err = keychain.AddPrvKey(runtimeID, prvKey)
		CheckError(err)

		fmt.Println(addedRuntime.ID)
	},
}

var lsRuntimesCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all Runtimes available in a Colony",
	Long:  "List all Runtimes available in a Colony",
	Run: func(cmd *cobra.Command, args []string) {
		keychain, err := security.CreateKeychain(KEYCHAIN_PATH)
		CheckError(err)

		if ColonyID == "" {
			ColonyID = os.Getenv("COLONYID")
		}
		if ColonyID == "" {
			CheckError(errors.New("Unknown Colony Id"))
		}

		if RuntimePrvKey == "" {
			if RuntimeID == "" {
				RuntimeID = os.Getenv("RUNTIMEID")
			}
			RuntimePrvKey, _ = keychain.GetPrvKey(RuntimeID)
		}

		if RuntimePrvKey == "" {
			if RuntimeID == "" {
				RuntimeID = os.Getenv("RUNTIMEID")
			}
			RuntimePrvKey, _ = keychain.GetPrvKey(RuntimeID)
		}

		client := client.CreateColoniesClient(ServerHost, ServerPort, true) // XXX: Insecure
		runtimesFromServer, err := client.GetRuntimes(ColonyID, RuntimePrvKey)
		if err != nil {
			// Try ColonyPrvKey instead
			if ColonyPrvKey == "" {
				if ColonyID == "" {
					ColonyID = os.Getenv("COLONYID")
				}
				ColonyPrvKey, err = keychain.GetPrvKey(ColonyID)
				CheckError(err)
			}
			runtimesFromServer, err = client.GetRuntimes(ColonyID, ColonyPrvKey)
			CheckError(err)
		}

		if JSON {
			jsonString, err := core.ConvertRuntimeArrayToJSON(runtimesFromServer)
			CheckError(err)
			fmt.Println(jsonString)
			os.Exit(0)
		}

		var data [][]string
		for _, runtime := range runtimesFromServer {
			state := ""
			switch runtime.State {
			case core.PENDING:
				state = "Pending"
			case core.APPROVED:
				state = "Approved"
			case core.REJECTED:
				state = "Rejected"
			default:
				state = "Unknown"
			}

			data = append(data, []string{runtime.ID, runtime.Name, state})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "State"})

		for _, v := range data {
			table.Append(v)
		}

		table.Render()
	},
}

var approveRuntimeCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve a Colony Runtime",
	Long:  "Approve a Colony Runtime",
	Run: func(cmd *cobra.Command, args []string) {
		keychain, err := security.CreateKeychain(KEYCHAIN_PATH)
		CheckError(err)

		if ColonyID == "" {
			ColonyID = os.Getenv("COLONYID")
		}
		if ColonyID == "" {
			CheckError(errors.New("Unknown Colony Id"))
		}

		if ColonyPrvKey == "" {
			ColonyPrvKey, err = keychain.GetPrvKey(ColonyID)
			CheckError(err)
		}

		client := client.CreateColoniesClient(ServerHost, ServerPort, true) // XXX: Insecure
		err = client.ApproveRuntime(RuntimeID, ColonyPrvKey)
		CheckError(err)

		fmt.Println("Colony Runtime with Id <" + RuntimeID + "> is now approved")
	},
}

var rejectRuntimeCmd = &cobra.Command{
	Use:   "reject",
	Short: "Reject a Colony Runtime",
	Long:  "Reject a Colony Runtime",
	Run: func(cmd *cobra.Command, args []string) {
		keychain, err := security.CreateKeychain(KEYCHAIN_PATH)
		CheckError(err)

		if ColonyID == "" {
			ColonyID = os.Getenv("COLONYID")
		}
		if ColonyID == "" {
			CheckError(errors.New("Unknown Colony Id"))
		}

		if ColonyPrvKey == "" {
			ColonyPrvKey, err = keychain.GetPrvKey(ColonyID)
			CheckError(err)
		}

		client := client.CreateColoniesClient(ServerHost, ServerPort, true) // XXX: Insecure
		err = client.RejectRuntime(RuntimeID, ColonyPrvKey)
		CheckError(err)

		fmt.Println("Colony Runtime with Id <" + RuntimeID + "> is now rejected")
	},
}

var deleteRuntimeCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister a Colony Runtime",
	Long:  "Unregister a Colony Runtime",
	Run: func(cmd *cobra.Command, args []string) {
		keychain, err := security.CreateKeychain(KEYCHAIN_PATH)
		CheckError(err)

		if ColonyID == "" {
			ColonyID = os.Getenv("COLONYID")
		}
		if ColonyID == "" {
			CheckError(errors.New("Unknown Colony Id"))
		}

		if ColonyPrvKey == "" {
			ColonyPrvKey, err = keychain.GetPrvKey(ColonyID)
			CheckError(err)
		}

		client := client.CreateColoniesClient(ServerHost, ServerPort, true) // XXX: Insecure
		err = client.DeleteRuntime(RuntimeID, ColonyPrvKey)
		CheckError(err)

		fmt.Println("Colony Runtime with Id <" + RuntimeID + "> is now unregistered")
	},
}
