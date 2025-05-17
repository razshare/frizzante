/**
 * @typedef DonePayload
 * @property {function(string):(function(Record<string,any>):void)} findNavigateByPath
 * @property {Record<string,any>} data
 * @property {string} query
 */

/**
 * @param {DonePayload} payload
 */
function done(payload) {
    const {findNavigateByPath, data, query} = payload

    /**
     * @param {Response} response
     */
    return function (response) {
        if (response.status >= 300) {
            console.error(`Submit request failed with status ${response.status} ${response.statusText}.`)
            return
        }

        response.json()
            .then(function (responseData) {
                history.replaceState(
                    window.history.state ?? {},
                    "",
                    `${window.location.pathname}${document.location.hash}${query}`,
                )

                for (const key in data) {
                    delete data[key]
                }

                for (const key in responseData) {
                    data[key] = responseData[key]
                }

                if (response.redirected) {
                    const navigate = findNavigateByPath(response.url.replace(window.location.origin, ""))
                    navigate(responseData)
                }
            })
            .catch(fail)
    }
}

/**
 * @param {any} reason
 */
function fail(reason) {
    console.error("Submit request failed.", reason)
}

/**
 * @typedef UpdatePayload
 * @property {function(string):(function(Record<string,any>):void)} findNavigateByPath
 * @property {Record<string,any>} data
 */

/**
 * @param {UpdatePayload} payload
 */
export function update(payload) {
    const {findNavigateByPath, data} = payload
    return function onsubmit(e) {
        e.preventDefault()
        /** @type {HTMLFormElement} */
        const formElement = e.target
        const formData = new FormData(formElement)
        const method = formElement.method.toUpperCase()
        const headers = {"Accept": "application/json"}

        if ("GET" === method) {
            const dataLocal = new URLSearchParams();
            for (const [key, value] of formData) {
                dataLocal.append(key, value.toString());
            }

            let query = dataLocal.toString()
            if ('' !== query) {
                query = `?${query}`
            }

            fetch(`${formElement.action}${query}`, {method, headers})
                .then(done({findNavigateByPath, data, query})).catch(fail)
            return
        }

        fetch(formElement.action, {method, headers, body: formData})
            .then(done({findNavigateByPath, data, query: ""})).catch(fail)
    }
}