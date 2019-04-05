"use strict";

const constants = {
    rowsCountID: 'rows_count',
    columnsCountID: 'columns_count',
    notificationPanelID: 'notification_panel',
    priceMatrixSizeLimit: 5,
    fadeOutTimeMS: 2000,
    errors: {
        INFO: "alert alert-info",
        ERROR: "alert alert-danger",
    },
    apiURL: '/transport-issue/',
}

const showError = (message) => showNotification(constants.errors.ERROR, message)
const showInfo = (message) => showNotification(constants.errors.INFO, message)

function showNotification(className, message) {
    const tempDivId = new Date().getTime();
    const notificationPanel = document.getElementById(constants.notificationPanelID)
    const divElement = document.createElement("div");       
    divElement.id = tempDivId;
    divElement.className = className;
    divElement.innerText = message;
    divElement.setAttribute('role', 'alert');
    notificationPanel.appendChild(divElement);                    
    setTimeout(() => {
        console.debug(`expire notification: ${tempDivId}`)
        notificationPanel.removeChild(divElement)
    }, constants.fadeOutTimeMS);

}


function createMatrix() {
    const rowsCountInput = document.getElementById(constants.rowsCountID)
    const columnCountInput = document.getElementById(constants.columnsCountID)

    const rowsCount = Number(rowsCountInput.value)
    const columnCount = Number(columnCountInput.value)

    if (rowsCount === NaN) {
        showError("rows count is not a number value")
        return 
    }

    if (rowsCount < constants.priceMatrixSizeLimit) {
        showError(
            `rows count should be equal or greater than ${constants.priceMatrixSizeLimit}. Got ${rowsCount}`
        )
        return 
    }

    if (columnCount === NaN) {
        showError("rows count is not a number value")
        return 
    }

    if (columnCount < constants.priceMatrixSizeLimit) {
        showError(
            `column count should be equal or greater than ${constants.priceMatrixSizeLimit}. Got ${columnCount}`
        )
        return 
    }

    showInfo(`create Matrix: ${rowsCount}x${columnCount}`);
    createPriceHTMLTable()
}


async function solveIssueRequest() {
    const requestData = readValuesFromTable()

    const requestOptions = {
        method: "POST",
        headers: {'Accept': 'application/json', 'Content-Type': 'application/json'},
        body: JSON.stringify(requestData)
     }
      const apiURL = constants.apiURL;
      try {
        const response = await fetch(apiURL, requestOptions);
        if (response.status !== 200) {
          const errorText = await response.text();
          showError(`[${response.status}]: ${errorText}`)
          return 
        }
        const data = await response.text();
        showInfo(data)
      } catch(err) {
        showError(err)
      }
}


function createPriceHTMLTable () {

}
function readValuesFromTable () {
    // Example data
    const requestData = {
        consumers_needs: [20, 30, 30, 10],
             producers_sources: [30, 40, 20],
             prices: [
                 [2, 3, 2, 4],
                 [3, 2, 5, 1],
                 [4, 3, 2, 6],
        ],
    }

    // TODO fetch it from price table
    // const requestData = {
    //     consumers_needs: [],
    //     producers_sources: [],
    //     prices: [],
    // }
    return requestData
}