package obscuroscan

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	pathNumRollups    = "/numrollups/"
	pathLatestRollups = "/latestrollups/"
	pathLatestTxs     = "/latesttxs/"
	pathBlock         = "/block/"
	pathRollup        = "/rollup/"
	pathDecryptTxBlob = "/decrypttxblob/"
	staticDir         = "static"
	pathRoot          = "/"
	httpCodeErr       = 500
)

//go:embed static
var staticFiles embed.FS

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server      *http.Server
	client      rpcclientlib.Client
	contractABI abi.ABI
}

func NewObscuroscan(address string) *Obscuroscan {
	client := rpcclientlib.NewClient(address)
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic("could not parse management contract ABI to decrypt rollups")
	}
	return &Obscuroscan{
		client:      client,
		contractABI: contractABI,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathNumRollups, o.getNumRollups)       // Get the number of published rollups.
	serveMux.HandleFunc(pathLatestRollups, o.getLatestRollups) // Get the latest rollup numbers.
	serveMux.HandleFunc(pathLatestTxs, o.getLatestTxs)         // Get the latest transaction hashes.
	serveMux.HandleFunc(pathRollup, o.getRollupByNumOrTxHash)  // Get the rollup given its number or the hash of a transaction it contains.
	serveMux.HandleFunc(pathBlock, o.getBlock)                 // Get the L1 block with the given number.
	serveMux.HandleFunc(pathDecryptTxBlob, o.decryptTxBlob)    // Decrypt a transaction blob.

	// Serves the web assets for the user interface.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMux.Handle(pathRoot, http.FileServer(http.FS(noPrefixStaticFiles)))

	o.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err = o.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (o *Obscuroscan) Shutdown() {
	if o.server != nil {
		err := o.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down Obscuroscan. Cause: %s", err)
		}
	}
}

// Retrieves the number of published rollups.
func (o *Obscuroscan) getNumRollups(resp http.ResponseWriter, _ *http.Request) {
	numOfRollups, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}

	numOfRollupsStr := strconv.Itoa(int(numOfRollups))
	_, err = resp.Write([]byte(numOfRollupsStr))
	if err != nil {
		log.Error("could not return number of rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}
}

// Retrieves the last five rollup numbers.
func (o *Obscuroscan) getLatestRollups(resp http.ResponseWriter, _ *http.Request) {
	latestRollupNum, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}

	// We walk the chain of rollups, getting the number for the most recent five.
	rollupNums := make([]string, 5)
	for idx := 0; idx < 5; idx++ {
		rollupNum := int(latestRollupNum) - idx
		if rollupNum < 0 {
			// If there are less than five rollups, we return an N/A.
			rollupNums[idx] = "N/A"
		} else {
			rollupNums[idx] = strconv.Itoa(rollupNum)
		}
	}

	jsonRollupNums, err := json.Marshal(rollupNums)
	if err != nil {
		log.Error("could not return latest rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
	_, err = resp.Write(jsonRollupNums)
	if err != nil {
		log.Error("could not return latest rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
}

// Retrieves the last five transaction hashes.
func (o *Obscuroscan) getLatestTxs(resp http.ResponseWriter, _ *http.Request) {
	rollupNum, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}

	// We walk the chain of rollups, getting the transaction hashes until we've hit at least five.
	var txHashes []string
	for {
		rollup, err := o.getRollupByNumber(int(rollupNum))
		if err != nil {
			log.Error(err.Error())
			logAndSendErr(resp, "Could not fetch latest transactions.")
			return
		}

		for _, txHash := range rollup.TxHashes {
			txHashes = append(txHashes, txHash.String())
		}
		if len(txHashes) >= 5 {
			// It's fine if we return more than five transaction hashes; the front-end will ignore the later ones.
			break
		}

		rollupNum--
		if rollupNum < 0 {
			for idx := 0; idx < 5-len(txHashes); idx++ {
				txHashes = append(txHashes, "N/A")
			}
			break
		}
	}

	jsonTxHashes, err := json.Marshal(txHashes)
	if err != nil {
		log.Error("could not return latest transaction hashes to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}
	_, err = resp.Write(jsonTxHashes)
	if err != nil {
		log.Error("could not return latest transaction hashes to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}
}

// Retrieves the L1 block header with the given number.
func (o *Obscuroscan) getBlock(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	blockHashStr := buffer.String()
	blockHash := gethcommon.HexToHash(blockHashStr)

	var blockHeader *types.Header
	err = o.client.Call(&blockHeader, rpcclientlib.RPCGetBlockHeaderByHash, blockHash)
	if err != nil {
		log.Error("could not retrieve block with hash %s. Cause: %s", blockHash, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}

	jsonBlock, err := json.Marshal(blockHeader)
	if err != nil {
		log.Error("could not return block to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	_, err = resp.Write(jsonBlock)
	if err != nil {
		log.Error("could not return block to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
}

// Retrieves a rollup given its number or the hash of a transaction it contains.
func (o *Obscuroscan) getRollupByNumOrTxHash(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}

	var rollup *common.ExtRollup
	if strings.HasPrefix(buffer.String(), "0x") {
		// A "0x" prefix indicates that we should retrieve the rollup by transaction hash.
		txHash := gethcommon.HexToHash(buffer.String())

		err = o.client.Call(&rollup, rpcclientlib.RPCGetRollupForTx, txHash)
		if err != nil {
			log.Error("could not retrieve rollup. Cause: %s", err)
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
	} else {
		// Otherwise, we treat the input as a rollup number.
		rollupNumber, err := strconv.Atoi(buffer.String())
		if err != nil {
			log.Error("could not parse \"%s\" as an integer", buffer.String())
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
		rollup, err = o.getRollupByNumber(rollupNumber)
		if err != nil {
			log.Error(err.Error())
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
	}

	jsonRollup, err := json.Marshal(rollup)
	if err != nil {
		log.Error("could not return rollup to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		log.Error("could not return rollup to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}
}

// Decrypts the provided transaction blob using the provided key.
// TODO - Use the passed-in key, rather than a hardcoded enclave key.
func (o *Obscuroscan) decryptTxBlob(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	jsonTxs, err := decryptTxBlob(buffer.Bytes())
	if err != nil {
		log.Error("could not decrypt transaction blob. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	_, err = resp.Write(jsonTxs)
	if err != nil {
		log.Error("could not write decrypted transactions to client. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}
}

// Returns the number of the latest rollup.
func (o *Obscuroscan) getLatestRollupNumber() (int64, error) {
	var rollupHeader *common.Header
	err := o.client.Call(&rollupHeader, rpcclientlib.RPCGetCurrentRollupHead)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve head rollup. Cause: %w", err)
	}

	latestRollupNum := rollupHeader.Number.Int64()
	return latestRollupNum, nil
}

// Parses numberStr as a number and returns the rollup with that number.
func (o *Obscuroscan) getRollupByNumber(rollupNumber int) (*common.ExtRollup, error) {
	// TODO - If required, consolidate the two calls below into a single RPCGetRollupByNumber call to minimise round trips.
	var rollupHeader *common.Header
	err := o.client.Call(&rollupHeader, rpcclientlib.RPCGetRollupHeaderByNumber, rollupNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with number %d. Cause: %w", rollupNumber, err)
	}

	rollupHash := rollupHeader.Hash()
	if rollupHash == (gethcommon.Hash{}) {
		return nil, fmt.Errorf("rollup was retrieved but hash was nil")
	}

	var rollup *common.ExtRollup
	err = o.client.Call(&rollup, rpcclientlib.RPCGetRollup, rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup. Cause: %w", err)
	}
	return rollup, nil
}

// Decrypts the transaction blob and returns it as JSON.
func decryptTxBlob(encryptedTxBytesBase64 []byte) ([]byte, error) {
	encryptedTxBytes, err := base64.StdEncoding.DecodeString(string(encryptedTxBytesBase64))
	if err != nil {
		return nil, fmt.Errorf("could not decode encrypted transaction blob from Base64. Cause: %w", err)
	}

	key := gethcommon.Hex2Bytes(crypto.RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave rollup key. Cause: %w", err)
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %w", err)
	}

	encodedTxs, err := transactionCipher.Open(nil, []byte(crypto.RollupCipherNonce), encryptedTxBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted L2 transactions. Cause: %w", err)
	}

	var cleartextTxs []*common.L2Tx
	if err = rlp.DecodeBytes(encodedTxs, &cleartextTxs); err != nil {
		return nil, fmt.Errorf("could not decode encoded L2 transactions. Cause: %w", err)
	}

	jsonRollup, err := json.Marshal(cleartextTxs)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt transaction blob. Cause: %w", err)
	}

	return jsonRollup, nil
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
