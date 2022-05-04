const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idGenerateViewingKey = "generateViewingKey";
const idStatus = "status";
const pathGenerateViewingKey = "/generateviewingkey/";
const pathSubmitViewingKey = "/submitviewingkey/";
const methodPost = "post";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";
const personalSignPrefix = "vk";

const initialize = () => {
    const generateViewingKeyButton = document.getElementById(idGenerateViewingKey);
    const statusArea = document.getElementById(idStatus);

    generateViewingKeyButton.addEventListener(eventClick, async () => {
        const viewingKeyResp = await fetch(pathGenerateViewingKey);
        if (!isOk(viewingKeyResp)) {
            statusArea.innerText = "Failed to generate viewing key."
            return
        }

        const viewingKey = await viewingKeyResp.text();

        const accounts = await ethereum.request({method: metamaskRequestAccounts});
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        const account = accounts[0]; // TODO - Allow use of accounts other than the first.

        const signedBytes = await ethereum.request({
            method: metamaskPersonalSign,
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: [personalSignPrefix + viewingKey, account]
        }).catch(_ => { return -1 })
        if (signedBytes === -1) {
            statusArea.innerText = "Failed to sign viewing key."
            return
        }

        const signedViewingKeyJson = {"signedBytes": signedBytes}
        const submitViewingKeyResp = await fetch(
            pathSubmitViewingKey, {
                method: methodPost,
                headers: jsonHeaders,
                body: JSON.stringify(signedViewingKeyJson)
            }
        );
        if (isOk(submitViewingKeyResp)) {
            statusArea.innerText = `Account: ${account}\nViewing key: ${viewingKey}\nSigned bytes: ${signedBytes}`
        } else {
            statusArea.innerText = "Failed to submit viewing key to enclave."
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);

function isOk(response) {
    return response.status >= 200 && response.status < 300
}