<script>
    //:app-imports
    import {setContext} from 'svelte'

    /**
     * @typedef Props
     * @property {string} view
     * @property {Record<string,any>} data
     * @property {Record<string,string>} views
     * @property {Record<string,string>} parameters
     */

    // Do not remove or discard `pageId`, it's being used by app-router.
    /** @type {Props} */
    let {view, data, views, parameters} = $props()
    setContext("data", data)
    setContext("navigate", function () {
        // Noop.
    })

    /**
     * @param {string} string
     */
    function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
    }

    setContext("path", path)
    setContext("view", _view)

    /**
     * @param {string} view
     * @param {Record<string,string>} [fields]
     */
    function path(view, fields = {}) {
        let result = views[view] ?? ""
        if (!views[view]) {
            return ""
        }

        for (let key in fields) {
            const value = fields[key]
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
</script>

<!--app-router-->