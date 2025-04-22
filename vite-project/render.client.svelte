<script>
    //:app-imports
    import {setContext} from "svelte";

    /**
     * @typedef Props
     * @property {string} view
     * @property {Record<string,any>} data
     * @property {Record<string,string>} views
     * @property {Record<string,string>} parameters
     */

    /** @type {Props} */
    let {view, data, views, parameters} = $props()
    let viewState = $state(view)
    let dataState = $state({...data})
    let navCounterPrevious = 0
    setContext("data", dataState)
    setContext("navigate",
        /**
         * @param {string} view
         * @param {Record<string,string>} [parameters]
         * @param {false|Record<string,any>} [data]
         */
        function (view, parameters, data = false) {
            navigate(view, "push", parameters, data)
        }
    )
    setContext("path", path)
    setContext("view", _view)

    window.history.replaceState({
        ...(window.history.state ?? {}),
        view,
        parameters,
        navCounter: navCounterPrevious,
    }, "", `${document.location.pathname}${document.location.hash}${document.location.search}`)

    window.addEventListener("popstate", (e) => {
        e.preventDefault()
        const viewLocal = e.state?.view ?? ""
        const parameters = e.state?.parameters ?? {}
        const navCounterLocal = e.state?.navCounter ?? 0
        if (navCounterLocal < navCounterPrevious) {
            navigate(viewLocal, "back", parameters)
            navCounterPrevious = navCounterLocal
        } else if (navCounterLocal > navCounterPrevious) {
            navigate(viewLocal, "forward", parameters)
            navCounterPrevious = navCounterLocal
        } else {
            navigate(viewLocal, "push", parameters)
        }
    });

    /**
     * @param {string} string
     */
    function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
    }

    /**
     * @param {string} view
     * @param {Record<string,string>} [parameters]
     */
    function path(view, parameters = {}) {
        let result = views[view] ?? ""
        if (!views[view]) {
            return ""
        }

        for (let key in parameters) {
            const value = parameters[key]
            const regex = escapeRegExp(`{${key}}`)
            result = result.replaceAll(new RegExp(regex, "g"), value)
        }

        return result
    }

    /**
     * @param {string} path
     * @returns {{view:string,parameters:Record<string,string>}}
     */
    function _view(path) {
        const partsGiven = path.split("/")
        for (const view in views) {
            const pathExpected = views[view]
            const partsExpected = pathExpected.split("/")
            if (partsExpected.length !== partsGiven.length) {
                continue
            }

            /** @type {Record<string,string>} */
            const parameters = {}

            let ok = true
            for (let index = 0; index < partsExpected.length; index++) {
                const expectedIsParameter = partsExpected[index].startsWith("{") && partsExpected[index].endsWith("}")
                const givenAndExpectedAreDifferent = partsGiven[index] !== partsExpected[index]

                if (givenAndExpectedAreDifferent) {
                    if(!expectedIsParameter){
                        ok = false
                        break
                    }
                    const key = partsExpected[index].substring(0,partsExpected[index].length-1).substring(1)
                    parameters[key] = partsGiven[index]
                } else if(expectedIsParameter) {
                    // Given part and expected part cannot be equal while expected part is a parameter.
                    // We reject that.
                    ok = false
                    break
                }
            }

            if (ok) {
                return {
                    view,
                    parameters,
                }
            }
        }

        return {
            view: "",
            parameters: {}
        }
    }

    /**
     *
     * @param {string} view
     * @param {"back"|"forward"|"push"} modifier
     * @param {Record<string,string>} [parameters]
     * @param {false|Record<string,any>} [data]
     */
    function navigate(view, modifier, parameters, data = false) {
        if (!views[view]) {
            return
        }

        const pathLocal = path(view, parameters)
        if ("push" === modifier) {
            window.history.pushState({
                view,
                parameters,
                navCounter: ++navCounterPrevious,
            }, "", pathLocal);
        }
        viewState = view

        if(false !== data){
            return
        }

        fetch(pathLocal, {headers: {"Accept": "application/json"}}).then(async (response) => {
            const data = await response.json()

            for (const key in dataState) {
                delete dataState[key]
            }

            for (const key in data) {
                dataState[key] = data[key]
            }
        })
    }
</script>

<!--app-router-->