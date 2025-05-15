/**
 * @typedef DonePayload
 * @property {function(string):{view:string,parameters:Record<string,string>}} view
 * @property {function(string,Record<string,string>,false|Record<string,any>)} navigate
 * @property {string} query
 * @property {Record<string,any>} data
 */

/**
 *
 * @param {DonePayload} payload
 */
function done(payload) {
    const {
        view,
        navigate,
        query,
        data,
    } = payload

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
                    const resolved = view(response.url.replace(window.location.origin, ""))
                    navigate(resolved.view, resolved.parameters, responseData)
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
 * @property {function(string):{view:string,parameters:Record<string,string>}} view
 * @property {function(string,Record<string,string>,false|Record<string,any>):void} navigate
 * @property {Record<string,any>} data
 */

/**
 * @param {UpdatePayload} payload
 */
export function update(payload) {
    const {view, navigate, data} = payload
    return function onsubmit(e) {
        e.preventDefault()
        /** @type {HTMLFormElement} */
        const formElement = e.target
        const formData = new FormData(formElement)
        const method = formElement.method
        const headers = {"Accept": "application/json"}

        if (method === "get" || method === "GET") {
            const dataLocal = new URLSearchParams();
            for (const [key, value] of formData) {
                dataLocal.append(key, value.toString());
            }

            const search = data.toString()
            const query = `?${search}`
            const init = {method, headers}
            const donePayload = {
                view,
                navigate,
                query,
                data,
            }

            fetch(query, init).then(done(donePayload)).catch(fail)
            return
        }

        const init = {method, headers, body: formData}
        const donePayload = {
            view,
            navigate,
            query: "",
            data,
        }

        fetch(formElement.action, init).then(done(donePayload)).catch(fail)
    }
}