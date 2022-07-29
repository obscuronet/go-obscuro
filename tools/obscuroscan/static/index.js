"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const typeClick = "click";
const methodPost = "POST";
const jsonKeyHeader = "Header";
const jsonKeyL1Proof = "L1Proof";
const jsonKeyEncryptedTxBlob = "EncryptedTxBlob";

const idNumRollups = "numRollups";
const idFormGetRollup = "form-get-rollup";
const idRollupID = "rollupID";
const idResult = "result";
const idBlock = "block";
const idRollup = "rollup";
const idDecryptedTxs = "decryptedTxs";
const idRollupOne = "rollupOne";
const idRollupTwo = "rollupTwo";
const idRollupThree = "rollupThree";
const idRollupFour = "rollupFour";
const idRollupFive = "rollupFive";
const idTxOne = "txOne";
const idTxTwo = "txTwo";
const idTxThree = "txThree";
const idTxFour = "txFour";
const idTxFive = "txFive";

const pathNumRollups = "/numrollups/";
const pathLatestRollups = "/latestrollups/";
const pathLatestTxs = "/latesttxs/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";
const pathDecryptTxBlob = "/decrypttxblob/";

// Updates the displayed stats.
async function updateStats() {
    const numRollupsField = document.getElementById(idNumRollups);

    const numRollupsResp = await fetch(pathNumRollups);

    if (numRollupsResp.ok) {
        numRollupsField.innerText = "Total rollups: " + await numRollupsResp.text();
    } else {
        numRollupsField.innerText = "Failed to fetch number of rollups. Cause: " + await numRollupsResp.text()
    }
}

// Updates the list of latest rollups.
async function updateLatestRollups() {
    const rollupOneField = document.getElementById(idRollupOne);
    const rollupTwoField = document.getElementById(idRollupTwo);
    const rollupThreeField = document.getElementById(idRollupThree);
    const rollupFourField = document.getElementById(idRollupFour);
    const rollupFiveField = document.getElementById(idRollupFive);

    const latestRollupsResp = await fetch(pathLatestRollups);

    if (latestRollupsResp.ok) {
        const latestRollupsJSON = JSON.parse(await latestRollupsResp.text());
        rollupOneField.innerText = latestRollupsJSON[0];
        rollupTwoField.innerText = latestRollupsJSON[1];
        rollupThreeField.innerText = latestRollupsJSON[2];
        rollupFourField.innerText = latestRollupsJSON[3];
        rollupFiveField.innerText = latestRollupsJSON[4];
    } else {
        rollupOneField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
        rollupTwoField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
        rollupThreeField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
        rollupFourField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
        rollupFiveField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
    }
}

// Updates the list of latest transactions.
async function updateLatestTxs() {
    const txOneField = document.getElementById(idTxOne);
    const txTwoField = document.getElementById(idTxTwo);
    const txThreeField = document.getElementById(idTxThree);
    const txFourField = document.getElementById(idTxFour);
    const txFiveField = document.getElementById(idTxFive);

    const latestTxsResp = await fetch(pathLatestTxs);

    if (latestTxsResp.ok) {
        const latestTxsJSON = JSON.parse(await latestTxsResp.text());
        txOneField.innerText = latestTxsJSON[0];
        txTwoField.innerText = latestTxsJSON[1];
        txThreeField.innerText = latestTxsJSON[2];
        txFourField.innerText = latestTxsJSON[3];
        txFiveField.innerText = latestTxsJSON[4];
    } else {
        txOneField.innerText = "Failed to fetch latest rollups. Cause: " + await latestTxsResp.text();
        txTwoField.innerText = "Failed to fetch latest rollups. Cause: " + await latestTxsResp.text();
        txThreeField.innerText = "Failed to fetch latest rollups. Cause: " + await latestTxsResp.text();
        txFourField.innerText = "Failed to fetch latest rollups. Cause: " + await latestTxsResp.text();
        txFiveField.innerText = "Failed to fetch latest rollups. Cause: " + await latestTxsResp.text();
    }
}

// Displays the rollup based on the rollup number or transaction hash that has been clicked on.
async function displayClickedItem(event) {
    event.preventDefault();
    const rollupNumber = document.getElementById(event.target.id).innerText;
    await displayRollup(rollupNumber);
}

// Displays the rollup with the given identifier (either a rollup number or a transaction hash that it contains).
async function displayRollup(rollupID) {
    const resultPane = document.getElementById(idResult);
    const blockArea = document.getElementById(idBlock);
    const rollupArea = document.getElementById(idRollup);
    const decryptedTxsArea = document.getElementById(idDecryptedTxs);

    const rollupResp = await fetch(pathRollup, {
        method: methodPost,
        body: rollupID
    });

    if (!rollupResp.ok) {
        rollupArea.innerText = "Failed to fetch rollup. Cause: " + await rollupResp.text()
        blockArea.innerText = "Failed to fetch block. Cause: Could not retrieve rollup."
        return
    }

    const rollupJSON = JSON.parse(await rollupResp.text())
    rollupArea.innerText = JSON.stringify(rollupJSON, null, "\t");

    const blockResp = await fetch(pathBlock, {
        method: methodPost,
        body: rollupJSON[jsonKeyHeader][jsonKeyL1Proof]
    });

    if (blockResp.ok) {
        const json = JSON.parse(await blockResp.text())
        blockArea.innerText = JSON.stringify(json, null, "\t");
    } else {
        blockArea.innerText = "Failed to fetch block. Cause: " + await blockResp.text()
    }

    const encryptedTxBlob = rollupJSON[jsonKeyEncryptedTxBlob]
    const decryptTxBlobResp = await fetch(pathDecryptTxBlob, {
        method: methodPost,
        body: encryptedTxBlob
    });

    if (decryptTxBlobResp.ok) {
        const json = JSON.parse(await decryptTxBlobResp.text())
        decryptedTxsArea.innerText = JSON.stringify(json, null, "\t");
    } else {
        decryptedTxsArea.innerText = "Failed to decrypt transaction blob. Cause: " + await decryptTxBlobResp.text()
    }

    resultPane.scrollIntoView();
}

const initialize = () => {
    // Updates the page's stats and latest rollups and transactions.
    setInterval(async () => {
        await updateStats();
        await updateLatestRollups();
        await updateLatestTxs();
    }, 1000);

    // Handles searches for rollups.
    document.getElementById(idFormGetRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();
        const rollupNumber = document.getElementById(idRollupID).value;
        await displayRollup(rollupNumber);
    });

    // Handles clicks on the latest rollups or transactions.
    document.getElementById(idRollupOne).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupTwo).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupThree).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupFour).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupFive).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxOne).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxTwo).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxThree).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxFour).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxFive).addEventListener(typeClick, displayClickedItem);
}

window.addEventListener(eventDomLoaded, initialize);