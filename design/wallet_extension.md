# Wallet extension design

## Scope

The design for the wallet extension, a component that is responsible for handling RPC requests from traditional Ethereum 
wallets (e.g. MetaMask, hardware wallets) and webapps to the Obscuro host.

## Requirements

* The wallet extension serves an endpoint that meets the [Ethereum JSON-RPC specification
  ](https://playground.open-rpc.org/?schemaUrl=https://raw.githubusercontent.com/ethereum/eth1.0-apis/assembled-spec/openrpc.json)
* The wallet extension is run locally by the end user
* The wallet extension does not broadcast any sensitive information (e.g. transactions, balances) in plaintext
* The encryption keys are only known to the client (e.g. wallet, webapp), and not to any third-parties
* The encryption is transparent to the client; from the client's perspective, they are interacting with a "standard" 
  non-encrypting implementation of the Ethereum JSON-RPC specification
* The wallet extension is usable by any webapp or wallet type, and in particular:
  * Hardware wallets that do not offer a decryption capability
  * MetaMask, for which the keys are only available when running in the browser

## Design

The wallet extension is a local server application that serves two endpoints:

* An endpoint for managing *viewing keys*
* An endpoint that meets the Ethereum JSON-RPC specification

The wallet extension also maintains an RPC connection to one or more Obscuro hosts.

### Viewing-keys endpoint

This endpoint serves a webpage where the end user can generate new viewing keys for their account. For each generation, 
the following steps are taken:

* The wallet extension generates a new keypair
* The wallet extension stores the private key locally
* The end user signs a payload containing the public key and some metadata using their wallet
* The wallet extension sends the public key to the Obscuro enclave via the Obscuro host over RPC

Whenever an enclave needs to send sensitive information to the end user (e.g. a transaction result or account balance), 
it encrypts the sensitive information with the viewing key of the account.

This ensures that the sensitive information can only be decrypted by the wallet extension. By generating new viewing 
keys through a webpage, we maintain compatibility with MetaMask.

If multiple viewing keys are registered for a single account, a separate encrypted payload is sent for each viewing key.

This endpoint will also have to handle the expiry of viewing keys.

### Ethereum JSON-RPC endpoint

The wallet extension serves a standard implementation of the Ethereum JSON-RPC specification, except in the following 
respects:

* Any request containing sensitive information is encrypted with the Obscuro enclave public key before being forwarded 
  to the Obscuro host
* Any response that is encrypted with a viewing key is decrypted with the corresponding private key before being 
  forwarded to the end user

This ensures that the encryption and decryption involved in the Obscuro protocol is transparent to the end user, and 
that we are not relying on decryption capabilities being available in the wallet.

## Known limitations

* An additional set of keys, the viewing keys, must be managed outside of the user's wallet
  * Note, however, that these keys are far less sensitive than the signing keys; they only allow an attacker to view, 
    but not modify, the user's ledger
* The end user must run an additional component; this precludes mobile-only wallets

## Alternatives considered

### Alternatives to a local server application

#### Chrome and Firefox extension

This is appealing because the user experience of installing and running a browser extension is better than the user 
experience of installing and running a local application.

However, this approach is unworkable because Chrome extensions can only make outbound network requests, and cannot 
serve an endpoint. Chrome has a [`chrome.socket` API](https://developer.chrome.com/docs/extensions/reference/socket/), 
but it is only available to Chrome apps, and not extensions.

The penetration of Firefox was considered to be too low for a Firefox-extension-only approach to be considered viable.

#### Chrome app

Chrome apps are deprecated, and their support status is uncertain (see 
[here](https://blog.chromium.org/2021/10/extending-chrome-app-support-on-chrome.html)).

#### MetaMask snap

Snaps are an experimental MetaMask capability.

Snaps have several downsides:

* Snaps need to be installed per-page, requiring a code change in every webapp to prompt the user to install the Obscuro 
  snap
* Snaps are only compatible with MetaMask
* Snaps are marked as experimental and require users to switch from MetaMask to the experimental MetaMask Flask

### Alternatives to using viewing keys

Instead of viewing keys, we could encrypt and decrypt using the end user's existing account key.

This would obviate the need for us to implement a key management solution.

However, it has several downsides:

* It does not support hardware wallets, since hardware wallets do not provide an API to decrypt using their account 
  private key
* It may not be possible with MetaMask - the local server application would somehow have to speak to the browser to 
  access the MetaMask API; this may be possible through a combination of a Chrome extension and websockets, but is 
  unproven
* It would lead to a high number of MetaMask prompts
