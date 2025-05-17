<script>
    //:app-imports
    import {setContext} from "svelte";

    /**
     * @typedef PageMetadata
     * @property {string} path
     * @property {string} viewName
     */

    /**
     * @typedef Props
     * @property {string} pageName
     * @property {Record<string,any>} data
     * @property {Record<string,PageMetadata>} pagesMetadata
     */

    /** @type {Props} */
    let {pageName, data, pagesMetadata} = $props()
    let pageNameState = $state(pageName)
    let dataState = $state({...data})
    let navCounterPrevious = 0
    setContext("data", dataState)
    setContext("navigate", navigate)
    setContext("findPathByPageName", findPathByPageName)
    setContext("findNavigateByPath", findNavigateByPath)

    /**
     * @param {string} pageName
     * @param {Record<string,string>} [parameters]
     * @param {false|Record<string,any>} [data]
     */
    function navigate(pageName, parameters, data = false) {
        swap(pageName, "push", parameters, data)
    }

    /**
     * @param {string} pageName
     * @param {Record<string,string>} [parameters]
     * @returns {string}
     */
    function findPathByPageName(pageName, parameters = {}) {
        if (!pagesMetadata[pageName]) {
            return ""
        }

        let result = pagesMetadata[pageName].path ?? ""

        for (let key in parameters) {
            const value = parameters[key]
            const regex = `{${key}}`.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
            result = result.replaceAll(new RegExp(regex, "g"), value)
        }

        return result
    }

    /**
     * @param {string} path
     * @returns {function(Record<string,any>):void}
     */
    function findNavigateByPath(path) {
        const partsGiven = path.split("/")
        for (const pageName in pagesMetadata) {
            const pathExpected = pagesMetadata[pageName].path
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
                    if (!expectedIsParameter) {
                        ok = false
                        break
                    }
                    const key = partsExpected[index].substring(0, partsExpected[index].length - 1).substring(1)
                    parameters[key] = partsGiven[index]
                } else if (expectedIsParameter) {
                    // Given part and expected part cannot be equal while expected part is a parameter.
                    // We reject that.
                    ok = false
                    break
                }
            }

            if (ok) {
                return function (data) {
                    swap(pageName, "push", parameters,  data)
                }
            }
        }

        return function () {
        }
    }

    /**
     *
     * @param {string} pageName
     * @param {"back"|"forward"|"push"} modifier
     * @param {Record<string,string>} [parameters]
     * @param {false|Record<string,any>} [data]
     */
    function swap(pageName, modifier, parameters, data = false) {
        if (!pagesMetadata[pageName]) {
            return
        }

        const pathLocal = findPathByPageName(pageName, parameters)
        if ("push" === modifier) {
            window.history.pushState({
                pageName,
                parameters,
                navCounter: ++navCounterPrevious,
            }, "", pathLocal);
        }
        pageNameState = pageName

        if (false !== data) {
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

    window.addEventListener("popstate", (e) => {
        e.preventDefault()
        const viewLocal = e.state?.pageName ?? ""
        const parameters = e.state?.parameters ?? {}
        const navCounterLocal = e.state?.navCounter ?? 0
        if (navCounterLocal < navCounterPrevious) {
            swap(viewLocal, "back", parameters)
            navCounterPrevious = navCounterLocal
        } else if (navCounterLocal > navCounterPrevious) {
            swap(viewLocal, "forward", parameters)
            navCounterPrevious = navCounterLocal
        } else {
            swap(viewLocal, "push", parameters)
        }
    });
</script>

<!--app-router-->