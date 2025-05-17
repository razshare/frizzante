<script>
    //:app-imports
    import {setContext} from 'svelte'

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
    setContext("data", dataState)
    setContext("navigate", navigate)
    setContext("findPathByPageName", findPathByPageName)
    setContext("findNavigateByPath", findNavigateByPath)

    function navigate(){
        // Noop.
    }

    /**
     * @param {string} pageName
     * @param {Record<string,string>} [parameters]
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

    function findNavigateByPath() {
        return function () {
            // Noop.
        }
    }
</script>

<!--app-router-->