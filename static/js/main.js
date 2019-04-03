const constants = {
    rowsCountID: 'rows_count',
    columnsCountID: 'columns_count',
    notificationPanelID: 'notification_panel',
    priceMatrixSizeLimit: 3,
    fadeOutTimeMS: 2000,
    errors: {
        INFO: "alert alert-info",
        ERROR: "alert alert-danger",
    }
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
        console.log(`expire notification: ${tempDivId}`)
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

    showInfo(`create Matrix: ${rowsCount}x${columnCount}`)
}

// var app = new Vue({
//     el: '#app',
//     delimiters: ['${', '}'],
//     data: {
//       message: 'Hello Vue!'
//     },
//       methods: {
//         reverseMessage: async function () {
//           const requestData = {
//             consumers_needs: [20, 30, 30, 10],
//                  producers_sources: [30, 40, 20],
//                  prices: [
//                      [2, 3, 2, 4],
//                      [3, 2, 5, 1],
//               [4, 3, 2, 6],
//             ],
//           }
//           const requestOptions = {
//             method: "POST",
//             headers: {'Accept': 'application/json', 'Content-Type': 'application/json'},
//             body: JSON.stringify(requestData)
//          }
//           const apiURL = `/transport-issue/`;
//           //const apiURL = `'https://api.coindesk.com/v1/bpi/currentprice.json'`
//           try {
//             //debugger;
//             const response = await fetch(apiURL, requestOptions);
//             if (response.status !== 200) {
//               const errorText = await response.text();
//               console.log('res text', errorText)
//               this.message = errorText;
//               return 
//             }

//             const data = await response.text();
//             console.log('res', data)
//             this.message = data;
//           } catch(err) {
//             console.log('err', err)
//             this.message = "error";
//           }

//         }
//   }
//   })